package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/logger"
	"log"
	"novaro-server/dao"
	"novaro-server/model"
	"sync"
	"time"
)

type TagsRecordService struct {
	dao     *dao.TagsRecordDao
	postDao *dao.PostsDao
	tagsDao *dao.TagsDao
	userDao *dao.UsersDao
	rdb     *redis.Client
	mq      *amqp091.Connection
}

func NewTagsRecordService() *TagsRecordService {
	db := model.GetDB()
	return &TagsRecordService{
		dao:     dao.NewTagsRecordDao(db),
		postDao: dao.NewPostsDao(db),
		tagsDao: dao.NewTagsDao(db),
		userDao: dao.NewUsersDao(db),
		rdb:     model.GetRedisCli(),
		mq:      model.GetRabbitMqCli(),
	}
}

func (s *TagsRecordService) Create(records *model.TagsRecords) error {
	tagExists, err := s.tagsDao.TagExists(records.TagId)
	if err != nil || !tagExists {
		return fmt.Errorf("tag with id %s does not exist", records.TagId)
	}

	post, err := s.postDao.GetPostsById(records.PostId)
	if err != nil {
		logger.Error("post is not exist", logger.Err(err))
		return fmt.Errorf("get post error: %v", err)
	}

	user, err := s.userDao.GetById(records.UserId)
	if err != nil {
		return fmt.Errorf("get user error: %v", err)
	}

	var err1 error
	exists, source := s.RecordsExists(records.TagId, records.PostId, records.UserId)
	if exists {
		err1 = s.removeRecords(records, source)
	} else {
		err1 = s.addRecords(records, post, user)
	}

	return err1
}

func (s *TagsRecordService) RecordsExists(tagId, postId, userId string) (bool, string) {

	result, _ := s.rdb.SMembers(context.Background(), fmt.Sprintf("user:tags:%s:%s", userId, postId)).Result()
	if len(result) == 0 {
		count := s.dao.GetRecord(tagId, postId, userId)
		return count > 0, "db"
	}
	return len(result) > 0, "cache"

}

func (s *TagsRecordService) SyncData() {
	s.consumeFromRabbitMQ()
	log.Println("sync data from rabbitmq")
}

func (s *TagsRecordService) addRecords(r *model.TagsRecords, post *model.Posts, user *model.Users) error {
	err2 := s.addCacheRecords(r.UserId, r.PostId, r.TagId)
	if err2 != nil {
		return fmt.Errorf("exec error: %v", err2)
	}

	record := model.TagRecordQueue{
		TagId:      r.TagId,
		PostId:     r.PostId,
		UserId:     r.UserId,
		Points:     0,
		PostPoints: 0,
		Operation:  "a",
		CreatedAt:  time.Now(),
	}
	return s.sendToRabbitMQ(&record)
}

func (s *TagsRecordService) removeRecords(r *model.TagsRecords, source string) error {
	if source == "db" {
		err := s.dao.Delete(r)
		return err
	}

	err := s.removeCacheRecords(r.UserId, r.PostId, r.TagId)
	if err != nil {
		return fmt.Errorf("exec error: %v", err)
	}
	queue := model.TagRecordQueue{
		TagId:      r.TagId,
		PostId:     r.PostId,
		UserId:     r.UserId,
		Points:     0,
		PostPoints: 0,
		Operation:  "r",
		CreatedAt:  time.Now(),
	}
	err = s.sendToRabbitMQ(&queue)
	return nil
}

func (s *TagsRecordService) sendToRabbitMQ(q *model.TagRecordQueue) error {
	ch, err := s.mq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(q)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",              // exchange
		"records_queue", // routing key
		false,           // mandatory
		false,           // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	log.Println("发送成功")
	return err
}

func (s *TagsRecordService) consumeFromRabbitMQ() error {
	ch, err := s.mq.Channel()
	if err != nil {
		return fmt.Errorf("无法打开通道: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"records_queue", // 队列名
		true,            // 持久化
		false,           // 不自动删除
		false,           // 非排他
		false,           // 不等待
		nil,             // 无额外参数
	)
	if err != nil {
		return fmt.Errorf("无法声明队列: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return fmt.Errorf("无法设置 QoS: %v", err)
	}

	// 开始消费消息
	msgs, err := ch.Consume(
		q.Name, // 队列
		"",     // 消费者
		false,  // 手动确认
		false,  // 非排他
		false,  // 不等待
		false,  // 无额外参数
		nil,
	)

	if err != nil {
		return fmt.Errorf("无法注册消费者: %v", err)
	}

	var rcords []model.TagRecordQueue
	var mu sync.Mutex
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 20)

	timeout := time.After(5 * time.Second) // 设置5秒超时

	// 使用匿名函数来处理消息
	processMessage := func(msg amqp091.Delivery) {
		defer wg.Done()
		defer func() { <-semaphore }() // 释放信号量

		var queue model.TagRecordQueue
		err := json.Unmarshal(msg.Body, &queue)
		if err != nil {
			mu.Lock()
			mu.Unlock()
			msg.Nack(false, false)
			return
		}

		m := &model.TagsRecords{
			TagId:      queue.TagId,
			PostId:     queue.PostId,
			UserId:     queue.UserId,
			Points:     queue.Points,
			PostPoints: queue.PostPoints,
			CreatedAt:  queue.CreatedAt,
		}

		if queue.Operation == "a" {
			s.dao.AddTagsRecords(m)
		} else {
			s.dao.Delete(m)
		}

		mu.Lock()
		rcords = append(rcords, queue)
		mu.Unlock()
		msg.Ack(false)
	}

loop:
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				break loop
			}
			semaphore <- struct{}{} // 获取信号量
			wg.Add(1)
			go processMessage(msg)

		case <-timeout:
			mu.Lock()
			mu.Unlock()
			break loop
		}
	}

	// 等待所有工作协程完成
	wg.Wait()

	// 关闭通道
	err = ch.Close()
	if err != nil {
		mu.Lock()
		mu.Unlock()
	}
	return err
}

func (s *TagsRecordService) addCacheRecords(userId, postId, tagId string) error {
	pipeline := s.rdb.Pipeline()
	ctx := context.Background()
	key := fmt.Sprintf("tags:count:%s", postId)
	key2 := fmt.Sprintf("user:tags:%s:%s", userId, postId)

	pipeline.ZIncrBy(ctx, key, 1, tagId)
	pipeline.SAdd(ctx, key2, tagId)

	pipeline.Expire(ctx, key, 5*time.Minute)
	pipeline.Expire(ctx, key2, 5*time.Minute)
	_, err := pipeline.Exec(ctx)
	return err
}

func (s *TagsRecordService) removeCacheRecords(userId, postId, tagId string) error {
	pipeline := s.rdb.Pipeline()
	ctx := context.Background()
	key := fmt.Sprintf("tags:count:%s", postId)
	key2 := fmt.Sprintf("user:tags:%s:%s", userId, postId)

	pipeline.ZIncrBy(ctx, key, -1, tagId)
	pipeline.SRem(ctx, key2, tagId)
	_, err := pipeline.Exec(ctx)
	return err
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/logger"
	"log"
	"math"
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
	exists, _ := s.RecordsExists(records.TagId, records.PostId, records.UserId)
	if exists {
		err1 = s.removeRecords(records)
	} else {
		err1 = s.addRecords(records, &post, user)
	}

	return err1
}

func (s *TagsRecordService) RecordsExists(tagId, postId, userId string) (bool, string) {
	key := fmt.Sprintf("tags:expiry:%s", postId)
	member := fmt.Sprintf("%s:%s", tagId, userId)
	_, err := s.rdb.ZScore(context.Background(), key, member).Result()
	if err == redis.Nil {
		count := s.dao.GetRecord(tagId, postId, userId)
		return count > 0, "db"
	}
	return true, "cache"
}

func (s *TagsRecordService) SyncData() {
	s.consumeFromRabbitMQ()
	log.Println("sync data from rabbitmq")
}

func (s *TagsRecordService) addRecords(r *model.TagsRecords, post *model.Posts, user *model.Users) error {
	count, err := s.addCacheRecords(r.UserId, r.PostId, r.TagId)
	if err != nil {
		logger.Errorf("add cache records error: %v", err)
		return err
	}

	points := s.TagPoints(user.WalletPublicKey, count)
	postUser, err := s.userDao.GetById(post.UserId)

	var postPoints int64
	if err != nil || postUser.NftInfo == nil {
		postPoints = 0
	} else {
		postPoints = s.dao.Points(postUser.WalletPublicKey, postUser.NftInfo.Level)
	}

	record := model.TagRecordQueue{
		TagId:      r.TagId,
		PostId:     r.PostId,
		UserId:     r.UserId,
		Points:     points,
		PostPoints: float64(postPoints),
		Operation:  "a",
		CreatedAt:  time.Now(),
	}
	return s.sendToRabbitMQ(&record)
}

func (s *TagsRecordService) removeRecords(r *model.TagsRecords) error {
	err := s.removeCacheRecords(r.UserId, r.PostId, r.TagId)
	if err != nil {
		logger.Errorf("remove cache records error: %v", err)
		return err
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

func (s *TagsRecordService) addCacheRecords(userId, postId, tagId string) (int64, error) {
	pipeline := s.rdb.Pipeline()
	key := fmt.Sprintf("tags:expiry:%s", postId)
	countKey := fmt.Sprintf("tags:count:%s", postId)
	ctx := context.Background()
	now := float64(time.Now().Unix())
	member := fmt.Sprintf("%s:%s", tagId, userId)
	pipeline.ZAdd(ctx, key, redis.Z{Score: now, Member: member}).Result()
	pipeline.HIncrBy(ctx, countKey, tagId, 1)

	key3 := fmt.Sprintf("user:tags:expiry:%s", userId)
	count, err2 := pipeline.Get(ctx, key3).Int()
	if err2 == redis.Nil {
		count = 0
	}
	incr := pipeline.Incr(ctx, key3)

	if count == 0 {
		now := time.Now()
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		ttl := tomorrow.Sub(now)
		pipeline.Expire(ctx, key3, ttl)
	}

	_, err3 := pipeline.Exec(ctx)

	if err3 != nil && err3 != redis.Nil {
		return 0, err3
	}
	return incr.Result()
}

func (s *TagsRecordService) removeCacheRecords(userId, postId, tagId string) error {
	pipeline := s.rdb.Pipeline()
	ctx := context.Background()
	key := fmt.Sprintf("tags:expiry:%s", postId)
	member := fmt.Sprintf("%s:%s", tagId, userId)
	pipeline.ZRem(ctx, key, member)

	key3 := fmt.Sprintf("user:tags:expiry:%s", userId)
	pipeline.Decr(ctx, key3)
	_, err := pipeline.Exec(ctx)
	return err
}

func (s *TagsRecordService) TagPoints(wattle *string, count int64) float64 {
	if wattle == nil {
		return 0
	}

	totalPoints := float64(1)

	if count > 0 {
		count = count / 10
	}

	coefficients := float64(count)
	rewards := (totalPoints - coefficients) * totalPoints
	return math.Max(0, rewards)
}

func (s *TagsRecordService) CleanExpiredTags() {
	expirationDuration := 1 * time.Minute
	ctx := context.Background()

	// 计算过期时间点
	expirationThreshold := time.Now().Add(-expirationDuration).Unix()

	var totalRemoved int64
	var mu sync.Mutex

	workerCount := 10
	jobs := make(chan string, 100)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for key := range jobs {
			removed, err := s.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", expirationThreshold)).Result()
			if err != nil {
				logger.Errorf("error removing expired tags for key %s: %v", key, err)
			} else {
				mu.Lock()
				totalRemoved += removed
				mu.Unlock()
			}
		}
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker()
	}

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = s.rdb.Scan(ctx, cursor, "tags:expiry:*", 100).Result()
		if err != nil {
			logger.Errorf("error scanning keys: %v", err)
			continue
		}

		for _, key := range keys {
			jobs <- key
		}

		if cursor == 0 {
			break
		}
	}

	close(jobs)
	wg.Wait()
	logger.Debugf("expired tags total:%d", totalRemoved)
}

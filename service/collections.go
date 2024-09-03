package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type CollectionsService struct {
	dao *dao.CollectionsDao
}

func NewCollectionsService() *CollectionsService {
	return &CollectionsService{
		dao: dao.NewCollectionsDao(config.DB),
	}
}

func (s *CollectionsService) AddOrRemove(c *model.Collections) error {
	exist, err := NewUserService().UserExists(c.UserId)
	if err != nil || exist == false {
		return fmt.Errorf("userId is not exist")
	}

	postExist, err := NewPostService().PostExists(c.PostId)
	if err != nil || postExist == false {
		return fmt.Errorf("postId is not exist")
	}

	var errs error
	if collectionsExist := s.CollectionsExist(c.UserId, c.PostId); collectionsExist == true {
		errs = unCollectionsTweet(c)
	} else {
		errs = collectionsTweet(c)
	}

	return errs
}

// 获取用户收藏的推文
func collectionsTweet(c *model.Collections) error {
	ctx := context.Background()
	pipeline := config.RDB.Pipeline()

	// 将用户添加到推文的收藏集合中
	key := fmt.Sprintf("tweet:%s:collections", c.PostId)
	pipeline.SAdd(ctx, key, c.UserId)
	pipeline.Expire(ctx, key, 5*time.Minute)
	// 计数
	pipeline.ZIncrBy(ctx, "tweet:collections:count", 1, c.PostId)

	_, err := pipeline.Exec(ctx)

	q := model.Queue{
		Collections: *c,
		Operation:   "add",
	}

	go func() {
		err := sendToRabbitMQ(q)
		if err != nil {
			// 在这里处理错误，可以记录日志或者重试
			fmt.Printf("Error sending to RabbitMQ: %v\n", err)
		}
	}()
	return err
}

// 从收藏中移除推文
func unCollectionsTweet(c *model.Collections) error {
	// 删除redis缓存
	ctx := context.Background()
	pipeline := config.RDB.Pipeline()

	// 将用户移除推文的收藏集合中
	key := fmt.Sprintf("tweet:%s:collections", c.PostId)
	pipeline.SRem(ctx, key, c.UserId)

	pipeline.ZIncrBy(ctx, "tweet:collections:count", -1, c.PostId)

	_, err := pipeline.Exec(ctx)

	q := model.Queue{
		Collections: *c,
		Operation:   "remove",
	}
	go func() {
		err := sendToRabbitMQ(q)
		if err != nil {
			// 在这里处理错误，可以记录日志或者重试
			fmt.Printf("Error sending to RabbitMQ: %v\n", err)
		}
	}()
	return err
}

func (s *CollectionsService) CollectionsExist(userId string, postId string) bool {
	key := fmt.Sprintf("tweet:%s:collections", postId)
	result, err := config.RDB.SIsMember(context.Background(), key, userId).Result()
	if err != nil {
		e, _ := s.dao.CollectionsExist(userId, postId)
		return e
	}
	return result
}

// 发送消息到RabbitMQ
func sendToRabbitMQ(q model.Queue) error {
	ch, err := config.RBQ.Channel()
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
		"",                  // exchange
		"collections_queue", // routing key
		false,               // mandatory
		false,               // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	log.Println("发送成功")
	return err
}

// 从RabbitMQ中消费消息
func consumeFromRabbitMQ() ([]model.Queue, error) {
	ch, err := config.RBQ.Channel()
	if err != nil {
		return nil, fmt.Errorf("无法打开通道: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"collections_queue", // 队列名
		true,                // 持久化
		false,               // 不自动删除
		false,               // 非排他
		false,               // 不等待
		nil,                 // 无额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("无法声明队列: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return nil, fmt.Errorf("无法设置 QoS: %v", err)
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
		return nil, fmt.Errorf("无法注册消费者: %v", err)
	}

	var collections []model.Queue
	timeout := time.After(5 * time.Second) // 设置5秒超时

	for {
		select {
		case msg, ok := <-msgs:
			log.Println("消费信息：", msg)
			if !ok {
				return collections, nil // 通道已关闭，返回已收集的数据
			}
			var qu model.Queue
			err := json.Unmarshal(msg.Body, &qu)
			if err != nil {
				log.Printf("解析消息失败: %v", err)
				msg.Nack(false, true) // 消息解析失败，重新入队
				continue
			}
			collections = append(collections, qu)
			msg.Ack(false) // 手动确认消息
		case <-timeout:
			return collections, nil // 超时，返回已收集的数据
		}
	}
}

// 将记录同步到数据库
func (s *CollectionsService) SyncToDatabase() {
	messages, err := consumeFromRabbitMQ()
	if messages == nil || err != nil {
		log.Println("Failed to consumer from rabbitmq:", err)
		return
	}

	err = s.loadingData(messages)

}

// 刷新数据
func (s *CollectionsService) loadingData(operations []model.Queue) error {
	err := s.dao.RefreshData(operations)
	return err
}

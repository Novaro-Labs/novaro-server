package model

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"log"
	"novaro-server/config"
	"strings"
	"time"
)

// 收藏，先记录在redis中，每五分钟更新一次数据库

type Collections struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

type queue struct {
	Collections Collections
	Operation   string // add: 收藏 remove: 取消收藏
}

func (Collections) TableName() string {
	return "collections"
}

func (c *Collections) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	c.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

// 获取用户收藏的推文
func CollectionsTweet(c *Collections) error {
	ctx := context.Background()
	rdb := config.RDB
	pipeline := rdb.Pipeline()

	// 将用户添加到推文的收藏集合中
	key := fmt.Sprintf("tweet:%s:collections", c.PostId)
	pipeline.SAdd(ctx, key, c.UserId)
	pipeline.Expire(ctx, key, 5*time.Minute)
	// 计数
	pipeline.ZIncrBy(ctx, "tweet:collections:count", 1, c.PostId)

	_, err := pipeline.Exec(ctx)

	q := queue{
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
func UnCollectionsTweet(c *Collections) error {

	// 删除redis缓存
	ctx := context.Background()
	rdb := config.RDB
	pipeline := rdb.Pipeline()

	// 将用户移除推文的收藏集合中
	key := fmt.Sprintf("tweet:%s:collections", c.PostId)
	pipeline.SRem(ctx, key, c.UserId)

	pipeline.ZIncrBy(ctx, "tweet:collections:count", -1, c.PostId)

	_, err := pipeline.Exec(ctx)

	q := queue{
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

// 同步到数据库
func SyncToDatabase() {
	//获取计数
	rdb := config.RDB
	ctx := context.Background()
	pipeline := rdb.Pipeline()

	scoreCmd := pipeline.ZScore(ctx, "tweet:collections:count", "1")
	_, err := pipeline.Exec(ctx)
	score, err := scoreCmd.Result()
	log.Println("redis缓存", score)

	messages, err := consumeFromRabbitMQ()
	if messages == nil || err != nil {
		log.Println("Failed to consumer from rabbitmq:", err)
		return
	}

	err = loadingData(messages)

}

// 发送消息到RabbitMQ
func sendToRabbitMQ(q queue) error {
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
func consumeFromRabbitMQ() ([]queue, error) {
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

	var collections []queue
	timeout := time.After(5 * time.Second) // 设置5秒超时

	for {
		select {
		case msg, ok := <-msgs:
			log.Println("消费信息：", msg)
			if !ok {
				return collections, nil // 通道已关闭，返回已收集的数据
			}
			var qu queue
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

// 刷新数据
func loadingData(operations []queue) error {

	db := config.DB
	// 开始事务
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, coll := range operations {
			if coll.Operation == "add" {
				log.Println("获取到数据", coll.Collections)
				err := tx.Create(&coll.Collections).Error
				return err
			} else {
				err := tx.Where("user_id = ? and post_id = ?", coll.Collections.UserId, coll.Collections.PostId).Delete(&Collections{}).Error
				return err
			}
		}
		return nil
	})
	return err
}

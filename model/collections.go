package model

import (
	"context"
	"fmt"
	"novaro-server/config"
	"time"
)

// 收藏，先记录在redis中，每五分钟更新一次数据库
// TODO 更新，每次收藏都更新数据库

type Collections struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	createdAt time.Time `json:"createdAt"`
}

// 获取用户收藏的推文
func CollectionsTweet(c *Collections) error {
	ctx := context.Background()
	rdb := config.RDB
	pipeline := rdb.Pipeline()

	// 将用户添加到推文的收藏集合中
	pipeline.SAdd(ctx, fmt.Sprintf("tweet:%s:collections", c.PostId), c.UserId)
	// 增加推文的收藏计数
	pipeline.ZIncrBy(ctx, "tweet:collections:count", 1, c.PostId)
	// 将推文添加到用户的收藏集合中
	pipeline.SAdd(ctx, fmt.Sprintf("user:%s:collections", c.UserId), c.PostId)

	_, err := pipeline.Exec(ctx)
	return err
}

// 从收藏中移除推文
func UnCollectionsTweet(c *Collections) error {
	ctx := context.Background()
	rdb := config.RDB
	pipe := rdb.Pipeline()

	// 减少推文的收藏计数
	pipe.ZIncrBy(ctx, "tweet:collections:count", -1, c.PostId)

	// 将用户从推文的收藏集合中移除
	pipe.SRem(ctx, fmt.Sprintf("tweet:%s:collections", c.PostId), c.UserId)

	// 将推文从用户的收藏集合中移除
	pipe.SRem(ctx, fmt.Sprintf("user:%s:collections", c.UserId), c.PostId)

	_, err := pipe.Exec(ctx)

	return err
}

package model

import (
	"context"
	"fmt"
	"novaro-server/config"
	"time"
)

// TODO 同样需要同步数据库

type RePosts struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	createdAt time.Time `json:"createdAt"`
}

func AddRePosts(c *RePosts) error {
	ctx := context.Background()
	rdb := config.RDB
	pipeline := rdb.Pipeline()

	// 将用户添加到推文的转发集合中
	pipeline.SAdd(ctx, fmt.Sprintf("tweet:%s:reposts", c.PostId), c.UserId)

	// 增加推文的转发计数
	pipeline.ZIncrBy(ctx, "tweet:reposts:count", 1, c.PostId)

	_, err := pipeline.Exec(ctx)
	return err
}

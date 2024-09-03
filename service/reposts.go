package service

import (
	"context"
	"fmt"
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
)

type RePostsService struct {
	dao *dao.RePostsDao
}

func NewRePostsService() *RePostsService {
	return &RePostsService{
		dao: dao.NewRePostsDao(config.DB),
	}
}

func (s *RePostsService) AddRePosts(c *model.RePosts) error {
	ctx := context.Background()
	pipeline := config.RDB.Pipeline()

	// 将用户添加到推文的转发集合中
	pipeline.SAdd(ctx, fmt.Sprintf("tweet:%s:reposts", c.PostId), c.UserId)

	// 增加推文的转发计数
	pipeline.ZIncrBy(ctx, "tweet:reposts:count", 1, c.PostId)

	_, err := pipeline.Exec(ctx)
	return err
}

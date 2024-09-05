package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"novaro-server/dao"
	"novaro-server/model"
)

type RePostsService struct {
	dao *dao.RePostsDao
	rdb *redis.Client
}

func NewRePostsService() *RePostsService {
	return &RePostsService{
		dao: dao.NewRePostsDao(model.GetDB()),
		rdb: model.GetRedisCli(),
	}
}

func (s *RePostsService) AddRePosts(c *model.RePosts) error {
	ctx := context.Background()
	pipeline := s.rdb.Pipeline()

	// 将用户添加到推文的转发集合中
	pipeline.SAdd(ctx, fmt.Sprintf("tweet:%s:reposts", c.PostId), c.UserId)

	// 增加推文的转发计数
	pipeline.ZIncrBy(ctx, "tweet:reposts:count", 1, c.PostId)

	_, err := pipeline.Exec(ctx)
	return err
}

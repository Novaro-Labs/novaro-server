package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"novaro-server/dao"
	"novaro-server/model"
)

type CommentsService struct {
	dao *dao.CommentsDao
	rdb *redis.Client
}

func NewCommentService() *CommentsService {
	return &CommentsService{
		dao: dao.NewCommentsDao(model.GetDB()),
		rdb: model.GetRedisCli(),
	}
}

func (s *CommentsService) Create(c *model.Comments) error {
	err := s.dao.Create(c)

	// 用redis来记录评论数量
	s.rdb.ZIncrBy(context.Background(), "tweet:comments:count", 1, c.PostId)
	return err
}

func (s *CommentsService) GetById(id string) (model.Comments, error) {
	return s.dao.GetById(id)
}

func (s *CommentsService) GetCount(postId string) int64 {
	return s.dao.GetCount(postId)
}

func (s *CommentsService) GetListByPostId(postId string) ([]model.Comments, error) {
	return s.dao.GetListByPostId(postId)
}

func (s *CommentsService) GetListByParentId(parentId string) ([]model.Comments, error) {
	if parentId == "" {
		return nil, fmt.Errorf("parentId cannot be empty")
	}

	return s.dao.GetListByParentId(parentId)
}

func (s *CommentsService) GetListByUserId(userId string) ([]model.Comments, error) {

	return s.dao.GetListByUserId(userId)
}

func (s *CommentsService) Delete(id string) error {
	err := s.dao.DeleteById(id)
	if err != nil {
		return err
	}
	resp, _ := s.dao.GetById(id)
	s.rdb.ZIncrBy(context.Background(), "tweet:comments:count", -1, resp.PostId)
	return err
}

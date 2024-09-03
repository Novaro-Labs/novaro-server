package service

import (
	"context"
	"fmt"
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
)

type CommentsService struct {
	dao *dao.CommentsDao
}

func NewCommentService() *CommentsService {
	return &CommentsService{
		dao: dao.NewCommentsDao(config.DB),
	}
}

func (s *CommentsService) Create(c *model.Comments) error {
	err := s.dao.Create(c)

	// 用redis来记录评论数量
	config.RDB.ZIncrBy(context.Background(), "tweet:comments:count", 1, c.PostId)
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
	config.RDB.ZIncrBy(context.Background(), "tweet:comments:count", -1, resp.PostId)
	return err
}

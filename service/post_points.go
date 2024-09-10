package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type PostPointsService struct {
	dao *dao.PostPointsDao
}

func NewPostPointsService() *PostPointsService {
	return &PostPointsService{
		dao: dao.NewPostPointsDao(model.GetDB()),
	}
}

func (s *PostPointsService) Save(m *model.PostPoints) error {
	return s.dao.Save(nil, m)
}

func (s *PostPointsService) Delete(postId string) error {
	return s.dao.Delete(nil, postId)
}

package service

import (
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
)

type TagsService struct {
	dao *dao.TagsDao
}

func NewTagsService() *TagsService {
	return &TagsService{
		dao: dao.NewTagsDao(config.DB),
	}
}

func (s *TagsService) GetListByPostId(postId string) ([]model.Tags, error) {
	resp, err := s.dao.GetTagListByPostId(postId)
	return resp, err
}

func (s *TagsService) GetTagsList() (resp []model.Tags, err error) {

	return s.dao.GetTagsList()
}

func (s *TagsService) TagExists(id string) (bool, error) {
	return s.dao.TagExists(id)
}

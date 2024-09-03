package service

import (
	"fmt"
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
)

type TagsRecordService struct {
	dao *dao.TagsRecordDao
}

func NewTagsRecordService() *TagsRecordService {
	return &TagsRecordService{
		dao: dao.NewTagsRecordDao(config.DB),
	}
}

func (s *TagsRecordService) Create(records *model.TagsRecords) error {

	exists, err := NewPostService().PostExists(records.PostId)
	if err != nil || !exists {
		return fmt.Errorf("post with id %s does not exist", records.PostId)
	}

	tagExists, err := NewTagsService().TagExists(records.TagId)
	if err != nil || !tagExists {
		return fmt.Errorf("tag with id %s does not exist", records.TagId)
	}

	return s.dao.AddTagsRecords(records)

}

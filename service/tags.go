package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"novaro-server/dao"
	"novaro-server/model"
)

type TagsService struct {
	dao *dao.TagsDao
	rdb *redis.Client
}

func NewTagsService() *TagsService {
	return &TagsService{
		dao: dao.NewTagsDao(model.GetDB()),
		rdb: model.GetRedisCli(),
	}
}

func (s *TagsService) GetListByPostId(postId string) ([]model.Tags, error) {
	resp, err := s.dao.GetTagListByPostId(postId)
	return resp, err
}

func (s *TagsService) GetTagsList() (resp []model.Tags, err error) {
	key := "tags_list"
	ctx := context.Background()
	result, err := s.rdb.HGetAll(ctx, key).Result()

	if err != nil || len(result) == 0 {
		list, err := s.dao.GetTagsList()
		m := map[string]string{}
		for _, item := range list {
			m[item.Id] = item.SourceId
		}
		err = s.rdb.HSet(ctx, key, m).Err()
		return list, err
	}

	return mapToStruct(result), nil
}

func mapToStruct(m map[string]string) []model.Tags {
	fmt.Println("1111111111111111")
	var tags []model.Tags
	for id, url := range m {
		tags = append(tags, model.Tags{Id: id, SourceId: url})
	}
	return tags
}

func (s *TagsService) TagExists(id string) (bool, error) {
	return s.dao.TagExists(id)
}

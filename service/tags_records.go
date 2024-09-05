package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/logger"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type TagsRecordService struct {
	dao     *dao.TagsRecordDao
	postDao *dao.PostsDao
	tagsDao *dao.TagsDao
	userDao *dao.UsersDao
	rdb     *redis.Client
}

func NewTagsRecordService() *TagsRecordService {
	db := model.GetDB()
	return &TagsRecordService{
		dao:     dao.NewTagsRecordDao(db),
		postDao: dao.NewPostsDao(db),
		tagsDao: dao.NewTagsDao(db),
		userDao: dao.NewUsersDao(db),
		rdb:     model.GetRedisCli(),
	}
}

func (s *TagsRecordService) Create(records *model.TagsRecords) error {
	tagExists, err := s.tagsDao.TagExists(records.TagId)
	if err != nil || !tagExists {
		return fmt.Errorf("tag with id %s does not exist", records.TagId)
	}
	return s.addTags(records)

}

func (s *TagsRecordService) addTags(r *model.TagsRecords) error {

	post, err1 := s.postDao.GetPostsById(r.PostId)
	if err1 != nil {
		logger.Error("post is not exist", logger.Err(err1))
		return fmt.Errorf("get post error: %v", err1)
	}

	user, err3 := s.userDao.GetById(r.UserId)
	if err3 != nil || user.WalletPublicKey == nil {
		return fmt.Errorf("get user error: %v", err3)
	}

	_, err4 := s.userDao.GetById(post.UserId)
	if err4 != nil {
		return fmt.Errorf("get post user error: %v", err4)
	}

	pipeline := s.rdb.Pipeline()
	ctx := context.Background()

	pipeline.SAdd(ctx, fmt.Sprintf("user:%s:tags", r.UserId), r.PostId, r.TagId)
	pipeline.Expire(ctx, fmt.Sprintf("user:%s:tags", r.UserId), 5*time.Minute)

	_, err2 := pipeline.Exec(ctx)
	if err2 != nil {
		return fmt.Errorf("exec error: %v", err2)
	}

	s.recordCount(r)
	//_ := model.TagRecordQueue{
	//	TagId:  r.TagId,
	//	PostId: r.PostId,
	//	UserId: r.UserId,
	//	Points: 1,
	//}
	return nil
}

func (s *TagsRecordService) recordCount(r *model.TagsRecords) (bool, error) {

	pipeline := s.rdb.Pipeline()
	ctx := context.Background()

	result, _ := pipeline.SMembers(ctx, fmt.Sprintf("user:%s:tags", r.UserId)).Result()

	for _, member := range result {
		fmt.Println(member)
	}
	return false, nil

}

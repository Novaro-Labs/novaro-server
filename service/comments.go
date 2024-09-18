package service

import (
	"context"
	"fmt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"novaro-server/dao"
	"novaro-server/model"
	"strings"
	"sync"
)

type CommentsService struct {
	dao     *dao.CommentsDao
	postDao *dao.PostsDao
	userDao *dao.UsersDao
}

func NewCommentService() *CommentsService {
	db := model.GetDB()
	return &CommentsService{
		dao:     dao.NewCommentsDao(db),
		postDao: dao.NewPostsDao(db),
		userDao: dao.NewUsersDao(db),
	}
}

func (s *CommentsService) Create(c *model.Comments) (*model.Comments, int64, error) {
	item, count, err := s.dao.Create(c)
	user, err := s.userDao.GetById(item.UserId)
	item.User = &model.Users{
		Id:       user.Id,
		UserName: user.UserName,
	}

	return item, count, err
}

func (s *CommentsService) GetById(id string) (*model.Comments, error) {
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

func (s *CommentsService) Delete(id, postId string) (int64, error) {
	count, err := s.dao.DeleteById(id, postId)

	return count, err
}

func (s *CommentsService) SyncCommentsToDB() {
	result, client := s.dao.SyncCommentsToDB()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for _, key := range result {
		wg.Add(1)
		semaphore <- struct{}{} // 获取信号量
		go func(key string) {
			defer wg.Done()
			defer func() { <-semaphore }() // 释放信号量

			postID := extractPostIDFromKey(key)
			count, err := client.Get(context.Background(), key).Int64()
			if err != nil {
				logger.Error("redis get comment count error", logger.Err(err), logger.String("key", key))
				return
			}

			err = s.postDao.UpdateCount(postID, count)
			if err != nil {
				logger.Error("update comment count in DB error", logger.Err(err), logger.String("postID", postID))
			} else {
				logger.Info("successfully synced comment count", logger.String("postID", postID), logger.Int64("count", count))
			}
		}(key)
	}

	wg.Wait() // 等待所有协程完成
}

func extractPostIDFromKey(key string) string {
	// 假设 key 的格式是 "post:123:comment_count"
	parts := strings.Split(key, ":")
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

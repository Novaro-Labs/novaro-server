package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type LikeService struct {
	dao         *dao.LikeDao
	postDao     *dao.PostsDao
	redisClient *redis.Client
}

func NewLikeService() *LikeService {
	return &LikeService{
		dao:         dao.NewLikeDao(model.GetDB()),
		postDao:     dao.NewPostsDao(model.GetDB()),
		redisClient: model.GetRedisCli(),
	}
}

func (s *LikeService) Like(like *model.LikesReq) error {
	key := fmt.Sprintf("likes:%s", like.PostId)
	ctx := context.Background()

	var err error
	isMember := s.IsLikeByUser(ctx, like.PostId, like.UserId, key)
	if isMember {
		err = s.Delete(ctx, like.PostId, like.UserId)
	} else {
		err = s.Add(ctx, like.PostId, like.UserId)
	}

	return err
}

func (s *LikeService) Add(ctx context.Context, postId, userId string) error {
	key := fmt.Sprintf("likes:%s", postId)

	pipe := s.redisClient.Pipeline()
	pipe.SAdd(ctx, key, userId)
	pipe.Expire(ctx, key, 5*time.Minute)
	pipe.SAdd(ctx, "likes:pending", postId)

	_, err := pipe.Exec(ctx)
	return err
}

func (s *LikeService) Delete(ctx context.Context, postId, userId string) error {
	key := fmt.Sprintf("likes:%s", postId)

	_, err := s.redisClient.SRem(ctx, key, userId).Result()
	if err != nil {
		return err
	}

	// del key
	count, err := s.redisClient.SCard(ctx, key).Result()
	if count == 0 {
		_, err = s.redisClient.Del(ctx, key).Result()
		if err != nil {
			return err
		}
	}

	likeId := s.dao.IsLikedByUser(postId, userId)
	if likeId != "" {
		s.dao.Delete(likeId)
		s.postDao.UpdateLikeAmount(postId, 1, "sub")
	}

	return nil
}

func (s *LikeService) IsLikeByUser(ctx context.Context, postId, userId, key string) bool {
	isMember, _ := s.redisClient.SIsMember(ctx, key, userId).Result()
	if isMember {
		return true
	}

	likeId := s.dao.IsLikedByUser(postId, userId)
	if likeId != "" {
		return true
	}
	return false
}

func (s *LikeService) FlushToDatabase() {
	ctx := context.Background()
	postIDs, err := s.redisClient.SMembers(ctx, "likes:pending").Result()
	if err != nil {
		return
	}

	for _, postID := range postIDs {
		key := fmt.Sprintf("likes:%s", postID)
		// 获取该postID的所有点赞用户
		userIDs, err := s.redisClient.SMembers(ctx, key).Result()
		if err != nil {
			log.Printf("Error getting likes for post %s: %v", postID, err)
			continue
		}

		err = s.dao.BatchAdd(postID, userIDs)
		if err != nil {
			log.Printf("Error updating database: %v,%s", err, postID)
		}

		err = s.postDao.UpdateLikeAmount(postID, len(userIDs), "add")
		if err != nil {
			log.Printf("Error updating amount: %v,%s,%d", err, postID, len(userIDs))
		}

		log.Printf("Updating database: Post %s liked by %d users", postID, len(userIDs))

		for _, userId := range userIDs {
			_, err2 := s.redisClient.SRem(ctx, key, userId).Result()
			if err2 != nil {
				continue
			}
		}

		count, _ := s.redisClient.SCard(ctx, key).Result()
		if count == 0 {
			s.redisClient.Del(ctx, key).Result()
		}

		s.redisClient.SRem(ctx, "likes:pending", postID)
	}

}

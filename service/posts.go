package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
	"strings"
	"sync"
)

type PostService struct {
	dao *dao.PostsDao
}

func NewPostService() *PostService {
	return &PostService{
		dao: dao.NewPostsDao(config.DB),
	}
}

func (s *PostService) GetList(p *model.PostsQuery) (resp []model.Posts, err error) {
	query := config.DB.Table("posts")
	if p.Id != "" {
		query = query.Where("id = ?", p.Id)
	}

	err = query.Find(&resp).Error

	var wg sync.WaitGroup
	// 使用缓冲通道作为信号量来限制并发 goroutine 的数量
	semaphore := make(chan struct{}, 10) // 最多10个并发 goroutine

	for i := range resp {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			// 检查收藏状态
			resp[i].IsCollected = NewCollectionsService().CollectionsExist(p.UserId, resp[i].Id)

			// 获取标签
			tags, err := NewTagsService().GetListByPostId(resp[i].Id)
			if err != nil {
				resp[i].Tags = nil
			} else {
				resp[i].Tags = tags
			}

			// 获取图片
			source, err2 := NewImgsService().GetBySourceId(resp[i].SourceId)
			if err2 != nil {
				resp[i].Imgs = nil
			} else {
				resp[i].Imgs = source
			}

		}(i)
	}

	wg.Wait()

	return resp, err
}

func (s *PostService) GetById(id string) (model.Posts, error) {
	resp, err := s.dao.GetPostsById(id)
	tags, err := NewTagsService().GetListByPostId(resp.Id)
	resp.Tags = tags
	return resp, err
}

func (s *PostService) PostExists(id string) (bool, error) {
	return s.dao.PostExists(id)
}

func (s *PostService) GetByUserId(userId string) ([]model.Posts, error) {
	resp, err := s.dao.GetPostsByUserId(userId)

	// 处理标签
	for i := range resp {
		tags, err := NewTagsService().GetListByPostId(resp[i].Id)
		if err != nil {
			resp[i].Tags = nil
		}
		resp[i].Tags = tags
	}
	return resp, err
}

func (s *PostService) Save(posts *model.Posts) error {
	return s.dao.Save(posts)
}

func (s *PostService) Update(posts *model.Posts) error {
	return s.dao.Update(posts)
}

func (s *PostService) UpdateBatch(posts []model.Posts) error {
	return s.dao.UpdateBatch(posts)
}

func (s *PostService) Delete(id string) error {
	return s.dao.Delete(id)
}

func (s *PostService) SyncData() error {
	rdb := config.RDB
	ctx := context.Background()
	result, err := rdb.ZRange(ctx, "tweet:collections:count", 0, -1).Result()

	if err != nil {
		return fmt.Errorf("failed to get tweet IDs from Redis: %v", err)
	}

	updateChan := make(chan model.Posts, len(result))
	errChan := make(chan error, len(result))
	var wg sync.WaitGroup

	for _, tweetID := range result {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()
			data, err := s.processTweet(ctx, rdb, id)
			if err != nil {
				errChan <- err
				return
			}
			updateChan <- data
		}(tweetID)
	}

	go func() {
		wg.Wait()
		close(updateChan)
		close(errChan)
	}()

	// 收集所有更新
	var updates []model.Posts
	for data := range updateChan {
		updates = append(updates, data)
	}

	// 检查是否有错误发生
	for err := range errChan {
		log.Printf("Error processing tweet: %v", err)
	}

	// 批量更新数据库
	if err := s.dao.UpdateBatch(updates); err != nil {
		return fmt.Errorf("error updating database: %v", err)
	}
	log.Println("Data sync completed")
	return err
}

func (s *PostService) processTweet(ctx context.Context, rdb *redis.Client, tweetID string) (model.Posts, error) {
	resp, err := s.dao.GetPostsById(tweetID)
	if err != nil {
		return model.Posts{}, fmt.Errorf("error getting tweet %s: %v", tweetID, err)
	}

	score, err := rdb.ZScore(ctx, "tweet:collections:count", tweetID).Result()
	repost, err := rdb.ZScore(ctx, "tweet:reposts:count", tweetID).Result()

	count := NewCommentService().GetCount(tweetID)
	return model.Posts{
		Id:                tweetID,
		CollectionsAmount: int(score) + resp.CollectionsAmount,
		RepostsAmount:     int(repost),
		CommentsAmount:    int(count),
	}, nil
}

func (s *PostService) SyncCountToDataBase() {
	rdb := config.RDB
	ctx := context.Background()

	// 同步收藏数量
	collectionsResults, err := rdb.ZRangeWithScores(ctx, "tweet:collections:count", 0, -1).Result()
	if err != nil {
		log.Printf("failed to get collection scores: %v", err)
	}

	// 同步评论数量
	commentsResults, err := rdb.ZRangeWithScores(ctx, "tweet:comments:count", 0, -1).Result()
	if err != nil {
		log.Printf("failed to get comment scores: %v", err)
	}

	// 同步转发数量
	repostsResults, err := rdb.ZRangeWithScores(ctx, "tweet:reposts:count", 0, -1).Result()
	if err != nil {
		log.Printf("failed to get repost scores: %v", err)
	}

	var wg sync.WaitGroup
	errorChan := make(chan error, len(collectionsResults)+len(commentsResults)+len(repostsResults))

	// 处理收藏数量
	for _, result := range collectionsResults {
		wg.Add(1)
		go func(r redis.Z) {
			defer wg.Done()
			// 更新收藏数量
			postId := r.Member.(string)

			err2 := s.dao.Update(&model.Posts{
				Id:                postId,
				CollectionsAmount: int(r.Score),
			})
			if err2 != nil {
				errorChan <- fmt.Errorf("failed to update post %s: %v", postId, err2)
			}

		}(result)
	}

	// 处理评论数量
	for _, result := range commentsResults {
		wg.Add(1)
		go func(r redis.Z) {
			defer wg.Done()
			// 更新评论数量
			postId := r.Member.(string)

			err2 := s.dao.Update(&model.Posts{
				Id:             postId,
				CommentsAmount: int(r.Score),
			})

			if err2 != nil {
				errorChan <- fmt.Errorf("failed to update post %s: %v", postId, err2)
			}

		}(result)
	}

	// 处理转发数量
	for _, result := range repostsResults {
		wg.Add(1)
		go func(r redis.Z) {
			defer wg.Done()
			// 更新转发数量
			postId := r.Member.(string)

			err2 := s.dao.Update(&model.Posts{
				Id:            postId,
				RepostsAmount: int(r.Score),
			})
			if err2 != nil {
				errorChan <- fmt.Errorf("failed to update post %s: %v", postId, err2)
			}
		}(result)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	var errors []string
	for err := range errorChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		fmt.Errorf("some updates failed: %s", strings.Join(errors, "; "))
	}
}

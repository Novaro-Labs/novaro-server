package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"novaro-server/config"
	"strings"
	"sync"
	"time"
)

type Posts struct {
	Id                string      `json:"id"`
	UserId            string      `json:"userId"`
	Content           string      `json:"content"`
	CommentsAmount    int         `json:"commentsAmount"`
	CollectionsAmount int         `json:"collectionsAmount"`
	RepostsAmount     int         `json:"repostsAmount"`
	CreatedAt         time.Time   `json:"createdAt"`
	OriginalId        string      `json:"originalId"`
	SourceId          string      `json:"sourceId"`
	Tags              []Tags      `json:"tags" gorm:"many2many:tags_records;"`
	PostsImgs         []PostsImgs `json:"postsImgs" gorm:"many2many:posts_imgs;"`
	IsCollected       bool        `json:"isCollected" gorm:"-"`
}

func (Posts) TableName() string {
	return "posts"
}

type PostsQuery struct {
	Id     string `form:"id" json:"id"`
	UserId string `form:"userId" json:"userId"`
}

func (u *Posts) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func GetPostsList(p *PostsQuery) (resp []Posts, err error) {
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
			resp[i].IsCollected = CollectionsExist(p.UserId, resp[i].Id)

			// 获取标签
			tags, err := GetTagListByPostId(resp[i].Id)
			if err != nil {
				resp[i].Tags = nil
			} else {
				resp[i].Tags = tags
			}

			// 获取图片
			source, err2 := GetPostImgsBySourceId(resp[i].SourceId)
			if err2 != nil {
				resp[i].PostsImgs = nil
			} else {
				resp[i].PostsImgs = source
			}

		}(i)
	}

	wg.Wait()

	return resp, err
}

func GetPostsById(id string) (resp Posts, err error) {
	if id == "" {
		return resp, errors.New("id is required")
	}
	err = config.DB.Table("posts").Where("id = ?", id).Find(&resp).Error

	// 处理标签
	tags, err := GetTagListByPostId(resp.Id)
	resp.Tags = tags
	return resp, err
}

func PostExists(id string) (bool, error) {
	var count int64
	err := config.DB.Model(&Posts{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func GetPostsByUserId(userId string) (resp []Posts, err error) {
	if userId == "" {
		return resp, errors.New("UserId is required")
	}
	err = config.DB.Table("posts").Where("user_id = ?", userId).Find(&resp).Error

	// 处理标签
	for i := range resp {
		tags, err := GetTagListByPostId(resp[i].Id)
		if err != nil {
			resp[i].Tags = nil
		}
		resp[i].Tags = tags
	}
	return resp, nil
}

func SavePosts(posts *Posts) error {
	var data = Posts{
		Id:         posts.Id,
		UserId:     posts.UserId,
		Content:    posts.Content,
		OriginalId: posts.OriginalId,
		SourceId:   posts.SourceId,
	}

	tx := config.DB.Create(&data)

	if data.OriginalId != "" {
		rdb := config.RDB
		ctx := context.Background()
		rdb.ZIncrBy(ctx, "tweet:reposts:count", 1, data.Id)
	}

	return tx.Error
}

func UpdatePosts(posts *Posts) error {
	tx := config.DB.Updates(&posts)
	return tx.Error
}

func UpdatePostsBatch(posts []Posts) error {
	// 开始事务
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, post := range posts {
			// 更新每个 post
			if err := tx.Model(&post).Updates(Posts{
				Content:           post.Content,
				CommentsAmount:    post.CommentsAmount,
				CollectionsAmount: post.CollectionsAmount,
				RepostsAmount:     post.RepostsAmount,
				Tags:              post.Tags,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func DelPostsById(id string) error {
	tx := config.DB.Where("id = ?", id).Delete(&Posts{})
	return tx.Error
}

func SyncData() error {
	rdb := config.RDB
	ctx := context.Background()
	result, err := rdb.ZRange(ctx, "tweet:collections:count", 0, -1).Result()

	if err != nil {
		return fmt.Errorf("failed to get tweet IDs from Redis: %v", err)
	}

	updateChan := make(chan Posts, len(result))
	errChan := make(chan error, len(result))
	var wg sync.WaitGroup

	for _, tweetID := range result {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()
			data, err := processTweet(ctx, rdb, id)
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
	var updates []Posts
	for data := range updateChan {
		updates = append(updates, data)
	}

	// 检查是否有错误发生
	for err := range errChan {
		log.Printf("Error processing tweet: %v", err)
	}

	// 批量更新数据库
	if err := UpdatePostsBatch(updates); err != nil {
		return fmt.Errorf("error updating database: %v", err)
	}
	log.Println("Data sync completed")
	return err
}

func processTweet(ctx context.Context, rdb *redis.Client, tweetID string) (Posts, error) {

	resp, err := GetPostsById(tweetID)
	if err != nil {
		return Posts{}, fmt.Errorf("error getting tweet %s: %v", tweetID, err)
	}

	score, err := rdb.ZScore(ctx, "tweet:collections:count", tweetID).Result()
	repost, err := rdb.ZScore(ctx, "tweet:reposts:count", tweetID).Result()

	count := GetCommentsCount(tweetID)
	return Posts{
		Id:                tweetID,
		CollectionsAmount: int(score) + resp.CollectionsAmount,
		RepostsAmount:     int(repost),
		CommentsAmount:    int(count),
	}, nil
}

func SyncCountToDataBase() {
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

			err2 := UpdatePosts(&Posts{
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

			err2 := UpdatePosts(&Posts{
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

			err2 := UpdatePosts(&Posts{
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

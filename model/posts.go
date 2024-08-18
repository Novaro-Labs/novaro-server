package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"novaro-server/config"
	"sync"
	"time"
)

type Posts struct {
	Id                string    `json:"id"`
	UserId            string    `json:"userId"`
	Content           string    `json:"content"`
	CommentsAmount    int       `json:"commentsAmount"`
	CollectionsAmount int       `json:"collectionsAmount"`
	RepostsAmount     int       `json:"repostsAmount"`
	CreatedAt         time.Time `json:"createdAt"`
	Tags              []Tags    `json:"tags" gorm:"-"`
}

type PostsQuery struct {
	Id     string `form:"id" json:"id"`
	UserId string `form:"userId" json:"userId"`
}

func GetPostsList(p *PostsQuery) (resp []Posts, err error) {
	db := config.DB
	query := db.Model(Posts{})

	if p.UserId != "" {
		query = query.Where("user_id = ?", p.UserId)
	}
	if p.Id != "" {
		query = query.Where("id = ?", p.Id)
	}

	err = query.Find(&resp).Error

	// 处理标签
	for i := range resp {
		tags, err := GetTagListByPostId(resp[i].Id)
		if err != nil {
			resp[i].Tags = nil
		}
		resp[i].Tags = tags
	}

	return resp, err
}

func GetPostsById(id string) (resp Posts, err error) {
	db := config.DB

	if id == "" {
		return resp, errors.New("id is required")
	}
	err = db.Table("posts").Where("id = ?", id).Find(&resp).Error

	// 处理标签
	tags, err := GetTagListByPostId(resp.Id)
	resp.Tags = tags
	return resp, err
}

func GetPostsByUserId(userId string) (resp []Posts, err error) {
	db := config.DB

	if userId == "" {
		return resp, errors.New("UserId is required")
	}
	err = db.Table("posts").Where("user_id = ?", userId).Find(&resp).Error

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
	db := config.DB
	var data = Posts{
		Id:      posts.Id,
		UserId:  posts.UserId,
		Content: posts.Content,
	}

	tx := db.Create(&data)
	return tx.Error
}

func UpdatePosts(posts *Posts) error {
	db := config.DB
	tx := db.Updates(&posts)
	return tx.Error
}

func UpdatePostsBatch(posts []Posts) error {
	db := config.DB
	// 开始事务
	err := db.Transaction(func(tx *gorm.DB) error {
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
	db := config.DB
	tx := db.Where("id = ?", id).Delete(&Posts{})
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

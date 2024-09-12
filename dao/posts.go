package dao

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"novaro-server/model"
	"time"
)

type PostsDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewPostsDao(db *gorm.DB) *PostsDao {
	return &PostsDao{
		db:  db,
		rdb: model.GetRedisCli(),
	}
}

func (d *PostsDao) GetPostsList(p *model.PostsQuery) (resp []model.Posts, err error) {
	query := d.db.Table("posts").Limit(p.Size).Offset(p.Page * p.Size)
	if p.Id != "" {
		query = query.Where("id = ?", p.Id)
	}

	err = query.Order("created_at desc").Find(&resp).Error
	return resp, err
}

func (d *PostsDao) GetPostsById(id string) (resp *model.Posts, err error) {
	err = d.db.Table("posts").Where("id = ?", id).Find(&resp).Error
	return resp, err
}

func (d *PostsDao) GetPostIdByUser(postId string) *model.Posts {
	var post model.Posts
	d.db.Model(model.Posts{}).Preload("User").Where("id = ?", postId).First(&post)
	return &post
}

func (d *PostsDao) PostExists(id string) (bool, error) {
	var count int64
	err := d.db.Model(&model.Posts{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (d *PostsDao) GetPostsByUserId(userId string) (resp []model.Posts, err error) {
	err = d.db.Model(&model.Posts{}).Where("user_id = ?", userId).Find(&resp).Error
	return resp, nil
}
func (d *PostsDao) GetCountByUserId(userId string) (count int64, err error) {

	// 获取当前日期的开始和结束时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	tx := d.db.Model(&model.Posts{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userId, startOfDay, endOfDay).
		Count(&count)

	return count, tx.Error
}

func (d *PostsDao) Save(tx *gorm.DB, posts *model.Posts) error {
	var data = model.Posts{
		Id:         posts.Id,
		UserId:     posts.UserId,
		Content:    posts.Content,
		OriginalId: posts.OriginalId,
		SourceId:   posts.SourceId,
	}

	if tx == nil {
		tx = d.db
	}

	var err error
	err = tx.Create(&data).Error
	if err == nil && data.OriginalId != "" {
		err = d.UpdateRePostCount(data.OriginalId)
	}
	return err
}
func (d *PostsDao) Update(posts *model.Posts) error {
	tx := d.db.Updates(&posts)
	return tx.Error
}
func (d *PostsDao) UpdateCount(id string, count int64) error {
	err := d.db.Model(&model.Posts{}).Where("id = ?", id).Update("comments_amount", count).Error
	return err
}
func (d *PostsDao) UpdateRePostCount(id string) error {
	err := d.db.Model(&model.Posts{}).Where("id = ?", id).UpdateColumn("reposts_count", gorm.Expr("reposts_count + ?", 1)).Error
	return err
}
func (d *PostsDao) UpdateBatch(posts []model.Posts) error {
	// 开始事务
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, post := range posts {
			// 更新每个 post
			if err := tx.Model(&post).Updates(model.Posts{
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
func (d *PostsDao) Delete(tx *gorm.DB, id string) error {
	if tx == nil {
		tx = d.db
	}

	resp, err2 := d.GetPostsById(id)
	if err2 == nil && resp.Id != "" {
		err2 = tx.Where("id = ?", id).Delete(&model.Posts{}).Error
		d.DeleteCache(resp.UserId)
	}
	return err2
}

func (d *PostsDao) CalculateCommission(userId string) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("user:%s:posts:count", userId)

	count, err := d.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		count = 0
	}

	result, err := d.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 0 {
		now := time.Now()
		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		ttl := tomorrow.Sub(now)
		d.rdb.Expire(ctx, key, ttl)
	}
	points := d.calculatePoints(result)
	return points, nil
}

// 计算返佣
func (d *PostsDao) calculatePoints(value int64) int64 {
	switch value {
	case 1:
		return 20
	case 2:
		return 10
	default:
		return 0
	}
}

func (d *PostsDao) DeleteCache(userId string) error {
	ctx := context.Background()
	key := fmt.Sprintf("user:%s:posts:count", userId)

	_, err := d.rdb.Get(ctx, key).Int()
	if err != redis.Nil {
		d.rdb.Decr(ctx, key)
	}
	return nil

}

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
	query := d.db.Table("posts").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, user_name, wallet_public_key")
	}).Limit(p.Size).Offset((p.Page - 1) * p.Size)

	if p.Id != "" {
		query = query.Where("id = ?", p.Id)
	}

	if p.UserId != "" {
		query = query.Where("user_id = ?", p.UserId)
	}

	err = query.Order("created_at desc").Find(&resp).Error
	return resp, err
}

func (d *PostsDao) GetPostsById(id string) (resp model.Posts, err error) {
	err = d.db.Table("posts").Preload("User").Where("id = ?", id).First(&resp).Error
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

func (d *PostsDao) GetLikeByUser(p *model.PostsQuery) ([]model.Posts, error) {
	var posts []model.Posts
	err := d.db.Table("posts").Select("posts.*").Joins("RIGHT JOIN likes on posts.id = likes.post_id").Where("likes.user_id = ?", p.UserId).
		Limit(p.Size).Offset((p.Page - 1) * p.Size).Scan(&posts).Error
	return posts, err
}

func (d *PostsDao) GetCommentByUser(p *model.PostsQuery) ([]model.Posts, error) {
	var comments []model.Comments
	var posts []model.Posts

	err := d.db.Table("comments").Where("user_id = ?", p.UserId).Limit(p.Size).Offset((p.Page - 1) * p.Size).Group("post_id").
		Order("created_at desc").Scan(&comments).Error

	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return []model.Posts{}, nil
	}

	var postIds []string
	commentMap := make(map[string]model.Comments)
	for _, comment := range comments {
		postIds = append(postIds, comment.PostId)
		commentMap[comment.PostId] = comment
	}

	err = d.db.Table("posts").Where("id in (?)", postIds).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	for i := range posts {
		if comment, ok := commentMap[posts[i].Id]; ok {
			posts[i].Comments = &comment
		}
	}

	userIds := make([]string, 0)
	for _, post := range posts {
		userIds = append(userIds, post.UserId)
		userIds = append(userIds, post.Comments.UserId)
	}
	userIds = removeDuplicates(userIds)

	var users []model.Users
	err = d.db.Table("users").
		Where("id IN (?)", userIds).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	userMap := make(map[string]model.Users)
	for _, user := range users {
		userMap[user.Id] = user
	}

	// 6. 关联用户信息
	for i := range posts {
		if user, ok := userMap[posts[i].UserId]; ok {
			posts[i].User = &user
		}
		if user, ok := userMap[posts[i].Comments.UserId]; ok {
			posts[i].Comments.User = &user
		}
	}

	return posts, err
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
	err := d.db.Model(&model.Posts{}).Where("id = ?", id).UpdateColumn("comments_amount", count).Error
	return err
}
func (d *PostsDao) UpdateRePostCount(id string) error {
	err := d.db.Model(&model.Posts{}).Where("id = ?", id).UpdateColumn("reposts_amount", gorm.Expr("reposts_amount + ?", 1)).Error
	return err
}

func (d *PostsDao) UpdateLikeAmount(id string, amount int, types string) error {
	var err error
	if types == "add" {
		err = d.db.Model(&model.Posts{}).Where("id = ?", id).UpdateColumn("likes_amount", gorm.Expr("likes_amount + ?", amount)).Error
	} else {
		err = d.db.Model(&model.Posts{}).Where("id = ?", id).UpdateColumn("likes_amount", gorm.Expr("likes_amount - ?", amount)).Error
	}
	return err
}

func (d *PostsDao) UpdateBatch(posts []model.Posts) error {
	// 开始事务
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, post := range posts {
			// 更新每个 post
			if err := tx.Model(&post).Updates(model.Posts{
				Content:        post.Content,
				CommentsAmount: post.CommentsAmount,
				RepostsAmount:  post.RepostsAmount,
				//Tags:              post.Tags,
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
func removeDuplicates(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

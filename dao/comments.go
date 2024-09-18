package dao

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/logger"
	"gorm.io/gorm"
	"log"
	"novaro-server/model"
)

type CommentsDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewCommentsDao(db *gorm.DB) *CommentsDao {
	return &CommentsDao{
		db:  db,
		rdb: model.GetRedisCli(),
	}
}

func (d *CommentsDao) Create(c *model.Comments) (*model.Comments, int64, error) {
	err := d.db.Create(&c).Error

	var count int64
	if err == nil {
		key := fmt.Sprintf("post:%s:comment_count", c.PostId)
		count, err = d.rdb.Incr(context.Background(), key).Result()
	}
	return c, count, err
}

func (d *CommentsDao) GetById(id string) (resp *model.Comments, err error) {
	tx := d.db.Where("id = ?", id).First(&resp)
	return resp, tx.Error
}

func (d *CommentsDao) GetCount(postId string) int64 {
	var count int64
	d.db.Table("comments").Where("post_id = ?", postId).Count(&count)
	return count
}

func (d *CommentsDao) GetListByPostId(postId string) (resp []model.Comments, err error) {
	err = d.db.Table("comments").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,user_name,wallet_public_key")
	}).Where("post_id = ?", postId).Find(&resp).Error
	return resp, nil
}

func (d *CommentsDao) GetListByParentId(parentId string) (resp []model.Comments, err error) {

	err = d.db.Table("comments").Preload("User").Where("parent_id = ?", parentId).Find(&resp).Error
	if err != nil {
		return resp, err
	}

	for i := range resp {
		children, err := d.GetListByParentId(fmt.Sprint(resp[i].Id))
		if err != nil {
			return nil, err
		}
		resp[i].Children = children
	}
	return resp, nil
}

func (d *CommentsDao) GetListByUserId(userId string) (resp []model.Comments, err error) {
	err = d.db.Table("comments").Where("user_id = ?", userId).Find(&resp).Error
	return resp, nil
}

func (d *CommentsDao) DeleteById(id, postId string) (int64, error) {
	err := d.db.Table("comments").Where("id = ?", id).Delete(&model.Comments{}).Error
	fmt.Println(err)
	var count int64
	if err == nil {
		key := fmt.Sprintf("post:%s:comment_count", postId)
		count, err = d.rdb.Decr(context.Background(), key).Result()
	}
	return count, err
}

func (d *CommentsDao) SyncCommentsToDB() ([]string, *redis.Client) {
	ctx := context.Background()
	result, err := d.rdb.Keys(ctx, "post:*:comment_count").Result()
	if err != nil {
		logger.Error("redis get comments keys error", logger.Err(err))
		return nil, nil
	}
	return result, d.rdb
}

func (d *CommentsDao) GetCommentCount(postId string) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("post:%s:comment_count", postId)
	i, err := d.rdb.Get(ctx, key).Int64()
	if err == nil {
		return i, nil
	}

	err = d.db.Model(&model.Comments{}).Where("post_id = ?", postId).Count(&i).Error
	_, err = d.rdb.Set(ctx, key, i, 0).Result()
	if err != nil {
		// 这里我们只记录错误，不返回，因为我们已经有了正确的计数
		log.Printf("Failed to set Redis cache: %v", err)
	}
	return i, err
}

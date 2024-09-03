package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type CollectionsDao struct {
	db *gorm.DB
}

func NewCollectionsDao(db *gorm.DB) *CollectionsDao {
	return &CollectionsDao{
		db: db,
	}
}

func (d *CollectionsDao) CollectionsExist(userId, postId string) (bool, error) {
	var count int64
	err := d.db.Model(&model.Collections{}).Where("user_id = ? and post_id = ?", userId, postId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *CollectionsDao) Creates(c *model.Collections) error {
	tx := d.db.Create(c)

	return tx.Error
}

func (d *CollectionsDao) DeleteByUserIdAndPostId(userId, postId string) error {
	tx := d.db.Where("user_id = ? and post_id = ?", userId, postId).Delete(&model.Collections{})
	return tx.Error
}

// 刷新数据库
func (d *CollectionsDao) RefreshData(operations []model.Queue) error {
	// 开始事务
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, coll := range operations {
			if coll.Operation == "add" {
				err := d.Creates(&coll.Collections)
				return err
			} else {
				err := d.DeleteByUserIdAndPostId(coll.Collections.UserId, coll.Collections.PostId)
				return err
			}
		}
		return nil
	})
	return err
}

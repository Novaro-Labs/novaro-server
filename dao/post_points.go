package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type PostPointsDao struct {
	db *gorm.DB
}

func NewPostPointsDao(db *gorm.DB) *PostPointsDao {
	return &PostPointsDao{
		db: db,
	}
}

func (d *PostPointsDao) Save(tx *gorm.DB, m *model.PostPoints) error {
	if tx == nil {
		tx = d.db
	}
	return tx.Create(&m).Error
}

func (d *PostPointsDao) Delete(tx *gorm.DB, postId string) error {
	if tx == nil {
		tx = d.db
	}
	return tx.Where("post_id = ?", postId).Delete(&model.PostPoints{}).Error
}

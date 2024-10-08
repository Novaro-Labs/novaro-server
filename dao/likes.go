package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type LikeDao struct {
	dao *gorm.DB
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{
		dao: db,
	}
}

func (d *LikeDao) Add(like *model.Likes) error {
	err := d.dao.Create(&like).Error
	return err
}

func (d *LikeDao) Delete(like *model.Likes) error {
	err := d.dao.Model(&model.Likes{}).Where("id = ?", like.Id).Delete(&like).Error
	return err
}

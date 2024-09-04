package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type PointsHistoryDao struct {
	db *gorm.DB
}

func NewPointsHistoryDao(db *gorm.DB) *PointsHistoryDao {
	return &PointsHistoryDao{
		db: db,
	}
}

func (d *PointsHistoryDao) Create(tx *gorm.DB, history *model.PointsHistory) error {
	if tx == nil {
		tx = d.db
	}
	return tx.Create(&history).Error
}

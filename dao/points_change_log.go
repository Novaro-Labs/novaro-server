package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type PointsChangeLogDao struct {
	db *gorm.DB
}

func NewPointsChangeLogDao(db *gorm.DB) *PointsChangeLogDao {
	return &PointsChangeLogDao{
		db: db,
	}
}

func (d *PointsChangeLogDao) GetList(log *model.PointsChangeLogRequest) ([]model.PointsChangeLog, error) {
	var logs []model.PointsChangeLog
	err := d.db.Model(&model.PointsChangeLog{}).Limit(log.Size).Offset((log.Page - 1) * log.Size).Find(&logs).Error
	return logs, err
}

func (d *PointsChangeLogDao) Create(tx *gorm.DB, log *model.PointsChangeLog) error {
	if tx == nil {
		tx = d.db
	}
	err := tx.Create(&log).Error
	return err
}

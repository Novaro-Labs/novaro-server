package dao

import (
	"fmt"
	"gorm.io/gorm"
	"novaro-server/model"
	"time"
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
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour)
	fmt.Println(thirtyDaysAgo)
	err := d.db.Model(&model.PointsChangeLog{}).Limit(log.Size).Offset((log.Page-1)*log.Size).Order("created_at desc").
		Where("wallet = ? and created_at >? ", log.Wallet, thirtyDaysAgo).Find(&logs)
	return logs, err.Error
}

func (d *PointsChangeLogDao) GetYesterdayPoints(wallet string) (int, error) {
	var points int

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour).Add(-time.Nanosecond)

	err := d.db.Model(&model.PointsChangeLog{}).Where("wallet = ? and created_at >= ? and created_at <= ?", wallet, todayStart, todayEnd).Select("sum(change_amount)").Scan(&points).Error
	return points, err
}

func (d *PointsChangeLogDao) Create(tx *gorm.DB, log *model.PointsChangeLog) error {
	if tx == nil {
		tx = d.db
	}
	err := tx.Create(&log).Error
	return err
}

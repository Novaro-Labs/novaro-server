package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
	"time"
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

func (d *PostPointsDao) BatchSave(points []model.PostPoints) bool {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, point := range points {
			err := d.Save(tx, &point)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return false
	}
	return true
}

func (d *PostPointsDao) GetYesterdayPostHistory() ([]model.PostPoints, error) {
	// 获取当前时间
	now := time.Now()

	// 计算昨天的日期（不考虑时分秒）
	yesterday := now.AddDate(0, 0, -1)

	// 设置昨天的开始时间（00:00:00）
	yesterdayStart := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())

	// 设置昨天的结束时间（23:59:59.999999999）
	yesterdayEnd := yesterdayStart.AddDate(0, 0, 1).Add(-time.Nanosecond)

	var records []model.PostPoints
	err := d.db.Where("created_at >= ? AND created_at <= ?", yesterdayStart, yesterdayEnd).Find(&records).Error
	return records, err
}

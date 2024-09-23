package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
	"time"
)

type PointsHistoryDao struct {
	db *gorm.DB
}

func NewPointsHistoryDao(db *gorm.DB) *PointsHistoryDao {
	return &PointsHistoryDao{
		db: db,
	}
}

func (d *PointsHistoryDao) GetList(p *model.PointsHistoryQuery) ([]model.PointsHistory, error) {
	var points []model.PointsHistory
	err := d.db.Model(&model.PointsHistory{}).Select("id,points").Where("wallet = ? and status = '0'", p.Wallet).Find(&points).Error
	return points, err
}

func (d *PointsHistoryDao) Statistics(wallet string, datetime *time.Time) ([]model.PointsHistory, error) {
	var results []model.PointsHistory

	// 如果 datetime 为 nil，使用当前时间
	endDate := time.Now()
	if datetime != nil {
		endDate = *datetime
	}

	startDate := endDate.AddDate(0, 0, -7)

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location()).Add(-time.Nanosecond)

	err := d.db.Model(&model.PointsHistory{}).
		Select("DATE(create_at) as date, SUM(points) as points").
		Where("wallet = ? AND create_at >= ? AND create_at < ?", wallet, startDate, endDate).
		Group("DATE(create_at)").
		Order("date ASC").
		Scan(&results).Error

	return results, err
}

func (d *PointsHistoryDao) Create(tx *gorm.DB, history *model.PointsHistory) error {
	if tx == nil {
		tx = d.db
	}
	return tx.Create(&history).Error
}

func (d *PointsHistoryDao) BatchSave(history []model.PointsHistory) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, h := range history {
			err := d.Create(tx, &h)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (d *PointsHistoryDao) UpdateStatus(tx *gorm.DB, id string) error {
	err := tx.Model(&model.PointsHistory{}).Where("id = ?", id).UpdateColumn("status", "1").Error
	return err

}

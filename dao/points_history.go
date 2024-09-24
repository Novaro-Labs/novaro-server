package dao

import (
	"fmt"
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

func (d *PointsHistoryDao) GetList(wallet string) ([]model.PointsHistory, error) {
	var points []model.PointsHistory
	err := d.db.Model(&model.PointsHistory{}).Select("id,points").Where("wallet = ? and status = '0'", wallet).Find(&points).Error
	return points, err
}

func (d *PointsHistoryDao) Statistics(wallet string, datetime *time.Time) ([]model.PointsHistoryStatistics, error) {

	var dbResults []struct {
		Date   string  `gorm:"column:date"`
		Points float64 `json:"points"`
	}

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
		Scan(&dbResults).Error

	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

	// 创建一个map来存储数据库结果，方便查找
	dataMap := make(map[string]float64)
	for _, r := range dbResults {
		t, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		dataMap[weekdays[t.Weekday()]] = r.Points
	}

	// 创建完整的一周数据
	results := make([]model.PointsHistoryStatistics, 7)
	for i, day := range weekdays {
		points, exists := dataMap[day]
		if !exists {
			points = 0 // 如果没有数据，设置为0
		}
		results[i] = model.PointsHistoryStatistics{
			Date:   day,
			Points: points,
		}
	}

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

func (d *PointsHistoryDao) UpdateStatus(tx *gorm.DB, id string) (points float64, err error) {
	var pointsHistory model.PointsHistory
	err = tx.Model(&model.PointsHistory{}).Where("id = ?", id).UpdateColumn("status", "1").Error

	err = tx.Where("id = ?", id).First(&pointsHistory).Error
	return pointsHistory.Points, err

}

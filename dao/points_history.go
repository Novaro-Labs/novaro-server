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
	err := d.db.Model(&model.PointsHistory{}).Select("id,points").Limit(p.Size).
		Offset((p.Page-1)*p.Size).Where("wallet = ? and status = '0'", p.Wallet).Find(&points).Error
	return points, err
}

func (d *PointsHistoryDao) Statistics(wallet string, datetime time.Time) {
	d.db.Model(&model.PointsHistory{}).Select("sum(points) as points").Where("wallet = ? and status = '0' and create_at > ?", wallet, datetime).
		Scan(&model.PointsHistory{})
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

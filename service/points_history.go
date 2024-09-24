package service

import (
	"gorm.io/gorm"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type PointsHistoryService struct {
	dao *dao.PointsHistoryDao
}

func NewPointsHistoryService() *PointsHistoryService {
	return &PointsHistoryService{
		dao: dao.NewPointsHistoryDao(model.GetDB()),
	}
}

func (s *PointsHistoryService) Create(tx *gorm.DB, history *model.PointsHistory) error {
	return s.dao.Create(tx, history)
}

func (s *PointsHistoryService) GetList(wallet string) ([]model.PointsHistory, error) {
	return s.dao.GetList(wallet)
}

func (s *PointsHistoryService) Statistics(wallet string, datetime *time.Time) ([]model.PointsHistoryStatistics, error) {

	return s.dao.Statistics(wallet, datetime)

}

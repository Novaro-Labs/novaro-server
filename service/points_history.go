package service

import (
	"gorm.io/gorm"
	"novaro-server/dao"
	"novaro-server/model"
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

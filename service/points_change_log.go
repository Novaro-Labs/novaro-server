package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type PointsChangeLogService struct {
	dao *dao.PointsChangeLogDao
}

func NewPointsChangeLogService() *PointsChangeLogService {
	return &PointsChangeLogService{
		dao: dao.NewPointsChangeLogDao(model.GetDB()),
	}
}

func (s *PointsChangeLogService) GetList(r *model.PointsChangeLogRequest) ([]model.PointsChangeLog, error) {
	return s.dao.GetList(r)
}

package service

import (
	"gorm.io/gorm"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type NftInfoService struct {
	dao                *dao.NftInfoDao
	pointsHistoryDao   *dao.PointsHistoryDao
	pointsChangeLogDao *dao.PointsChangeLogDao
}

func NewNftInfoService() *NftInfoService {
	db := model.GetDB()
	return &NftInfoService{
		dao:                dao.NewNftInfoDao(db),
		pointsHistoryDao:   dao.NewPointsHistoryDao(db),
		pointsChangeLogDao: dao.NewPointsChangeLogDao(db),
	}
}

func (s *NftInfoService) Create(info *model.NftInfo) (string, error) {
	return s.dao.Create(info)
}

func (s *NftInfoService) UpdateScore(wallet string, score float64) error {
	return s.dao.UpdateScore(wallet, score)
}

func (s *NftInfoService) GetByWallet(wallet string) (model.NftInfo, error) {
	return s.dao.GetByWallet(wallet)
}

func (s *NftInfoService) UpdatePoints(info *model.NftInfoRequest) (points float64, err error) {

	err = model.GetDB().Transaction(func(tx *gorm.DB) error {
		points, err = s.dao.UpdatePoints(tx, info)
		if err != nil {
			return err
		}

		err = s.pointsHistoryDao.UpdateStatus(tx, info.PointId)
		if err != nil {
			return err
		}

		err = s.pointsChangeLogDao.Create(tx, &model.PointsChangeLog{
			Wallet:       info.Wallet,
			ChangeAmount: info.Points,
			ChangeType:   0,
			Reason:       "",
			CreatedAt:    time.Now(),
		})
		if err != nil {
			return err
		}

		return nil
	})
	return points, err
}

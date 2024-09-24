package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type NftInfoService struct {
	dao                *dao.NftInfoDao
	pointsHistoryDao   *dao.PointsHistoryDao
	pointsChangeLogDao *dao.PointsChangeLogDao
	rdb                *redis.Client
}

func NewNftInfoService() *NftInfoService {
	db := model.GetDB()
	return &NftInfoService{
		dao:                dao.NewNftInfoDao(db),
		pointsHistoryDao:   dao.NewPointsHistoryDao(db),
		pointsChangeLogDao: dao.NewPointsChangeLogDao(db),
		rdb:                model.GetRedisCli(),
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

func (s *NftInfoService) UpdatePoints(info *model.NftInfoRequest) (map[string]any, error) {
	maps := make(map[string]any, 2)
	var err error
	var points float64

	err = model.GetDB().Transaction(func(tx *gorm.DB) error {
		point, err2 := s.pointsHistoryDao.UpdateStatus(tx, info.PointId)
		if err2 != nil {
			return err2
		}
		info.Points = point

		points, err = s.dao.UpdatePoints(tx, info)
		if err != nil {
			return err
		}

		err = s.pointsChangeLogDao.Create(tx, &model.PointsChangeLog{
			Wallet:       info.Wallet,
			ChangeAmount: point,
			ChangeType:   0,
			Reason:       "",
			CreatedAt:    time.Now(),
		})

		if err != nil {
			return err
		}

		return nil
	})

	maps["totalPoints"] = points
	yesterdayPoints, err := s.pointsChangeLogDao.GetYesterdayPoints(info.Wallet)
	if err != nil {
		yesterdayPoints = 0
	}
	maps["currentPoints"] = yesterdayPoints
	return maps, err
}

//func (s *NftInfoService) CheckUpgrade(wallet string, points float64) (bool, error) {
//	//s.rdb.ZScore()
//

//}

package service

import (
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
)

type NftInfoService struct {
	dao *dao.NftInfoDao
}

func NewNftInfoService() *NftInfoService {
	return &NftInfoService{
		dao: dao.NewNftInfoDao(config.DB),
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

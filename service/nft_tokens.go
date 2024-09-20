package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type NftTokensService struct {
	dao *dao.NftTokensDao
}

func NewNftTokensService() *NftTokensService {
	return &NftTokensService{
		dao: dao.NewNftTokensDao(model.GetDB()),
	}
}

func (s *NftTokensService) GetTokensByWallet(wallet string) ([]model.NftTokens, error) {
	return s.dao.GetTokensByWallet(wallet)
}

func (s *NftTokensService) SaveNftToken(nftToken *model.NftTokens) error {
	return s.dao.SaveNftToken(nftToken)
}

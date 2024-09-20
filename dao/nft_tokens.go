package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type NftTokensDao struct {
	db *gorm.DB
}

func NewNftTokensDao(db *gorm.DB) *NftTokensDao {
	return &NftTokensDao{
		db: db,
	}
}

func (d *NftTokensDao) GetTokensByWallet(wallet string) ([]model.NftTokens, error) {
	var nftTokens []model.NftTokens
	err := d.db.Model(&model.NftTokens{}).Preload("Img").Where("wallet = ?", wallet).Find(&nftTokens).Error
	return nftTokens, err
}

func (d *NftTokensDao) SaveNftToken(nftToken *model.NftTokens) error {
	return d.db.Save(nftToken).Error
}

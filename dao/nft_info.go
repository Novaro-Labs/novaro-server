package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type NftInfoDao struct {
	db *gorm.DB
}

func NewNftInfoDao(db *gorm.DB) *NftInfoDao {
	return &NftInfoDao{
		db: db,
	}
}

func (d *NftInfoDao) Create(info *model.NftInfo) (string, error) {
	tx := d.db.Create(&info)
	return info.Id, tx.Error
}

func (d *NftInfoDao) UpdateScore(wallet string, score float64) error {
	tx := d.db.Where("wallet = ? ", wallet, score).UpdateColumn("score", score)
	return tx.Error
}

func (d *NftInfoDao) GetByWallet(wallet string) (model.NftInfo, error) {
	var resp model.NftInfo
	err := d.db.Table("nft_info").Where("wallet = ?", wallet).Find(&resp).Error
	return resp, err

}

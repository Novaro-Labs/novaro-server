package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err := d.db.Model(&model.NftInfo{}).Preload("Nft", func(db *gorm.DB) *gorm.DB {
		return db.Select("Id, Right,Url ")
	}).Where("wallet = ?", wallet).First(&resp).Error
	return resp, err
}

func (d *NftInfoDao) UpdatePoints(tx *gorm.DB, info *model.NftInfoRequest) (float64, error) {
	var nftInfo model.NftInfo

	if tx == nil {
		tx = d.db
	}

	result := tx.Model(&model.NftInfo{}).Clauses(clause.Returning{Columns: []clause.Column{{Name: "points"}}}).Where("wallet = ?", info.Wallet).
		UpdateColumn("points", gorm.Expr("points+?", info.Points)).
		Scan(&nftInfo)

	if result.Error != nil {
		return 0, result.Error
	}
	return nftInfo.Points, nil
}

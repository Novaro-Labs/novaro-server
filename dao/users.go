package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type UsersDao struct {
	db *gorm.DB
}

func NewUsersDao(db *gorm.DB) *UsersDao {
	return &UsersDao{
		db: db,
	}
}

func (d *UsersDao) SaveUsers(users *model.Users) (string, error) {
	var data = model.Users{
		TwitterId:       users.TwitterId,
		UserName:        users.UserName,
		CreatedAt:       users.CreatedAt,
		Avatar:          users.Avatar,
		WalletPublicKey: users.WalletPublicKey,
	}

	tx := d.db.Table("users").Where("twitter_id = ?", users.TwitterId).FirstOrCreate(&data)

	return data.Id, tx.Error
}

func (d *UsersDao) UserExists(userId string) (bool, error) {
	var count int64
	err := d.db.Model(&model.Users{}).Where("id = ?", userId).Count(&count).Error
	return count > 0, err
}

func (d *UsersDao) GetById(id string) (model.Users, error) {
	var user model.Users
	tx := d.db.Model(&model.Users{}).Joins("left join nft_info on nft_info.wallet = users.wallet_public_key").Where("users.id = ?", id).First(&user)
	return user, tx.Error
}

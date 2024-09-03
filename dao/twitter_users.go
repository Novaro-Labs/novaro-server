package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type TwitterUserDao struct {
	db *gorm.DB
}

func NewTwitterUserDao(db *gorm.DB) *TwitterUserDao {
	return &TwitterUserDao{
		db: db,
	}
}

func (d *TwitterUserDao) SaveTwitterUsers(users *model.TwitterUsers) error {

	tx := d.db.Table("twitter_user").Where("twitter_id = ?", users.TwitterId).FirstOrCreate(&users)
	return tx.Error
}

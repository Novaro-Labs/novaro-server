package model

import (
	"novaro-server/config"
	"time"
)

type TwitterUsers struct {
	TwitterId        string     `json:"twitterId"`
	TwitterUserName  string     `json:"twitterUserName"`
	TwitterAvatar    *string    `json:"twitterAvatar"`
	TwitterFollowers *int       `json:"twitterFollowers"`
	TwitterCreatedAt *time.Time `json:"twitterCreatedAt"`
}

func (TwitterUsers) TableName() string {
	return "twitter_user"
}

func SaveTwitterUsers(users *TwitterUsers) error {
	db := config.DB

	var data = TwitterUsers{
		TwitterId:        users.TwitterId,
		TwitterUserName:  users.TwitterUserName,
		TwitterAvatar:    users.TwitterAvatar,
		TwitterFollowers: users.TwitterFollowers,
		TwitterCreatedAt: users.TwitterCreatedAt,
	}
	tx := db.Table("twitter_user").Where("twitter_id = ?", users.TwitterId).FirstOrCreate(&data)
	return tx.Error
}

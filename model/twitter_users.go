package model

import (
	"time"
)

type TwitterUsers struct {
	TwitterId        string     `json:"twitterId"`
	TwitterUserName  string     `json:"twitterUserName"`
	TwitterAvatar    *string    `json:"twitterAvatar"`
	TwitterFollowers *int       `json:"twitterFollowers"`
	Wallet           string     `json:"wallet"`
	TwitterCreatedAt *time.Time `json:"twitterCreatedAt"`
}

func (TwitterUsers) TableName() string {
	return "twitter_user"
}

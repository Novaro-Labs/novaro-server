package model

import (
	"time"
)

type TwitterUser struct {
	TwitterId        string     `json:"twitterId"`
	TwitterUserName  string     `json:"twitterUserName"`
	TwitterAvatar    *string    `json:"twitterAvatar"`
	TwitterFollowers *int       `json:"twitterFollowers"`
	TwitterCreatedAt *time.Time `json:"twitterCreatedAt"`
}

func (TwitterUser) TableName() string {
	return "twitter_user"
}

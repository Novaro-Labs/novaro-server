package model

import "time"

type TwitterUserInfo struct {
	Id       string    `json:"id"`
	Avatar   string    `json:"profile_image_url"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created_at"`
	Username string    `json:"username"`
}

func (userInfo *TwitterUserInfo) ToUsers() *Users {
	users := &Users{
		TwitterId: userInfo.Id,
		UserName:  userInfo.Username,
		CreatedAt: userInfo.Created,
		Avatar:    &userInfo.Avatar,
	}
	return users
}

func (userInfo *TwitterUserInfo) ToTwitterUsers() *TwitterUsers {
	twitterUsers := &TwitterUsers{
		TwitterId:        userInfo.Id,
		TwitterUserName:  userInfo.Username,
		TwitterAvatar:    &userInfo.Avatar,
		TwitterCreatedAt: &userInfo.Created,
	}
	return twitterUsers
}

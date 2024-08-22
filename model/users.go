package model

import (
	"novaro-server/config"
	"time"
)

type Users struct {
	Id              string    `json:"id"`
	TwitterId       string    `json:"twitterID"`
	UserName        string    `json:"userName"`
	Avatar          *string   `json:"avatar"`
	Followers       uint      `json:"followers"`
	Following       uint      `json:"following"`
	WalletPublicKey *string   `json:"walletPublicKey"`
	InvitationCode  *string   `json:"invitationCode"`
	UserLevel       int       `json:"userLevel"`
	UserScore       float64   `json:"userScore"`
	CreditScore     float64   `json:"creditScore"`
	CreatedAt       time.Time `json:"createdAt"`
	LastLogin       time.Time `json:"lastLogin"`
}

func (Users) TableName() string {
	return "users"
}

func SaveUsers(users *Users) error {
	db := config.DB
	var data = Users{
		Id:              users.Id,
		TwitterId:       users.TwitterId,
		UserName:        users.UserName,
		CreatedAt:       users.CreatedAt,
		Avatar:          users.Avatar,
		WalletPublicKey: users.WalletPublicKey,
	}

	tx := db.Create(&data)
	return tx.Error
}

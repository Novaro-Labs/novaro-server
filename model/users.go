package model

import (
	"novaro-server/config"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (u *Users) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func SaveUsers(users *Users) (string, error) {
	db := config.DB

	var data = Users{
		TwitterId:       users.TwitterId,
		UserName:        users.UserName,
		CreatedAt:       users.CreatedAt,
		Avatar:          users.Avatar,
		WalletPublicKey: users.WalletPublicKey,
	}

	tx := db.Table("users").Where("twitter_id = ?", users.TwitterId).FirstOrCreate(&data)

	return data.Id, tx.Error
}

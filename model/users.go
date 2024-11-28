package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Users struct {
	Id              string     `json:"id"`
	TwitterId       string     `json:"twitterID,omitempty"`
	UserName        string     `json:"userName,omitempty"`
	Avatar          *string    `json:"avatar,omitempty"`
	Followers       uint       `json:"followers,omitempty"`
	Following       uint       `json:"following,omitempty"`
	WalletPublicKey string     `json:"walletPublicKey,omitempty"`
	InvitationCode  *string    `json:"invitationCode,omitempty"`
	UserLevel       int        `json:"userLevel,omitempty"`
	UserScore       float64    `json:"userScore,omitempty"`
	CreditScore     float64    `json:"creditScore,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	LastLogin       *time.Time `json:"lastLogin,omitempty"`
	NftInfo         *NftInfo   `json:"nftInfo,omitempty" gorm:"foreignKey:Wallet;references:WalletPublicKey"`
}

func (Users) TableName() string {
	return "users"
}

func (u *Users) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

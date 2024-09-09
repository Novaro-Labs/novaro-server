package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Users struct {
	Id              string    `json:"id" `
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
	NftInfo         *NftInfo  `json:"nftInfo" gorm:"foreignKey:Wallet;references:WalletPublicKey"`
}

func (Users) TableName() string {
	return "users"
}

func (u *Users) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func (u *Users) Points(wattle *string, nftLevel int) int {
	if wattle == nil {
		return 0
	}
	defaultPoints := 5
	rewards := nftLevel

	return (nftLevel * defaultPoints) + rewards
}

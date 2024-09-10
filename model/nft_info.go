package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type NftInfo struct {
	Id        string    `json:"id"`
	Wallet    string    `json:"wallet"`
	Level     int       `json:"level"`
	points    float64   `json:"points"`
	Create_at time.Time `json:"createAt"`
}

func (NftInfo) TableName() string {
	return "nft_info"
}
func (u *NftInfo) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

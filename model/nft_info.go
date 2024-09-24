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
	Points    float64   `json:"points"`
	Nft       *NftLevel `json:"nft" gorm:"foreignKey:id;references:Level" `
	Create_at time.Time `json:"createAt"`
}
type NftInfoRequest struct {
	PointId string  `json:"pointId" binding:"required"`
	Wallet  string  `json:"wallet" binding:"required"`
	Points  float64 `json:"points"`
}

func (NftInfo) TableName() string {
	return "nft_info"
}
func (u *NftInfo) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

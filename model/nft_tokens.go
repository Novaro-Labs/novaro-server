package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type NftTokens struct {
	Id          string     `json:"id" `
	Wallet      string     `json:"wallet" binding:"required"`
	TokenName   string     `json:"tokenName" binding:"required"`
	TokenSymbol string     `json:"tokenSymbol" binding:"required"`
	SourceId    string     `json:"sourceId" binding:"required"`
	Img         *Imgs      `json:"img" gorm:"foreignKey:SourceId;references:SourceId;"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func (NftTokens) TableName() string {
	return "nft_tokens"
}

func (c *NftTokens) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	c.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

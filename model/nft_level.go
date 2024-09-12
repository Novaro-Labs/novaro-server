package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type NftLevel struct {
	Id        string     `json:"level,omitempty"`
	Left      int        `json:"left,omitempty"`
	Right     int        `json:"right,omitempty"`
	Url       string     `json:"url,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

func (NftLevel) TableName() string {
	return "nft_level"
}

func (u *NftLevel) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

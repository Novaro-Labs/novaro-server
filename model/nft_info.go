package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type NftInfo struct {
	Id        string    `json:"id"`
	Wallet    string    `json:"wallet"`
	Level     int       `json:"level"`
	Score     float64   `json:"score"`
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

func (u *NftInfo) TagPoints(count int64) float64 {
	totalPoints := float64(1)

	if count > 0 {
		count = count / 10
	}

	coefficients := float64(count)

	rewards := (totalPoints - coefficients) * totalPoints
	return math.Max(0, rewards)
}

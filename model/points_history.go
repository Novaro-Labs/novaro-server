package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PointsHistory struct {
	Id       string    `json:"id"`
	Wallet   *string   `json:"wallet"`
	Points   float64   `json:"points"`
	Status   int       `json:"status"`
	CreateAt time.Time `json:"createAt"`
}

func (PointsHistory) TableName() string {
	return "points_history"
}
func (u *PointsHistory) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

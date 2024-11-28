package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PointsHistory struct {
	Id       string     `json:"id,omitempty"`
	Wallet   string     `json:"wallet,omitempty"`
	Points   float64    `json:"points,omitempty"`
	Status   int        `json:"status,omitempty"`
	CreateAt *time.Time `json:"createAt,omitempty"`
}

type PointsHistoryQuery struct {
	Wallet string `json:"wallet"  binding:"required"`
	Page   int    `json:"page"  binding:"required"`
	Size   int    `json:"size"  binding:"required"`
}

type PointsHistoryStatistics struct {
	Date   string  `json:"date"`
	Points float64 `json:"points"`
}

func (PointsHistory) TableName() string {
	return "points_history"
}
func (u *PointsHistory) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

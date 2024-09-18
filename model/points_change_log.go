package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PointsChangeLog struct {
	Id           string    `json:"id"`
	Wallet       string    `json:"wallet"`
	ChangeAmount float64   `json:"changeAmount"`
	ChangeType   int       `json:"changeType"`
	Reason       string    `json:"reason"`
	CreatedAt    time.Time `json:"createdAt"`
}

type PointsChangeLogRequest struct {
	Wallet string `json:"wallet" binding:"required"`
	Page   int    `json:"page" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}

func (PointsChangeLog) TableName() string {
	return "points_change_log"
}
func (u *PointsChangeLog) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

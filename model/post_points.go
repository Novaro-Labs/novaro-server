package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PostPoints struct {
	Id        string    `json:"id"`
	PostId    string    `json:"postId"`
	Points    float64   `json:"points"`
	CreatedAt time.Time `json:"createdAt"`
}

func (PostPoints) TableName() string {
	return "post_points"
}

func (u *PostPoints) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

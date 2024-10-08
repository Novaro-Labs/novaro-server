package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Likes struct {
	Id      string    `json:"id"`
	UserId  string    `json:"userId"`
	PostId  string    `json:"postId"`
	Created time.Time `json:"created"`
}

func (Likes) TableName() string {
	return "likes"
}
func (u *Likes) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Comments struct {
	Id        string     `json:"id"`
	UserId    string     `json:"userId"`
	PostId    string     `json:"postId"`
	ParentId  string     `json:"parentId"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	User      *Users     `json:"user" gorm:"foreignKey:id;references:UserId;"`
	Children  []Comments `json:"children,omitempty" gorm:"-"`
}

func (Comments) TableName() string {
	return "comments"
}

func (u *Comments) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

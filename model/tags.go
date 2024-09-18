package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Tags struct {
	Id        string     `json:"id"`
	SourceId  string     `json:"sourceId"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	Posts     []Posts    `json:"posts,omitempty" gorm:"-"`
}

func (u *Tags) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func (Tags) TableName() string {
	return "tags"
}

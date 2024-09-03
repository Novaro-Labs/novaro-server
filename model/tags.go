package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Tags struct {
	Id        string    `json:"id"`
	TagType   string    `json:"tagType"`
	TagColor  string    `json:"tagColor"`
	CreatedAt time.Time `json:"createdAt"`
	Posts     []Posts   `json:"posts" gorm:"-"`
}

func (u *Tags) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func (Tags) TableName() string {
	return "tags"
}

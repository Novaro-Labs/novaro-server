package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TagsRecords struct {
	Id        string    `json:"id"`
	TagId     string    `json:"tagId"`
	PostId    string    `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *TagsRecords) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func (TagsRecords) TableName() string {
	return "tags_records"
}

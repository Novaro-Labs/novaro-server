package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

// 收藏，先记录在redis中，每五分钟更新一次数据库

type Collections struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

type Queue struct {
	Collections Collections
	Operation   string // add: 收藏 remove: 取消收藏
}

func (Collections) TableName() string {
	return "collections"
}

func (c *Collections) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	c.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

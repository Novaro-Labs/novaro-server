package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Posts struct {
	Id                string    `json:"id"`
	UserId            string    `json:"userId"`
	Content           string    `json:"content"`
	CommentsAmount    int       `json:"commentsAmount"`
	CollectionsAmount int       `json:"collectionsAmount"`
	RepostsAmount     int       `json:"repostsAmount"`
	CreatedAt         time.Time `json:"createdAt"`
	OriginalId        string    `json:"originalId"`
	SourceId          string    `json:"sourceId"`
	Tags              []Tags    `json:"tags" gorm:"-"`
	Imgs              []Imgs    `json:"Imgs" gorm:"-"`
	IsCollected       bool      `json:"isCollected" gorm:"-"`
}

func (Posts) TableName() string {
	return "posts"
}

type PostsQuery struct {
	Id     string `form:"id" json:"id"`
	UserId string `form:"userId" json:"userId"`
}

func (u *Posts) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

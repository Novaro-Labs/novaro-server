package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Posts struct {
	Id                string              `json:"id"`
	UserId            string              `json:"userId"`
	Content           string              `json:"content"`
	CommentsAmount    int64               `json:"commentsAmount"`
	CollectionsAmount int                 `json:"collectionsAmount"`
	RepostsAmount     int                 `json:"repostsAmount"`
	TagsAmount        int                 `json:"tagsAmount"`
	CreatedAt         time.Time           `json:"createdAt"`
	OriginalId        string              `json:"originalId"`
	ViewAmount        int                 `json:"viewAmount"`
	SourceId          string              `json:"sourceId"`
	User              *Users              `json:"user" gorm:"foreignKey:id;references:UserId"`
	Tags              []TagRecordResponse `json:"tags" gorm:"-"`
	Imgs              []Imgs              `json:"Imgs" gorm:"-"`
}

func (Posts) TableName() string {
	return "posts"
}

type PostsQuery struct {
	Id     string `form:"id" json:"id"`
	UserId string `form:"userId" json:"userId"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

func (u *Posts) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

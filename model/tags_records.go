package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type TagsRecords struct {
	Id         string    `json:"id"`
	TagId      string    `json:"tagId"`
	PostId     string    `json:"postId"`
	UserId     string    `json:"userId"`
	Points     float64   `json:"points"`
	PostPoints float64   `json:"postPoints"`
	CreatedAt  time.Time `json:"createdAt"`
}

type TagRecordQueue struct {
	TagId      string    `json:"tagId"`
	PostId     string    `json:"postId"`
	UserId     string    `json:"userId"`
	Points     float64   `json:"points"`
	PostPoints float64   `json:"postPoints"`
	CreatedAt  time.Time `json:"createdAt"`
	Operation  string    `json:"operation"`
}

type TagRecordResponse struct {
	Id    string `json:"id"`
	Count int    `json:"count"`
}

func (u *TagsRecords) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func (TagsRecords) TableName() string {
	return "tags_records"
}

func (u *TagsRecords) TagPoints(wattle *string, count int64) float64 {
	if wattle == nil {
		return 0
	}

	totalPoints := float64(1)

	if count > 0 {
		count = count / 10
	}

	coefficients := float64(count)
	rewards := (totalPoints - coefficients) * totalPoints
	return math.Max(0, rewards)
}

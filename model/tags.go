package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"novaro-server/config"
	"strings"
	"time"
)

type Tags struct {
	Id        string    `json:"id"`
	TagType   string    `json:"tagType"`
	TagColor  string    `json:"tagColor"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *Tags) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func (Tags) TableName() string {
	return "tags"
}

func GetTagsList() (resp []Tags, err error) {
	err = config.DB.Model(Tags{}).Find(&resp).Error
	return resp, err
}

func GetTagListByPostId(postId string) (resp []Tags, err error) {
	err = config.DB.Distinct("tags.*").Model(&Tags{}).
		Joins("JOIN tags_records ON tags.id = tags_records.tag_id").
		Where("tags_records.post_id = ?", postId).
		Find(&resp).Error
	return resp, err
}

func TagExists(id string) (bool, error) {
	var count int64
	tx := config.DB.Model(&Tags{}).Where("id = ?", id).Count(&count)
	return count > 0, tx.Error
}

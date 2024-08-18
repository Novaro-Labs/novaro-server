package model

import (
	"novaro-server/config"
	"time"
)

type Tags struct {
	Id        string    `json:"id"`
	TagType   string    `json:"tagType"`
	TagColor  string    `json:"tagColor"`
	CreatedAt time.Time `json:"createdAt"`
}

func GetTagsList() (resp []Tags, err error) {
	db := config.DB
	err = db.Model(Tags{}).Find(&resp).Error
	return resp, err
}

func GetTagListByPostId(postId string) (resp []Tags, err error) {
	db := config.DB
	err = db.Distinct("tags.*").Model(&Tags{}).
		Joins("JOIN tags_records ON tags.id = tags_records.tag_id").
		Where("tags_records.post_id = ?", postId).
		Find(&resp).Error
	return resp, err
}

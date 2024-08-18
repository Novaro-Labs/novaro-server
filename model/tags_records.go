package model

import (
	"novaro-server/config"
	"time"
)

type TagsRecords struct {
	Id        string    `json:"id"`
	TagId     string    `json:"tagId"`
	PostId    string    `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

func AddTagsRecords(t *TagsRecords) error {
	db := config.DB
	err := db.Create(&t).Error
	return err
}

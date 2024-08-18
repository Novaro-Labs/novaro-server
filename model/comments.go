package model

import (
	"fmt"
	"novaro-server/config"
	"time"
)

type Comments struct {
	Id        string     `json:"id"`
	UserId    string     `json:"userId"`
	PostId    string     `json:"postId"`
	ParentId  string     `json:"parentId"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	Children  []Comments `json:"children" gorm:"-"`
}

func AddComments(c *Comments) error {
	db := config.DB
	tx := db.Create(&c)
	return tx.Error
}

func GetCommentsCount(postId string) int64 {
	db := config.DB
	var count int64
	db.Table("comments").Where("post_id = ?", postId).Count(&count)
	return count
}

func GetCommentsListByPostId(postId string) (resp []Comments, err error) {
	db := config.DB
	err = db.Table("comments").Where("post_id = ?", postId).Find(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetCommentsListByParentId(parentId string) (resp []Comments, err error) {
	if parentId == "" {
		return nil, fmt.Errorf("parentId cannot be empty")
	}

	db := config.DB
	err = db.Table("comments").Where("parent_id = ?", parentId).Find(&resp).Error
	if err != nil {
		return resp, err
	}

	for i := range resp {
		children, err := GetCommentsListByParentId(fmt.Sprint(resp[i].Id))
		if err != nil {
			return nil, err
		}
		resp[i].Children = children
	}
	return resp, nil
}

func GetCommentsListByUserId(userId string) (resp []Comments, err error) {
	db := config.DB
	err = db.Table("comments").Where("user_id = ?", userId).Find(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, nil
}

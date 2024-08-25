package model

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"novaro-server/config"
	"strings"
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

func (Comments) TableName() string {
	return "comments"
}

func (u *Comments) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	u.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func AddComments(c *Comments) error {
	db := config.DB
	tx := db.Create(&c)

	// 用redis来记录评论数量
	rdb := config.RDB
	ctx := context.Background()
	rdb.ZIncrBy(ctx, "tweet:comments:count", 1, c.PostId)

	return tx.Error
}

func GetCommentsById(id string) (resp Comments, err error) {
	tx := config.DB.Where("id = ?", id).First(&resp)
	return resp, tx.Error
}

func GetCommentsCount(postId string) int64 {
	var count int64
	config.DB.Table("comments").Where("post_id = ?", postId).Count(&count)
	return count
}

func GetCommentsListByPostId(postId string) (resp []Comments, err error) {
	err = config.DB.Table("comments").Where("post_id = ?", postId).Find(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetCommentsListByParentId(parentId string) (resp []Comments, err error) {
	if parentId == "" {
		return nil, fmt.Errorf("parentId cannot be empty")
	}

	err = config.DB.Table("comments").Where("parent_id = ?", parentId).Find(&resp).Error
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
	err = config.DB.Table("comments").Where("user_id = ?", userId).Find(&resp).Error
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func DeleteById(id string) error {
	tx := config.DB.Table("comments").Where("id = ?", id).Delete(&Comments{})

	resp, _ := GetCommentsById(id)
	rdb := config.RDB
	ctx := context.Background()
	rdb.ZIncrBy(ctx, "tweet:comments:count", -1, resp.PostId)

	return tx.Error
}

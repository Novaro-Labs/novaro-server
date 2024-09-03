package model

import (
	"time"
)

// TODO 同样需要同步数据库

type RePosts struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	PostId    string    `json:"postId"`
	createdAt time.Time `json:"createdAt"`
}

func (RePosts) TableName() string {
	return "reposts"
}

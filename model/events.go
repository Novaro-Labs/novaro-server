package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type Events struct {
	Id        string `json:"id"`
	SourceId  string `json:"sourceId"`
	UserId    string `json:"userId"`
	Title     string `json:"title"`
	Website   string `json:"website"`
	ExpiredAt string `json:"expiredAt"`
	Status    int    `json:"status"`
	CreatedAt string `json:"createdAt"`
	Img       Imgs   `json:"img" gorm:"foreignKey:SourceId;references:SourceId;"`
}

func (Events) TableName() string {
	return "events"
}

func (e *Events) BeforeCreate(tx *gorm.DB) error {
	e.Id = strings.ReplaceAll(uuid.New().String(), "-", "")
	return nil
}

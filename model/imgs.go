package model

import (
	"time"
)

type Imgs struct {
	Id        string    `json:"id" `
	Path      string    `json:"path"`
	SourceId  string    `json:"sourceId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Imgs) TableName() string {
	return "imgs"
}

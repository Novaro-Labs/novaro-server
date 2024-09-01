package model

import (
	"github.com/google/uuid"
	"novaro-server/config"
	"strings"
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
func UploadFile(path string, sourceId string) error {
	tx := config.DB.Create(&Imgs{
		Path:     path,
		SourceId: sourceId,
		Id:       strings.ReplaceAll(uuid.New().String(), "-", ""),
	})
	return tx.Error
}
func GetImgsBySourceId(sourceId string) ([]Imgs, error) {
	var postsImgs []Imgs
	tx := config.DB.Where("source_id = ?", sourceId).Find(&postsImgs)
	return postsImgs, tx.Error
}

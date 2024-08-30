package model

import (
	"github.com/google/uuid"
	"novaro-server/config"
	"strings"
	"time"
)

type PostsImgs struct {
	Id        string    `json:"id" `
	Path      string    `json:"path"`
	SourceId  string    `json:"sourceId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (PostsImgs) TableName() string {
	return "posts_imgs"
}
func UploadFile(path string, sourceId string) error {
	tx := config.DB.Create(&PostsImgs{
		Path:     path,
		SourceId: sourceId,
		Id:       strings.ReplaceAll(uuid.New().String(), "-", ""),
	})
	return tx.Error
}
func GetPostImgsBySourceId(sourceId string) ([]PostsImgs, error) {
	var postsImgs []PostsImgs
	tx := config.DB.Where("source_id = ?", sourceId).Find(&postsImgs)
	return postsImgs, tx.Error
}

package dao

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"novaro-server/model"
	"strings"
)

type ImgsDao struct {
	db *gorm.DB
}

func NewImgsDao(db *gorm.DB) *ImgsDao {
	return &ImgsDao{
		db: db,
	}
}

func (d *ImgsDao) GetImgsBySourceId(sourceId string) ([]model.Imgs, error) {
	var imgs []model.Imgs
	tx := d.db.Where("source_id = ?", sourceId).Find(&imgs)
	return imgs, tx.Error
}

func (d *ImgsDao) UploadFile(path string, sourceId string) (*model.Imgs, error) {
	imgs := model.Imgs{
		Path:     path,
		SourceId: sourceId,
		Id:       strings.ReplaceAll(uuid.New().String(), "-", ""),
	}
	tx := d.db.Create(&imgs)
	return &imgs, tx.Error
}

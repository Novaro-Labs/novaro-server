package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type TagsRecordDao struct {
	db *gorm.DB
}

func NewTagsRecordDao(db *gorm.DB) *TagsRecordDao {

	return &TagsRecordDao{
		db: db,
	}
}

func (d *TagsRecordDao) AddTagsRecords(t *model.TagsRecords) error {
	err := d.db.Create(&t).Error
	return err
}

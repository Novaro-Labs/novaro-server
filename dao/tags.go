package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type TagsDao struct {
	db *gorm.DB
}

func NewTagsDao(db *gorm.DB) *TagsDao {
	return &TagsDao{
		db: db,
	}
}

func (d *TagsDao) GetTagListByPostId(postId string) (resp []model.Tags, err error) {
	err = d.db.Distinct("tags.*").Model(&model.Tags{}).
		Joins("JOIN tags_records ON tags.id = tags_records.tag_id").
		Where("tags_records.post_id = ?", postId).
		Find(&resp).Error
	return resp, err
}

func (d *TagsDao) GetTagsList() (resp []model.Tags, err error) {
	err = d.db.Model(model.Tags{}).Find(&resp).Error
	return resp, err
}

func (d *TagsDao) TagExists(id string) (bool, error) {
	var count int64
	tx := d.db.Model(&model.Tags{}).Where("id = ?", id).Count(&count)
	return count > 0, tx.Error
}

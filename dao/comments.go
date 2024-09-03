package dao

import (
	"fmt"
	"gorm.io/gorm"
	"novaro-server/model"
)

type CommentsDao struct {
	db *gorm.DB
}

func NewCommentsDao(db *gorm.DB) *CommentsDao {
	return &CommentsDao{
		db: db,
	}
}

func (d *CommentsDao) Create(c *model.Comments) error {
	tx := d.db.Create(&c)
	return tx.Error
}

func (d *CommentsDao) GetById(id string) (resp model.Comments, err error) {
	tx := d.db.Where("id = ?", id).First(&resp)
	return resp, tx.Error
}

func (d *CommentsDao) GetCount(postId string) int64 {
	var count int64
	d.db.Table("comments").Where("post_id = ?", postId).Count(&count)
	return count
}

func (d *CommentsDao) GetListByPostId(postId string) (resp []model.Comments, err error) {
	err = d.db.Table("comments").Where("post_id = ?", postId).Find(&resp).Error
	return resp, nil
}

func (d *CommentsDao) GetListByParentId(parentId string) (resp []model.Comments, err error) {

	err = d.db.Table("comments").Where("parent_id = ?", parentId).Find(&resp).Error
	if err != nil {
		return resp, err
	}

	for i := range resp {
		children, err := d.GetListByParentId(fmt.Sprint(resp[i].Id))
		if err != nil {
			return nil, err
		}
		resp[i].Children = children
	}
	return resp, nil
}

func (d *CommentsDao) GetListByUserId(userId string) (resp []model.Comments, err error) {
	err = d.db.Table("comments").Where("user_id = ?", userId).Find(&resp).Error
	return resp, nil
}

func (d *CommentsDao) DeleteById(id string) error {
	tx := d.db.Table("comments").Where("id = ?", id).Delete(&model.Comments{})
	return tx.Error
}

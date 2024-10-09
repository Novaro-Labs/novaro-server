package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type LikeDao struct {
	dao *gorm.DB
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{
		dao: db,
	}
}

func (d *LikeDao) Add(like *model.Likes) (*model.Likes, error) {
	err := d.dao.Create(&like).Error
	return like, err
}

func (d *LikeDao) BatchAdd(postId string, userIds []string) error {
	for _, userId := range userIds {
		d.dao.Create(&model.Likes{
			UserId: userId,
			PostId: postId,
		})
	}
	return nil
}

func (d *LikeDao) Delete(id string) bool {
	err := d.dao.Where("id=?", id).Delete(&model.Likes{}).Error
	return err == nil
}

func (d *LikeDao) IsLikedByUser(postId string, userId string) string {
	var like model.Likes
	d.dao.Where("post_id = ? AND user_id = ?", postId, userId).First(&like)
	return like.Id
}

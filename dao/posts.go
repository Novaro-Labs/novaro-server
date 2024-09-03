package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type PostsDao struct {
	db *gorm.DB
}

func NewPostsDao(db *gorm.DB) *PostsDao {
	return &PostsDao{
		db: db,
	}
}

func (d *PostsDao) GetPostsList(p *model.PostsQuery) (resp []model.Posts, err error) {
	query := d.db.Table("posts")
	if p.Id != "" {
		query = query.Where("id = ?", p.Id)
	}

	err = query.Find(&resp).Error
	return resp, err
}

func (d *PostsDao) GetPostsById(id string) (resp model.Posts, err error) {
	err = d.db.Table("posts").Where("id = ?", id).Find(&resp).Error
	return resp, err
}

func (d *PostsDao) PostExists(id string) (bool, error) {
	var count int64
	err := d.db.Model(&model.Posts{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (d *PostsDao) GetPostsByUserId(userId string) (resp []model.Posts, err error) {
	err = d.db.Table("posts").Where("user_id = ?", userId).Find(&resp).Error
	return resp, nil
}

func (d *PostsDao) Save(posts *model.Posts) error {
	var data = model.Posts{
		Id:         posts.Id,
		UserId:     posts.UserId,
		Content:    posts.Content,
		OriginalId: posts.OriginalId,
		SourceId:   posts.SourceId,
	}

	tx := d.db.Create(&data)
	return tx.Error
}

func (d *PostsDao) Update(posts *model.Posts) error {
	tx := d.db.Updates(&posts)
	return tx.Error
}

func (d *PostsDao) UpdateBatch(posts []model.Posts) error {
	// 开始事务
	err := d.db.Transaction(func(tx *gorm.DB) error {
		for _, post := range posts {
			// 更新每个 post
			if err := tx.Model(&post).Updates(model.Posts{
				Content:           post.Content,
				CommentsAmount:    post.CommentsAmount,
				CollectionsAmount: post.CollectionsAmount,
				RepostsAmount:     post.RepostsAmount,
				Tags:              post.Tags,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (d *PostsDao) Delete(id string) error {
	tx := d.db.Where("id = ?", id).Delete(&model.Posts{})
	return tx.Error
}

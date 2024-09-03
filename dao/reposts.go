package dao

import "gorm.io/gorm"

type RePostsDao struct {
	db *gorm.DB
}

func NewRePostsDao(db *gorm.DB) *RePostsDao {
	return &RePostsDao{
		db: db,
	}
}

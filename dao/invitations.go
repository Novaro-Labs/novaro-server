package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
)

type InvitationsDao struct {
	db *gorm.DB
}

func NewInvitationsDao(db *gorm.DB) *InvitationsDao {
	return &InvitationsDao{db: db}
}

func (d *InvitationsDao) Save(invitations *model.Invitations) error {
	tx := d.db.Create(&invitations)
	return tx.Error
}

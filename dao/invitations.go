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

func (d *InvitationsDao) CheckInvitationRewards(inviteeId, invitationCode string) (bool, error) {
	var count int64
	tx := d.db.Model(&model.Invitations{}).Where("invitee_id = ? AND invitation_code = ?", inviteeId, invitationCode).Count(&count)
	return count > 0, tx.Error

}

func (d *InvitationsDao) CheckInvitee(inviteeId string) (bool, error) {
	var count int64
	tx := d.db.Model(&model.Invitations{}).Where("invitee_id = ?", inviteeId).Count(&count)
	return count > 0, tx.Error

}

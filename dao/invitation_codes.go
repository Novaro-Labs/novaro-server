package dao

import (
	"gorm.io/gorm"
	"novaro-server/model"
	"time"
)

type InvitationCodesDao struct {
	db *gorm.DB
}

func NewInvitationCodesDao(db *gorm.DB) *InvitationCodesDao {
	return &InvitationCodesDao{
		db: db,
	}
}

func (d *InvitationCodesDao) MakeInvitationCodes(codes *model.InvitationCodes) error {
	tx := d.db.Create(&codes)
	return tx.Error
}

func (d *InvitationCodesDao) CheckInvitationCodes(code string) (bool, error) {
	var invitationCodes model.InvitationCodes
	tx := d.db.Table("invitation_codes").Where("code = ?", code).Find(&invitationCodes)
	if tx.Error != nil {
		return false, tx.Error
	}
	if tx.RowsAffected == 0 {
		return false, nil
	}
	if invitationCodes.ExpiresAt.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

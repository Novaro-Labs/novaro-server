package model

import (
	"crypto/rand"
	"encoding/hex"
	"novaro-server/config"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvitationCodes struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	Status    int       `json:"status"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (InvitationCodes) TableName() string {
	return "invitation_codes"
}

func (i *InvitationCodes) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	i.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

func MakeInvitationCodes(userId string) (*string, *time.Time, error) {
	db := config.DB

	var code string
	var err error
	for {
		code, err = MakeInvitationCode(config.InvitatioCodeLength)
		if err != nil {
			return nil, nil, err
		}

		exist, err := CheckInvitationCodes(code)
		if err != nil {
			return nil, nil, err
		}

		if !exist {
			break
		}
	}

	var data = InvitationCodes{
		UserId:    userId,
		Code:      code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(config.InvitatioCodeExpiration),
	}

	tx := db.Create(&data)
	return &code, &data.ExpiresAt, tx.Error
}

func CheckInvitationCodes(code string) (bool, error) {
	db := config.DB
	var invitationCodes InvitationCodes
	tx := db.Table("invitation_codes").Where("code = ?", code).Find(&invitationCodes)
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

func MakeInvitationCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes)[:length], nil
}

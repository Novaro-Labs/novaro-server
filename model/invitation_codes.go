package model

import (
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

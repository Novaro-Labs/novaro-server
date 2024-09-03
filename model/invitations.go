package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invitations struct {
	Id             string    `json:"id"`
	InviterId      string    `json:"inviterId"`
	InviteeId      string    `json:"inviteeId"`
	InvitationCode string    `json:"invitationCode"`
	InvitedAt      time.Time `json:"invitedAt"`
}

func (Invitations) TableName() string {
	return "invitations"
}

func (i *Invitations) BeforeCreate(tx *gorm.DB) error {
	u2 := uuid.New()
	i.Id = strings.ReplaceAll(u2.String(), "-", "")
	return nil
}

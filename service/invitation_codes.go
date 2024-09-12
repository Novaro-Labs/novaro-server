package service

import (
	"crypto/rand"
	"encoding/hex"
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
	"time"
)

type InvitationCodesService struct {
	dao *dao.InvitationCodesDao
}

func NewInvitationCodesService() *InvitationCodesService {
	return &InvitationCodesService{
		dao: dao.NewInvitationCodesDao(model.GetDB()),
	}
}

func (s *InvitationCodesService) MakeInvitationCodes(userId string) (*string, *time.Time, error) {

	client := config.Get().Client

	var code string
	var err error
	for {
		code, err = s.MakeInvitationCode(client.InvitationCodeLength)
		if err != nil {
			return nil, nil, err
		}

		exist, err := s.dao.CheckInvitationCodes(code)
		if err != nil {
			return nil, nil, err
		}

		if !exist {
			break
		}
	}

	var data = &model.InvitationCodes{
		UserId:    userId,
		Code:      code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(client.InvitationCodeExpireDay),
	}
	err = s.dao.MakeInvitationCodes(data)

	return &code, &data.ExpiresAt, err
}

func (s *InvitationCodesService) CheckInvitationCodes(code string) (bool, error) {
	return s.dao.CheckInvitationCodes(code)
}

func (s *InvitationCodesService) MakeInvitationCode(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes)[:length], nil
}

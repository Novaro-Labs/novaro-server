package service

import (
	"novaro-server/config"
	"novaro-server/dao"
	"novaro-server/model"
)

type InvitationsService struct {
	dao *dao.InvitationsDao
}

func NewInvitationsService() *InvitationsService {
	return &InvitationsService{
		dao: dao.NewInvitationsDao(config.DB),
	}
}

func (s *InvitationsService) Save(i *model.Invitations) error {
	var data = model.Invitations{
		InviterId:      i.InviterId,
		InviteeId:      i.InviteeId,
		InvitationCode: i.InvitationCode,
		InvitedAt:      i.InvitedAt,
	}

	return s.dao.Save(&data)
}

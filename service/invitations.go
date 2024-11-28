package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type InvitationsService struct {
	dao *dao.InvitationsDao
}

func NewInvitationsService() *InvitationsService {
	return &InvitationsService{
		dao: dao.NewInvitationsDao(model.GetDB()),
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

func (s *InvitationsService) CheckInvitationRewards(inviteeId, code string) (bool, error) {
	return s.dao.CheckInvitationRewards(inviteeId, code)
}

func (s *InvitationsService) CheckInvitee(inviteeId string) (bool, error) {
	return s.dao.CheckInvitee(inviteeId)
}

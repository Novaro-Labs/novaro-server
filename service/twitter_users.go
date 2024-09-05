package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type TwitterUserService struct {
	dao *dao.TwitterUserDao
}

func NewTwiitterUserService() *TwitterUserService {
	return &TwitterUserService{
		dao: dao.NewTwitterUserDao(model.GetDB()),
	}
}

func (s *TwitterUserService) SaveTwitterUsers(users *model.TwitterUsers) error {
	var data = model.TwitterUsers{
		TwitterId:        users.TwitterId,
		TwitterUserName:  users.TwitterUserName,
		TwitterAvatar:    users.TwitterAvatar,
		TwitterFollowers: users.TwitterFollowers,
		TwitterCreatedAt: users.TwitterCreatedAt,
	}

	return s.dao.SaveTwitterUsers(&data)
}

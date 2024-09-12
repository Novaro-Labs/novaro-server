package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type UserService struct {
	dao *dao.UsersDao
}

func NewUserService() *UserService {
	return &UserService{
		dao: dao.NewUsersDao(model.GetDB()),
	}
}

func (s *UserService) SaveUsers(users *model.Users) (string, error) {

	return s.dao.SaveUsers(users)
}

func (s *UserService) UserExists(userId string) (bool, error) {
	return s.dao.UserExists(userId)
}

func (s *UserService) GetById(id string) (*model.Users, error) {
	return s.dao.GetById(id)
}

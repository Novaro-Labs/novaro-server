package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type EventsService struct {
	dao *dao.EventsDao
}

func NewEventsService() *EventsService {
	return &EventsService{
		dao: dao.NewEventsDao(model.GetDB()),
	}
}

func (s *EventsService) Create(r *model.Events) error {

	return s.dao.Create(r)
}

func (s *EventsService) Delete(id string) error {
	return s.dao.Delete(id)
}

func (s *EventsService) Update(r *model.Events) error {
	if r.Id == "" {
		err := s.dao.Create(r)
		return err
	}

	tx := s.dao.Updates(r)
	return tx
}
func (s *EventsService) GetById(id string) (model.Events, error) {
	return s.dao.Get(id)
}

func (s *EventsService) GetList(r *model.Events) ([]model.Events, error) {
	return s.dao.GetList(r)
}

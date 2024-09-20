package service

import (
	"novaro-server/dao"
	"novaro-server/model"
)

type ImgsService struct {
	dao *dao.ImgsDao
}

func NewImgsService() *ImgsService {
	return &ImgsService{
		dao: dao.NewImgsDao(model.GetDB()),
	}
}

func (s *ImgsService) GetBySourceId(sourceId string) ([]model.Imgs, error) {
	return s.dao.GetImgsBySourceId(sourceId)
}

func (s *ImgsService) UploadFile(path string, sourceId string) (*model.Imgs, error) {
	return s.dao.UploadFile(path, sourceId)
}

package api

import "novaro-server/service"

type PostPointsApi struct {
	service *service.PostPointsService
}

func NewPostPointsApi() *PostPointsApi {
	return &PostPointsApi{
		service: service.NewPostPointsService(),
	}
}

func (api *PostPointsApi) SyncData() {
	api.service.SyncData()
}

package api

import (
	"novaro-server/service"
)

type RePostsApi struct {
	service *service.RePostsService
}

func NewRePostApi() *RePostsApi {
	return &RePostsApi{
		service: service.NewRePostsService(),
	}
}

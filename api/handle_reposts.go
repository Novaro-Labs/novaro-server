package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
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

func (api *RePostsApi) AddRePosts(c *gin.Context) {
	var rePosts model.RePosts
	if err := c.ShouldBindJSON(&rePosts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}
	if err := api.service.AddRePosts(&rePosts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added reposts"})
}

package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type CollectionsApi struct {
	UserId  string `json:"userId" binding:"required"`
	PostId  string `json:"postId" binding:"required"`
	service *service.CollectionsService
}

func NewCollectionsApi() *CollectionsApi {
	return &CollectionsApi{
		service: service.NewCollectionsService(),
	}
}

func (api *CollectionsApi) Create(c *gin.Context) {
	var collections model.Collections

	if err := c.ShouldBindJSON(&collections); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if collections.UserId == "" || collections.PostId == "" {
		c.JSON(400, gin.H{"error": "userId and postId is required"})
		return
	}

	// 收藏
	if err := api.service.AddOrRemove(&collections); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "success"})
}

func (api *CollectionsApi) Sync() {
	api.service.SyncToDatabase()
}

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
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	if collections.UserId == "" || collections.PostId == "" {
		c.JSON(400, gin.H{"msg": "userId and postId is required"})
		return
	}

	// 收藏
	if err := api.service.AddOrRemove(&collections); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func (api *CollectionsApi) Sync() {
	api.service.SyncToDatabase()
}

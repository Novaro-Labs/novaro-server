package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type LikeApi struct {
	service *service.LikeService
}

func NewLikeApi() *LikeApi {
	return &LikeApi{
		service: service.NewLikeService(),
	}
}

func (api *LikeApi) Like(c *gin.Context) {
	var like model.LikesReq
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	if err := api.service.Like(&like); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"data": false,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": true,
		"msg":  "success",
	})
}

func (api *LikeApi) FlushToDatabase() {
	api.service.FlushToDatabase()
}

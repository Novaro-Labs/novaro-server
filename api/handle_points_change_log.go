package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type PointsChangeLogApi struct {
	service *service.PointsChangeLogService
}

func NewPointsChangeLogApi() *PointsChangeLogApi {
	return &PointsChangeLogApi{service: service.NewPointsChangeLogService()}
}

func (api *PointsChangeLogApi) GetList(c *gin.Context) {
	var request model.PointsChangeLogRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	list, err := api.service.GetList(&request)

	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": list,
		"msg":  "success",
	})
}

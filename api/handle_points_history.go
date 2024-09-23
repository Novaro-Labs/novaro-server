package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type PointsHistoryApi struct {
	service *service.PointsHistoryService
}

func NewPointsHistoryApi() *PointsHistoryApi {
	return &PointsHistoryApi{
		service: service.NewPointsHistoryService(),
	}
}

func (api *PointsHistoryApi) GetList(c *gin.Context) {
	var p model.PointsHistoryQuery
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	history, err := api.service.GetList(&p)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": history,
		"msg":  "success",
	})
}

func (api *PointsHistoryApi) Statistics(c *gin.Context) {
	value := c.Query("wallet")
	if value == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "wallet is required",
		})
		return
	}

	statistics, err := api.service.Statistics(value, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": statistics,
		"msg":  "success",
	})
}

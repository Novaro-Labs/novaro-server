package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/service"
)

type NftInfoApi struct {
	service *service.NftInfoService
}

func NewNftInfoApi() *NftInfoApi {
	return &NftInfoApi{
		service: service.NewNftInfoService(),
	}
}

func (api *NftInfoApi) GetNftInfo(c *gin.Context) {
	value := c.Query("wattle")

	if value == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "wattle is required",
		})
		return
	}

	wallet, err := api.service.GetByWallet(value)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": wallet,
		"msg":  "success",
	})
}

package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type NftTokensApi struct {
	service *service.NftTokensService
}

func NewNftTokensApi() *NftTokensApi {
	return &NftTokensApi{
		service: service.NewNftTokensService(),
	}
}

func (api *NftTokensApi) GetTokensByWallet(c *gin.Context) {
	value := c.Query("wallet")

	if value != "" {
		tokens, err := api.service.GetTokensByWallet(value)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 200,
			"data": tokens,
			"msg":  "success",
		})
	} else {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "wallet is empty",
		})
	}

}

func (api *NftTokensApi) SaveNftToken(c *gin.Context) {
	var tokens model.NftTokens
	if err := c.ShouldBindJSON(&tokens); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	if err := api.service.SaveNftToken(&tokens); err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

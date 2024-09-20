package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/utils"
)

type UploadApi struct {
}

func NewUploadApi() *UploadApi {
	return &UploadApi{}
}

func (api *UploadApi) UploadFile(c *gin.Context) {
	files, err := utils.UploadFiles(c.Writer, c.Request)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": files,
		"msg":  "success",
	})
}

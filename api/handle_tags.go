package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type TagsApi struct {
}

func (TagsApi) GetTagsList(c *gin.Context) {
	tags, err := model.GetTagsList()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tags)
}

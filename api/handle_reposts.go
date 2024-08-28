package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type RePostsApi struct {
}

func (RePostsApi) AddRePosts(c *gin.Context) {
	var rePosts model.RePosts
	if err := c.ShouldBindJSON(&rePosts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return

	}
	if err := model.AddRePosts(&rePosts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added reposts"})
}

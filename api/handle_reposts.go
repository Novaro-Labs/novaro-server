package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type RePostsApi struct {
}

// AddRePosts godoc
// @Summary Add a new repost
// @Description Add a new repost to the system
// @Tags reposts
// @Accept json
// @Produce json
// @Param repost body model.RePosts true "Repost object"
// @Success 200 " Successfully added reposts"
// @Failure 400
// @Router /v1/api/reposts/add [post]
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

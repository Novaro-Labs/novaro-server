package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type CollectionsApi struct {
	UserId string `json:"userId" binding:"required"`
	PostId string `json:"postId" binding:"required"`
}

// CollectionsTweet godoc
// @Summary Collect a tweet
// @Description Add a tweet to user's collection
// @Tags collections
// @Accept json
// @Produce json
// @Param  body body CollectionsApi true "Collection information"
// @Success 200
// @Failure 400
// @Router /vi/api/collections/add [post]
func (CollectionsApi) CollectionsTweet(c *gin.Context) {
	var collections model.Collections

	if err := c.ShouldBindJSON(&collections); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if collections.UserId == "" || collections.PostId == "" {
		c.JSON(400, gin.H{"error": "userId and postId is required"})
		return
	}

	exist, err := model.UserExists(collections.UserId)
	if err != nil || exist == false {
		c.JSON(400, gin.H{"error": "userId is not exist"})
		return
	}

	postExist, err := model.PostExists(collections.PostId)
	if err != nil || postExist == false {
		c.JSON(400, gin.H{"error": "postId is not exist"})
		return
	}

	if err := model.CollectionsTweet(&collections); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "success"})
}

// UnCollectionsTweet godoc
// @Summary Cancel collect a tweet
// @Description Cancel add a tweet to user's collection
// @Tags collections
// @Accept json
// @Produce json
// @Param  body body CollectionsApi true "Collection information"
// @Success 200
// @Failure 400
// @Router /vi/api/collections/remove [post]
func (CollectionsApi) UnCollectionsTweet(c *gin.Context) {
	var coll model.Collections

	if err := c.ShouldBindJSON(&coll); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if coll.UserId == "" || coll.PostId == "" {
		c.JSON(400, gin.H{"error": "userId and postId is required"})
		return
	}

	exist, err := model.UserExists(coll.UserId)
	if err != nil || exist == false {
		c.JSON(400, gin.H{"error": "userId is not exist"})
		return
	}

	postExist, err := model.PostExists(coll.PostId)
	if err != nil || postExist == false {
		c.JSON(400, gin.H{"error": "postId is not exist"})
		return
	}

	if err := model.UnCollectionsTweet(&coll); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "success"})
}

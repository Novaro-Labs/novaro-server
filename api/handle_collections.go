package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type CollectionsApi struct {
}

// CollectionsTweet godoc
// @Summary Collect a tweet
// @Description Add a tweet to user's collections
// @Tags collections
// @Accept json
// @Produce json
// @Param collection body model.Collections true "Collection information"
// @Success 200 {object} model.Collections
// @Failure 400
// @Router /v1/api/collections/add [post]
func (CollectionsApi) CollectionsTweet(c *gin.Context) {
	var collections model.Collections

	if err := c.ShouldBindJSON(&collections); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := model.CollectionsTweet(&collections); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, collections)
}

// UnCollectionsTweet godoc
// @Summary Remove a tweet from collections
// @Description Remove a tweet from user's collections
// @Tags collections
// @Accept json
// @Produce json
// @Param collection body model.Collections true "Collection information"
// @Success 200 {object} model.Collections
// @Failure 400
// @Router /v1/api/collections/remove [delete]
func (CollectionsApi) UnCollectionsTweet(c *gin.Context) {
	var coll model.Collections

	if err := c.ShouldBindJSON(&coll); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := model.UnCollectionsTweet(&coll); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, coll)
}

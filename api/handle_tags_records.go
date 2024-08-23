package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type TagsRecordsApi struct {
}

func (TagsRecordsApi) AddTagsRecords(c *gin.Context) {
	var tagsRecords model.TagsRecords
	if err := c.ShouldBindJSON(&tagsRecords); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := model.AddTagsRecords(&tagsRecords); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added tags records"})

}

package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type TagsRecordsApi struct {
	TagId  string `json:"tagId"`
	PostId string `json:"postId"`
}

// AddTagsRecords godoc
// @Summary Add new tags records
// @Description Add new tags records to the database
// @Tags tags-records
// @Accept json
// @Produce json
// @Param tagsRecords body TagsRecordsApi true "Tags records to add"
// @Success 200 "Successfully added tags records"
// @Failure 400
// @Router /v1/api/tags/records/add [post]
func (TagsRecordsApi) AddTagsRecords(c *gin.Context) {
	var tagsRecords model.TagsRecords
	if err := c.ShouldBindJSON(&tagsRecords); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	exists, err := model.PostExists(tagsRecords.PostId)
	if err != nil || !exists {
		c.JSON(400, gin.H{"error": "Post does not exist"})
		return
	}

	tagExists, err := model.TagExists(tagsRecords.TagId)
	if err != nil || !tagExists {
		c.JSON(400, gin.H{"error": "Tag does not exist"})
		return
	}

	if err := model.AddTagsRecords(&tagsRecords); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added tags records"})
}

package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type TagsRecordsApi struct {
}

// AddTagsRecords godoc
// @Summary Add new tags records
// @Description Add new tags records to the system
// @Tags tags-records
// @Accept json
// @Produce json
// @Param tagsRecords body model.TagsRecords true "Tags Records object"
// @Success 200 " Successfully added tags records"
// @Failure 400 " Error adding tags records"
// @Router /v1/api/tags/records/add [post]
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

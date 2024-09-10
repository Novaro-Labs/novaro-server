package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type TagsRecordsApi struct {
	TagId   string                     `json:"tagId"`
	PostId  string                     `json:"postId"`
	UserId  string                     `json:"userId"`
	service *service.TagsRecordService `json:"-"`
}

func NewTagsRecordApi() *TagsRecordsApi {
	return &TagsRecordsApi{
		service: service.NewTagsRecordService(),
	}
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
func (api *TagsRecordsApi) AddTagsRecords(c *gin.Context) {
	var tagsRecords model.TagsRecords
	if err := c.ShouldBindJSON(&tagsRecords); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	if err := api.service.Create(&tagsRecords); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "Successfully added tags records"})
}

func (api *TagsRecordsApi) SyncData() {
	api.service.SyncData()
	api.service.CleanExpiredTags()
}

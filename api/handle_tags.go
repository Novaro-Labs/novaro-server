package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/service"
)

type TagsApi struct {
	service *service.TagsService
}

func NewTagsApi() *TagsApi {
	return &TagsApi{
		service: service.NewTagsService(),
	}
}

// GetTagsList godoc
// @Summary Get list of tags
// @Description Retrieve a list of all tags
// @Tags tags
// @Produce json
// @Success 200 {array} model.Tags "Successful operation"
// @Failure 400
// @Router /v1/api/tags/list [get]
func (api *TagsApi) GetTagsList(c *gin.Context) {
	tags, err := api.service.GetTagsList()
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": tags,
	})
}

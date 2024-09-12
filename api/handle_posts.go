package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type PostsApi struct {
	UserId  string               `json:"userId"`
	Content string               `json:"content"`
	service *service.PostService `json:"-"`
}

type RePosts struct {
	PostsApi   PostsApi
	OriginalId string `json:"originalId"`
}

func NewPostsApi() *PostsApi {
	return &PostsApi{
		service: service.NewPostService(),
	}
}

func (api *PostsApi) GetPostsList(c *gin.Context) {
	var postsQuery model.PostsQuery
	if err := c.ShouldBind(&postsQuery); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	resp, err := api.service.GetList(&postsQuery)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}

func (api *PostsApi) SavePosts(c *gin.Context) {
	var posts model.Posts

	if err := c.ShouldBindJSON(&posts); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	if err := api.service.Save(&posts); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func (api *PostsApi) SaveReposts(c *gin.Context) {
	var posts model.Posts

	if err := c.ShouldBindJSON(&posts); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	if posts.OriginalId == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "OriginalId is required",
		})
		return
	}

	if err := api.service.Save(&posts); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

func (api *PostsApi) DelPostsById(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	if err := api.service.Delete(id); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

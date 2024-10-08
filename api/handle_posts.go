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

func (api *PostsApi) GetLikeByUser(c *gin.Context) {
	var postsQuery model.PostsQuery
	if err := c.ShouldBind(&postsQuery); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	if postsQuery.UserId == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "userId is required",
		})
		return
	}

	user, err := api.service.GetList(&postsQuery)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": user,
		"msg":  "sussess",
	})
}

func (api *PostsApi) GetPostById(c *gin.Context) {
	value := c.Query("postId")
	if value == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "postId is required",
		})
		return
	}

	res, err := api.service.GetById(value)

	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": res,
		"msg":  "sussess",
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

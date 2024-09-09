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

// GetPostsList godoc
// @Summary Get a list of posts
// @Description Retrieves a list of posts based on the provided query parameters
// @Tags posts
// @Accept json
// @Produce json
// @Param query body model.PostsQuery false "Query parameters"
// @Success 200 {array} model.Posts
// @Failure 400
// @Router /v1/api/posts/getPostsList [post]
func (api *PostsApi) GetPostsList(c *gin.Context) {
	var postsQuery model.PostsQuery
	if err := c.ShouldBind(&postsQuery); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := api.service.GetList(&postsQuery)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": resp,
	})
}

// SavePosts godoc
// @Summary Save a new post
// @Description Creates a new post in the database
// @Tags posts
// @Accept json
// @Produce json
// @Param post body PostsApi true "Post object"
// @Success 200 {object} model.Posts
// @Failure 400
// @Router /v1/api/posts/savePosts [post]
func (api *PostsApi) SavePosts(c *gin.Context) {
	var posts model.Posts

	if err := c.ShouldBindJSON(&posts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := api.service.Save(&posts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// SavePosts godoc
// @Summary Save a new repost
// @Description Creates a new repost in the database
// @Tags posts
// @Accept json
// @Produce json
// @Param post body RePosts true "Post object"
// @Success 200 {object} model.Posts
// @Failure 400
// @Router /v1/api/posts/saveRePosts [post]
func (api *PostsApi) SaveReposts(c *gin.Context) {
	var posts model.Posts

	if err := c.ShouldBindJSON(&posts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := api.service.Save(&posts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// DelPostsById godoc
// @Summary Delete a post by ID
// @Description Deletes a post from the database based on the provided ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id query string true "Post ID"
// @Success 200
// @Failure 400
// @Router /v1/api/posts/delPostsById [delete]
func (api *PostsApi) DelPostsById(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	if err := api.service.Delete(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

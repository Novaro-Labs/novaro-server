package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type PostsApi struct {
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

// GetPostsById godoc
// @Summary Get a post by ID
// @Description Retrieves a post from the database based on the provided ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id query string true "PostID"
// @Success 200 {object} model.Posts
// @Failure 400
// @Router /v1/api/posts/getPostsById [get]
func (PostsApi) GetPostsById(c *gin.Context) {
	value := c.Query("id")

	resp, err := model.GetPostsById(value)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

// GetPostsByUserId godoc
// @Summary Get posts by user ID
// @Description Retrieves all posts from the database for a specific user
// @Tags posts
// @Accept json
// @Produce json
// @Param userId query string true "UserID"
// @Success 200 {array} model.Posts
// @Failure 400
// @Router /v1/api/posts/getPostsByUserId [get]
func (PostsApi) GetPostsByUserId(c *gin.Context) {
	id := c.Query("userId")
	resp, err := model.GetPostsByUserId(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
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
func (PostsApi) GetPostsList(c *gin.Context) {
	var postsQuery model.PostsQuery
	if err := c.ShouldBind(&postsQuery); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := model.GetPostsList(&postsQuery)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
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
func (PostsApi) SavePosts(c *gin.Context) {
	var posts model.Posts

	if err := c.ShouldBindJSON(&posts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := model.SavePosts(&posts); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, posts)
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
func (PostsApi) DelPostsById(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	if err := model.DelPostsById(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "success"})
}

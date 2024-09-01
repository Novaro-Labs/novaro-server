package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
)

type CommentsApi struct {
	UserId   string `json:"userId"`
	PostId   string `json:"postId"`
	ParentId string `json:"parentId"`
	Content  string `json:"content"`
}

// AddComments godoc
// @Summary Add a new comment
// @Description Add a new comment to the system
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body CommentsApi true "Comment object"
// @Success 200
// @Failure 400
// @Router /v1/api/comments/add [post]
func (CommentsApi) AddComments(c *gin.Context) {
	var comments model.Comments

	if err := c.ShouldBindJSON(&comments); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := model.AddComments(&comments); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "success"})
}

// GetCommentsListByPostId godoc
// @Summary Get comments by post ID
// @Description Get a list of comments for a specific post
// @Tags comments
// @Accept json
// @Produce json
// @Param postId query string true "Post ID"
// @Success 200 {array} model.Comments
// @Failure 400
// @Router /v1/api/comments/getCommentsListByPostId [get]
func (CommentsApi) GetCommentsListByPostId(c *gin.Context) {
	postId := c.Query("postId")
	comments, err := model.GetCommentsListByPostId(postId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, comments)
}

// GetCommentsListByParentId godoc
// @Summary Get comments by parent ID
// @Description Get a list of child comments for a specific parent comment
// @Tags comments
// @Accept json
// @Produce json
// @Param parentId query string true "Parent Comment ID"
// @Success 200 {array} model.Comments
// @Failure 400
// @Router /v1/api/comments/getCommentsListByParentId [get]
func (CommentsApi) GetCommentsListByParentId(c *gin.Context) {
	parentId := c.Query("parentId")
	comments, err := model.GetCommentsListByParentId(parentId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, comments)
}

// GetCommentsListByUserId godoc
// @Summary Get comments by user ID
// @Description Get a list of comments made by a specific user
// @Tags comments
// @Accept json
// @Produce json
// @Param userId query string true "UserID"
// @Success 200 {array} model.Comments
// @Failure 400
// @Router /v1/api/comments/getCommentsListByUserId [get]
func (CommentsApi) GetCommentsListByUserId(c *gin.Context) {
	userId := c.Query("userId")
	comments, err := model.GetCommentsListByUserId(userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, comments)
}

// DeleteById godoc
// @Summary Delete a comment by ID
// @Description Deletes a comment from the database based on the provided ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id query string true "Comment ID"
// @Success 200
// @Failure 400
// @Router /v1/api/comments/delete [delete]
func (CommentsApi) DeleteById(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}

	err := model.DeleteById(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "success"})
}

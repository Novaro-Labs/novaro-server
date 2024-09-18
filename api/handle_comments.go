package api

import (
	"github.com/gin-gonic/gin"
	"novaro-server/model"
	"novaro-server/service"
)

type CommentsApi struct {
	UserId   string                   `json:"userId"`
	PostId   string                   `json:"postId"`
	ParentId string                   `json:"parentId"`
	Content  string                   `json:"content"`
	service  *service.CommentsService `json:"-"`
}

func NewCommentApi() *CommentsApi {
	return &CommentsApi{
		service: service.NewCommentService(),
	}
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
func (api *CommentsApi) GetCommentsListByPostId(c *gin.Context) {
	postId := c.Query("postId")
	comments, err := api.service.GetListByPostId(postId)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
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
func (api *CommentsApi) GetCommentsListByParentId(c *gin.Context) {
	parentId := c.Query("parentId")
	comments, err := api.service.GetListByParentId(parentId)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
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
func (api *CommentsApi) GetCommentsListByUserId(c *gin.Context) {
	userId := c.Query("userId")
	comments, err := api.service.GetListByUserId(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": comments,
	})
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
func (api *CommentsApi) AddComments(c *gin.Context) {
	var comments model.Comments

	if err := c.ShouldBindJSON(&comments); err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	comment, count, err := api.service.Create(&comments)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":  200,
		"data":  comment,
		"total": count,
		"msg":   "success",
	})
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
func (api *CommentsApi) DeleteById(c *gin.Context) {
	id := c.Query("id")
	postId := c.Query("postId")
	if id == "" || postId == "" {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "id or postId is empty",
		})
		return
	}

	byId, err2 := api.service.GetById(id)
	if err2 != nil || byId == nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err2.Error(),
		})
		return
	}

	count, err := api.service.Delete(id, postId)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": count, "msg": "success"})
}

func (api *CommentsApi) SyncCommentsToDB() {
	api.service.SyncCommentsToDB()
}

package routes

import (
	"novaro-server/api"

	"github.com/gin-gonic/gin"
)

func AddOtherRoutes(r *gin.RouterGroup) {
	collections := r.Group("/api/collections")
	{
		collections.POST("/add", api.CollectionsApi{}.CollectionsTweet)
		collections.POST("/remove", api.CollectionsApi{}.UnCollectionsTweet)
	}

	comments := r.Group("/api/comments")
	{
		comments.GET("/getCommentsListByPostId", api.CommentsApi{}.GetCommentsListByPostId)
		comments.GET("/getCommentsListByParentId", api.CommentsApi{}.GetCommentsListByParentId)
		comments.GET("/getCommentsListByUserId", api.CommentsApi{}.GetCommentsListByUserId)
		comments.POST("/add", api.CommentsApi{}.AddComments)
		comments.DELETE("/delete", api.CommentsApi{}.DeleteById)
	}

	posts := r.Group("/api/posts")
	{
		posts.GET("/getPostsById", api.PostsApi{}.GetPostsById)
		posts.GET("/getPostsByUserId", api.PostsApi{}.GetPostsByUserId)
		posts.POST("/getPostsList", api.PostsApi{}.GetPostsList)
		posts.POST("/savePosts", api.PostsApi{}.SavePosts)
		posts.POST("/saveRePosts", api.PostsApi{}.SavePosts)
		posts.DELETE("/delPostsById", api.PostsApi{}.DelPostsById)
	}

	reposts := r.Group("/api/reposts")
	{
		reposts.POST("/add", api.RePostsApi{}.AddRePosts)
	}

	tags := r.Group("/api/tags")
	{
		tags.GET("/list", api.TagsApi{}.GetTagsList)
	}

	records := r.Group("/api/tags/records")
	{
		records.POST("/add", api.TagsRecordsApi{}.AddTagsRecords)
	}

	invitationCodes := r.Group("/api/invitation/codes")
	{
		invitationCodes.GET("/add", api.InvitationCodesApi{}.MakeInvitationCodes)
	}
}

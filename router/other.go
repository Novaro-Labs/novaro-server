package router

import (
	"github.com/gin-gonic/gin"
	"novaro-server/api"
	"novaro-server/utils"
)

func AddOtherRoutes(r *gin.RouterGroup) {
	collections := r.Group("/api/collections")
	{
		collectionsApi := api.CollectionsApi{}

		collections.POST("/add", collectionsApi.CollectionsTweet)
		collections.POST("/remove", collectionsApi.UnCollectionsTweet)
	}

	comments := r.Group("/api/comments")
	{
		commentsApi := api.CommentsApi{}
		comments.GET("/getCommentsListByPostId", commentsApi.GetCommentsListByPostId)
		comments.GET("/getCommentsListByParentId", commentsApi.GetCommentsListByParentId)
		comments.GET("/getCommentsListByUserId", commentsApi.GetCommentsListByUserId)
		comments.POST("/add", commentsApi.AddComments)
		comments.DELETE("/delete", commentsApi.DeleteById)
	}

	posts := r.Group("/api/posts")
	{
		postsApi := api.PostsApi{}
		posts.GET("/getPostsById", postsApi.GetPostsById)
		posts.GET("/getPostsByUserId", postsApi.GetPostsByUserId)
		posts.POST("/getPostsList", postsApi.GetPostsList)
		posts.POST("/savePosts", postsApi.SavePosts)
		posts.POST("/saveRePosts", postsApi.SavePosts)
		posts.DELETE("/delPostsById", postsApi.DelPostsById)
	}

	reposts := r.Group("/api/reposts")
	{
		postsApi := api.RePostsApi{}
		reposts.POST("/add", postsApi.AddRePosts)
	}

	tags := r.Group("/api/tags")
	{
		tagsApi := api.TagsApi{}
		tags.GET("/list", tagsApi.GetTagsList)
	}

	records := r.Group("/api/tags/records")
	{
		recordsApi := api.TagsRecordsApi{}
		records.POST("/add", recordsApi.AddTagsRecords)
	}

	files := r.Group("/upload")
	{
		files.POST("/files", func(context *gin.Context) {
			utils.UploadFiles(context.Writer, context.Request)
		})

	}

	events := r.Group("/api/event")
	{
		eventsApi := api.EventsApi{}
		events.POST("/add", eventsApi.CreateEvents)
		events.DELETE("/delete", eventsApi.DeleteEvents)
		events.PUT("/update", eventsApi.UpdateEvents)
		events.POST("/list", eventsApi.GetList)
	}
}

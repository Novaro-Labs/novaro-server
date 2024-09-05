package routers

import (
	"github.com/gin-gonic/gin"
	"novaro-server/api"
	"novaro-server/utils"
)

func AddOtherRoutes(r *gin.RouterGroup) {

	cron := utils.NewCronManager()

	collections := r.Group("/api/collections")
	{
		collectionsApi := api.NewCollectionsApi()
		// 定时器
		cron.AddJob("@every 5m", func() {
			collectionsApi.Sync()
		})

		collections.POST("/add", collectionsApi.Create)
	}

	comments := r.Group("/api/comments")
	{
		commentsApi := api.NewCommentApi()

		comments.GET("/getCommentsListByPostId", commentsApi.GetCommentsListByPostId)
		comments.GET("/getCommentsListByParentId", commentsApi.GetCommentsListByParentId)
		comments.GET("/getCommentsListByUserId", commentsApi.GetCommentsListByUserId)
		comments.POST("/add", commentsApi.AddComments)
		comments.DELETE("/delete", commentsApi.DeleteById)
	}

	posts := r.Group("/api/posts")
	{
		postsApi := api.NewPostsApi()

		cron.AddJob("@every 3s", func() {
			postsApi.SyncData()
		})

		posts.GET("/getPostsById", postsApi.GetPostsById)
		posts.GET("/getPostsByUserId", postsApi.GetPostsByUserId)
		posts.POST("/getPostsList", postsApi.GetPostsList)
		posts.POST("/savePosts", postsApi.SavePosts)
		posts.POST("/saveRePosts", postsApi.SavePosts)
		posts.DELETE("/delPostsById", postsApi.DelPostsById)
	}

	reposts := r.Group("/api/reposts")
	{
		postsApi := api.NewRePostApi()
		reposts.POST("/add", postsApi.AddRePosts)
	}

	tags := r.Group("/api/tags")
	{
		tagsApi := api.NewTagsApi()
		tags.GET("/list", tagsApi.GetTagsList)
	}

	records := r.Group("/api/tags/records")
	{

		recordsApi := api.NewTagsRecordApi()
		records.POST("/add", recordsApi.AddTagsRecords)
	}
	invitationCodes := r.Group("/api/invitation/codes")
	{
		invitationCodesApi := api.NewInvitationCodesApi()
		invitationCodes.GET("/add", invitationCodesApi.MakeInvitationCodes)
	}

	files := r.Group("/upload")
	{
		files.POST("/files", func(context *gin.Context) {
			utils.UploadFiles(context.Writer, context.Request)
		})

	}

	events := r.Group("/api/event")
	{
		eventsApi := api.NewEventApi()
		events.POST("/add", eventsApi.CreateEvents)
		events.DELETE("/delete", eventsApi.DeleteEvents)
		events.PUT("/update", eventsApi.UpdateEvents)
		events.POST("/list", eventsApi.GetList)
	}

	cron.Start()
}

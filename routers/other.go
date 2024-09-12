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
		collections.POST("/add", collectionsApi.Create)
	}

	comments := r.Group("/api/comments")
	{
		commentsApi := api.NewCommentApi()

		cron.AddJob("@every 30s", func() {
			//commentsApi.SyncCommentsToDB()
		})

		comments.GET("/getCommentsListByPostId", commentsApi.GetCommentsListByPostId)
		comments.GET("/getCommentsListByParentId", commentsApi.GetCommentsListByParentId)
		comments.GET("/getCommentsListByUserId", commentsApi.GetCommentsListByUserId)
		comments.PUT("/add", commentsApi.AddComments)
		comments.DELETE("/delete", commentsApi.DeleteById)
	}

	posts := r.Group("/api/posts")
	{
		postsApi := api.NewPostsApi()
		//posts.GET("/getPostsById", postsApi.GetPostsById)
		//posts.GET("/getPostsByUserId", postsApi.GetPostsByUserId)
		posts.POST("/list", postsApi.GetPostsList)
		posts.PUT("/save", postsApi.SavePosts)
		posts.PUT("/resave", postsApi.SavePosts)
		posts.DELETE("/delete", postsApi.DelPostsById)
	}

	tags := r.Group("/api/tags")
	{
		tagsApi := api.NewTagsApi()
		tags.GET("/list", tagsApi.GetTagsList)
	}

	records := r.Group("/api/tags/records")
	{

		recordsApi := api.NewTagsRecordApi()
		cron.AddJob("@every 30s", func() {
			//recordsApi.SyncData()
		})

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

	nft := r.Group("/api/nft")
	{
		nftApi := api.NewNftInfoApi()
		nft.GET("/info", nftApi.GetNftInfo)
		nft.POST("/updatePoints", nftApi.UpdatePoints)
	}

	r.Group("/api/postPoints")
	{
		//pointsApi := api.NewPostPointsApi()
		cron.AddJob("@every 30s", func() {
			//pointsApi.SyncData()
		})
	}

	pointsHistory := r.Group("/api/pointsHistory")
	{
		historyApi := api.NewPointsHistoryApi()
		pointsHistory.POST("/list", historyApi.GetList)
	}
	cron.Start()
}

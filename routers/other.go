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
		comments.POST("/add", commentsApi.AddComments)
		comments.DELETE("/delete", commentsApi.DeleteById)
	}

	posts := r.Group("/api/posts")
	{
		postsApi := api.NewPostsApi()
		posts.POST("/list", postsApi.GetPostsList)
		posts.POST("/save", postsApi.SavePosts)
		posts.POST("/resave", postsApi.SavePosts)
		posts.DELETE("/delete", postsApi.DelPostsById)
		posts.POST("/listByUser", postsApi.GetListByUser)
		posts.POST("/likes", postsApi.GetLikeByUser)
		posts.POST("/comments", postsApi.GetCommentByUser)
		posts.GET("/get", postsApi.GetPostById)
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

		records.PUT("/add", recordsApi.AddTagsRecords)
	}
	invitationCodes := r.Group("/api/invitation/codes")
	{
		invitationCodesApi := api.NewInvitationCodesApi()
		invitationCodes.GET("/add", invitationCodesApi.MakeInvitationCodes)
	}

	files := r.Group("/upload")
	{
		uploadApi := api.NewUploadApi()
		files.POST("/files", uploadApi.UploadFile)
		files.GET("/novaro", uploadApi.LoadSql)
		files.GET("/getTokenImg", uploadApi.TokenImg)

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
		nft.PUT("/updatePoints", nftApi.UpdatePoints)
	}

	r.Group("/api/postPoints")
	{
		pointsApi := api.NewPostPointsApi()
		cron.AddJob("@every 30s", func() {
			pointsApi.SyncData()
		})
	}

	pointsHistory := r.Group("/api/pointsHistory")
	{
		historyApi := api.NewPointsHistoryApi()
		pointsHistory.GET("/list", historyApi.GetList)
		pointsHistory.GET("/statistics", historyApi.Statistics)
	}

	pointsChangeLog := r.Group("/api/pointsChangeLog")
	{
		changeLogApi := api.NewPointsChangeLogApi()
		pointsChangeLog.POST("/list", changeLogApi.GetList)
	}

	tokens := r.Group("/api/tokens")
	{
		tokensApi := api.NewNftTokensApi()
		tokens.GET("/getTokensByWallet", tokensApi.GetTokensByWallet)
		tokens.POST("/save", tokensApi.SaveNftToken)
	}

	likes := r.Group("/api/likes")
	{
		likesApi := api.NewLikeApi()
		cron.AddJob("@every 30s", func() {
			likesApi.FlushToDatabase()
		})

		likes.POST("/like", likesApi.Like)
	}
	cron.Start()
}

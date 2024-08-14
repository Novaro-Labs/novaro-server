package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.String(http.StatusOK, "NOVARO")
}

func AddHomeRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/home/index")

	group.GET("/", index)
}

package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const Userkey = "user"

func CurrentUser(c *gin.Context) string {
	session := sessions.Default(c)
	user := session.Get(Userkey)
	return user.(string)
}

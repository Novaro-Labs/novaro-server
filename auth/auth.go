package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const Userkey = "user"

func CurrentUser(c *gin.Context) *string {
	session := sessions.Default(c)
	user := session.Get(Userkey)
	if user != nil {
		userId := user.(string)
		return &userId
	} else {
		return nil
	}
}

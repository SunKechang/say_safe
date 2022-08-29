package middleware

import (
	"gin-test/handler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get(handler.Role)
	if role == nil {
		c.Redirect(http.StatusFound, "/page/index")
		c.Abort()
		return
	}

	c.Set(handler.UserName, session.Get(handler.UserName))
	c.Set(handler.Password, session.Get(handler.Password))
	c.Set(handler.Role, session.Get(handler.Role))
}

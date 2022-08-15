package middleware

import (
	"gin-test/handler"
	"gin-test/handler/response"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	role := session.Get(handler.Role)
	if role == nil {
		res := response.NewResponse()
		res["code"] = http.StatusUnauthorized
		res[handler.Message] = "redirect to /login"
		c.JSON(res["code"].(int), res)
		c.Abort()
		return
	}

	c.Set(handler.UserName, session.Get(handler.UserName))
	c.Set(handler.Password, session.Get(handler.Password))
	c.Set(handler.Role, session.Get(handler.Role))
}

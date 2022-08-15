package handler

import (
	"gin-test/handler/response"
	"github.com/gin-gonic/gin"
)

func Hello() func(*gin.Context) {
	return func(context *gin.Context) {
		hello(context)
	}
}
func hello(c *gin.Context) {
	res := response.NewResponse()
	res[Message] = "hello"
	c.JSON(res["code"].(int), res)
}

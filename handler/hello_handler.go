package handler

import (
	"gin-test/handler/response"
	"gin-test/util/flag"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func Hello() func(*gin.Context) {
	return func(context *gin.Context) {
		hello(context)
	}
}
func GetHello() func(*gin.Context) {
	return func(context *gin.Context) {
		getHello(context)
	}
}
func GetPublic() func(c *gin.Context) {
	return func(c *gin.Context) {
		getPublic(c)
	}
}
func hello(c *gin.Context) {
	res := response.NewResponse()
	res[Message] = "hello"
	c.JSON(res["code"].(int), res)
}

func getHello(c *gin.Context) {
	res := response.NewResponse()
	res[Message] = "hello"
	c.JSON(res["code"].(int), res)
}

func getPublic(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)
	file, err := ioutil.ReadFile(flag.PubPath)
	if err != nil {
		return
	}
	res["data"] = string(file)
}

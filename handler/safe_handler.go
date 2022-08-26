package handler

import (
	"fmt"
	safe2 "gin-test/database/safe"
	"gin-test/handler/response"
	"gin-test/service/safe"
	"gin-test/util/log"
	"github.com/gin-gonic/gin"
	"io"
)

func SaySafe() func(*gin.Context) {
	return func(context *gin.Context) {
		saySafe(context)
	}
}

func AddSafe() func(*gin.Context) {
	return func(context *gin.Context) {
		addSafe(context)
	}
}

func GetSafe() func(*gin.Context) {
	return func(context *gin.Context) {
		getSafe(context)
	}
}
func GetSafeList() func(*gin.Context) {
	return func(context *gin.Context) {
		getSafeList(context)
	}
}
func saySafe(c *gin.Context) {
	res := response.NewResponse()
	temp, _ := c.Get(UserName)
	username := temp.(string)
	temp, _ = c.Get(Password)
	password := temp.(string)
	safeService := safe.NewSafeService()
	safeRes, err := safeService.SendSafe(username, password)
	if err != nil {
		res[Message] = err.Error()
	} else {
		res[Message] = safeRes
	}
	c.JSON(res["code"].(int), res)
}

//将删除该用户所有已存在的报平安任务
func addSafe(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)

	temp, _ := c.Get(UserName)
	username := temp.(string)

	//读取request body中数据
	bodyTemp, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Log(fmt.Sprintf("[AddSafe] failed: %s\n", err.Error()))
	}

	service := safe.NewSafeService()
	err = service.AddSafe(username, bodyTemp)
	if err != nil {
		res[Message] = err.Error()
		return
	}
	res[Message] = "添加成功"
}

func getSafe(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)

	temp, _ := c.Get(UserName)
	username := temp.(string)
	service := safe.NewSafeService()
	content, err := service.GetSafe(username)
	if err != nil {
		res[Message] = err.Error()
		return
	}
	res["data"] = content
}

func getSafeList(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)

	temp, _ := c.Get(UserName)
	username := temp.(string)
	dao := safe2.NewSafeLogDao()
	content, err := dao.GetSafeLog(username)
	if err != nil {
		res[Message] = err.Error()
		return
	}
	res["data"] = content
}

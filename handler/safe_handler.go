package handler

import (
	safe2 "gin-test/database/safe"
	"gin-test/handler/response"
	"gin-test/service/safe"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddSafe1() func(*gin.Context) {
	return func(context *gin.Context) {
		addSafe1(context)
	}
}

func SaySafe1() func(*gin.Context) {
	return func(context *gin.Context) {
		saySafe1(context)
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
func getSafe(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)

	temp, _ := c.Get(UserName)
	username := temp.(string)
	service := safe.NewSafeService()
	content, count, err := service.GetSafe(username)
	if err != nil {
		res[Message] = err.Error()
		return
	}
	res["data"] = content
	res["count"] = count
}

func getSafeList(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)

	temp, _ := c.Get(UserName)
	username := temp.(string)

	noTemp := c.DefaultQuery("pageNo", "1")
	sizeTemp := c.DefaultQuery("pageSize", "20")
	pageNo, err := strconv.Atoi(noTemp)
	if err != nil || pageNo < 1 {
		res[Message] = "传入参数格式有误"
		return
	}
	pageSize, err := strconv.Atoi(sizeTemp)
	if err != nil || pageSize < 0 {
		res[Message] = "传入参数格式有误"
		return
	}
	dao := safe2.NewSafeLogDao()
	content, err := dao.GetSafeLog(username, pageNo, pageSize)
	if err != nil {
		res[Message] = err.Error()
		return
	}
	res["data"] = content
}

func addSafe1(c *gin.Context) {
	res := response.NewResponse()
	defer c.JSON(res["code"].(int), res)

	temp, _ := c.Get(UserName)
	username := temp.(string)

	service := safe.NewSafeService()
	err := service.AddSafe1(username)
	if err != nil {
		res[Message] = err.Error()
		return
	}
	res[Message] = "添加成功"
}

func saySafe1(c *gin.Context) {
	res := response.NewResponse()
	temp, _ := c.Get(UserName)
	username := temp.(string)
	temp, _ = c.Get(Password)
	password := temp.(string)
	safeService := safe.NewSafeService()
	safeRes, err := safeService.SendSafe1(username, password)
	if err != nil {
		res[Message] = err.Error()
	} else {
		res[Message] = safeRes
	}
	c.JSON(res["code"].(int), res)
}

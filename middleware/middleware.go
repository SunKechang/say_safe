package middleware

import (
	"github.com/gin-gonic/gin"
)

func InitMiddlewares(r *gin.Engine) {
	var middlewares []gin.HandlerFunc
	middlewares = append(middlewares, func(context *gin.Context) {
		AuthMiddleware(context)
	})
	r.Use(middlewares...)
}

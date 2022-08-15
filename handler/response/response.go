package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewResponse() gin.H {
	return gin.H{"code": http.StatusOK, "message": "", "data": ""}
}

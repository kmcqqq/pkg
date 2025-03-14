package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 返回前端
func BsResponse(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "message": message})
}

// 返回前端-成功
func BsSuccess(c *gin.Context, data gin.H) {
	BsResponse(c, http.StatusOK, 0, data, "请求成功")
}

func BsSuccess1(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data, "message": "请求成功"})
}

// 返回前端-失败
func BsFail(c *gin.Context, data gin.H, message string) {
	BsResponse(c, http.StatusOK, 400, data, message)
}

package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 返回前端
func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "msg": message})
}

// 返回前端-成功
func Success(c *gin.Context, data gin.H) {
	Response(c, http.StatusOK, 100, data, "请求成功")
}

func Success1(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 100, "data": data, "msg": "请求成功"})
}

// 返回前端-失败
func Fail(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 400, data, message)
}

func NoData(c *gin.Context, data gin.H) {
	Response(c, http.StatusOK, 106, data, "暂无数据")
}

func Error(c *gin.Context) {
	Response(c, http.StatusInternalServerError, 500, nil, "Internal Server Error")
}

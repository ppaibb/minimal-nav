package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 API 响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误返回
func Error(c *gin.Context, code int, msg string) {
	httpStatus := http.StatusOK
	if code >= 400 && code < 600 {
		httpStatus = code
	}
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

package common

import "github.com/gin-gonic/gin"

// Response 标准返回结构
type Response struct {
	Code int         `json:"code"` // 自定义状态码
	Data interface{} `json:"data"` // 数据内容
	Msg  string      `json:"msg"`  // 提示信息
}

// Result 统一封装返回函数
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(200, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// Success 成功快捷返回
func Success(data interface{}, msg string, c *gin.Context) {
	Result(200, data, msg, c)
}

// Fail 失败快捷返回
func Fail(code int, msg string, c *gin.Context) {
	Result(code, nil, msg, c)
}

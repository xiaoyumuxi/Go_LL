package controller

import (
	"gin-crud/common"
	"gin-crud/service"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, s *service.UserService) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		common.Fail(400, "参数错误", c)
		return
	}

	// 直接调用 Service 层的 Login
	token, err := s.Login(loginData.Username, loginData.Password)
	if err != nil {
		common.Fail(401, err.Error(), c)
		return
	}

	common.Success(gin.H{"token": token}, "登录成功", c)
}

// AuthMiddleware 拦截器，拦截没有权限的访问
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		token := c.GetHeader("Authorization")

		if token == "" {
			common.Fail(401, "未登录，请先提供 Token", c)
			c.Abort()
			return
		}

		claims, err := common.ParseToken(token)
		if err != nil {
			common.Fail(401, "Token 无效或已过期", c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

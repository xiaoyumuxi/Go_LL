package controller

import (
	"gin-crud/common"
	"gin-crud/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *gin.Context, db *gorm.DB) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		common.Fail(400, "参数错误", c)
		return
	}

	var user models.User
	// 1. 查找用户
	if err := db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		common.Fail(404, "用户不存在", c)
		return
	}

	// 2. 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		common.Fail(401, "密码错误", c)
		return
	}

	// 3. TODO: 生成 JWT Token (这里先用模拟字符串代替)
	token := "valid-token-for-" + user.Username
	common.Success(gin.H{"token": token}, "登录成功", c)
}

// AuthMiddleware 拦截器，拦截没有权限的访问
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 Token
		token := c.GetHeader("Authorization")

		// 简单的 Token 验证逻辑
		if token == "" || token != "valid-token-for-admin" {
			c.JSON(401, gin.H{"msg": "权限不足，请先登录"})
			c.Abort() // 拦截，不许往后走
			return
		}

		c.Next() // 放行
	}
}

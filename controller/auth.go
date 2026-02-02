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

	// 3. 生成token
	token, err := common.GenerateToken(user.ID, user.Username)
	if err != nil {
		common.Fail(500, "生成 Token 失败", c)
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

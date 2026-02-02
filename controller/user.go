package controller

import (
	"gin-crud/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"net/http"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	// 1. 参数绑定与校验 (binding 标签生效)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 2. 存入数据库
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "创建成功", "data": user})
}

// GetUser 获取用户详情
func GetUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id") // 路由参数
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User
	db.Delete(&user, id)
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功", "data": user})
}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User
	db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"username": user.Username,
		"age":      0,
	})
}

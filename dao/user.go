package dao

import (
	"errors"
	"gin-crud/common"
	"gin-crud/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	// 1. 参数绑定与校验 (binding 标签生效)
	if err := c.ShouldBindJSON(&user); err != nil {
		common.Fail(400, "参数校验失败: "+err.Error(), c)
		return
	}
	// 2. 存入数据库
	if err := db.Create(&user).Error; err != nil {
		common.Fail(500, "存储失败", c)
		return
	}
	common.Success(user, "创建成功", c)
}

// GetUser 获取用户详情
func GetUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.Fail(404, "用户不存在", c)
		} else {
			common.Fail(500, "数据库查询异常", c)
		}
		return
	}
	common.Success(user, "获取成功", c)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user models.User
	if err := db.Where("id = ?", id).Delete(&user).Error; err != nil {
		common.Fail(500, "删除失败", c)
		return
	}
	common.Success(nil, "删除成功", c)
}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		common.Fail(400, "无效的 JSON", c)
		return
	}

	result := db.Model(&models.User{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		common.Fail(500, "更新异常", c)
		return
	}
	if result.RowsAffected == 0 {
		common.Fail(404, "未找到该记录或数据无变化", c)
		return
	}

	common.Success(nil, "更新成功", c)
}

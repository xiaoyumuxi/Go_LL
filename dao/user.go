package dao

import (
	"gin-crud/common"
	"gin-crud/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser 创建用户 (暂时保留旧写法，后续可重构)
func CreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		common.Fail(400, "参数校验失败: "+err.Error(), c)
		return
	}
	if err := db.Create(&user).Error; err != nil {
		common.Fail(500, "存储失败", c)
		return
	}
	common.Success(user, "创建成功", c)
}

// GetUserByID 根据 ID 获取用户
func GetUserByID(id string, db *gorm.DB) (*models.User, error) {
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUserByID 根据 ID 删除用户
func DeleteUserByID(id string, db *gorm.DB) error {
	var user models.User
	// GORM 的 Delete 需要传入模型实例或带 Where 条件
	// 这里我们先查一下是否存在，或者直接根据 ID 删除
	// 直接删除更高效，但无法知道 ID 是否存在（除非检查 RowsAffected）
	result := db.Where("id = ?", id).Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateUserByID 根据 ID 更新用户
func UpdateUserByID(id string, updateData map[string]interface{}, db *gorm.DB) error {
	// Updates 方法会自动忽略零值，非常适合 PATCH/PUT 操作
	result := db.Model(&models.User{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// --- 以下旧方法已废弃，待 main.go 彻底移除引用后可删除 ---

// GetUser (旧)
func GetUser(c *gin.Context, db *gorm.DB) {
	// ... (已在 Controller 中实现)
}

// DeleteUser (旧)
func DeleteUser(c *gin.Context, db *gorm.DB) {
	// ... (已在 Controller 中实现)
}

// UpdateUser (旧)
func UpdateUser(c *gin.Context, db *gorm.DB) {
	// ... (已在 Controller 中实现)
}

package service

import (
	"errors"
	"gin-crud/common"
	"gin-crud/dao"
	"gin-crud/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

// Register 注册业务逻辑
func (s *UserService) Register(user *models.User) error {
	var count int64
	s.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}
	return s.DB.Create(user).Error
}

// Login 登录业务逻辑
func (s *UserService) Login(username, password string) (string, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("密码错误")
	}
	return common.GenerateToken(user.ID, user.Username)
}

// GetUser 获取单个用户
func (s *UserService) GetUser(id string) (*models.User, error) {
	user, err := dao.GetUserByID(id, s.DB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id string) error {
	err := dao.DeleteUserByID(id, s.DB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	return nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id string, updateData map[string]interface{}) error {
	// 这里可以添加业务逻辑，例如：
	// 1. 如果更新了密码，需要重新加密
	// 2. 如果更新了用户名，需要检查是否重复

	// 简单示例：如果包含 password 字段，进行加密处理
	if pwd, ok := updateData["password"].(string); ok && pwd != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updateData["password"] = string(hash)
	}

	err := dao.UpdateUserByID(id, updateData, s.DB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	return nil
}

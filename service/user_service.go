package service

import (
	"errors"
	"gin-crud/common"
	"gin-crud/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

// Register 注册业务逻辑
func (s *UserService) Register(user *models.User) error {
	// 1. 业务逻辑：检查用户名是否已占用
	var count int64
	s.DB.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 2. 存入数据库 (注意：密码加密由 models/user.go 的 BeforeSave 钩子自动完成)
	return s.DB.Create(user).Error
}

// Login 登录业务逻辑
func (s *UserService) Login(username, password string) (string, error) {
	var user models.User

	// 1. 查找用户
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", err
	}

	// 2. 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("密码错误")
	}

	// 3. 生成 Token
	return common.GenerateToken(user.ID, user.Username)
}

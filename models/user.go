package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"` // Gin 参数校验
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// BeforeSave 加密逻辑：在创建或更新数据前自动执行
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// 只有当密码字段不为空时才加密（防止更新其他字段时把已加密的密码再次加密）
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

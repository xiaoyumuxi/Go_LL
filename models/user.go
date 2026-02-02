package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"` // Gin 参数校验
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-crud/common"
	"gin-crud/dao"
	"gin-crud/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB  *gorm.DB
	RDB *redis.Client
}

// TokenResponse 登录返回结构
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

// Login 登录业务逻辑 (返回双 Token)
func (s *UserService) Login(username, password string) (*TokenResponse, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 1. 生成 Access Token
	accessToken, err := common.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 2. 生成 Refresh Token
	refreshToken, err := common.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 3. 将 Refresh Token 存入 Redis (有效期 7 天)
	// Key: refresh_token:{token} -> Value: userID
	err = s.RDB.Set(context.Background(), "refresh_token:"+refreshToken, user.ID, 7*24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken 刷新 Access Token
func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	// 1. 查 Redis
	val, err := s.RDB.Get(context.Background(), "refresh_token:"+refreshToken).Result()
	if err == redis.Nil {
		return "", errors.New("Refresh Token 无效或已过期")
	}
	if err != nil {
		return "", err
	}

	// 2. 获取 UserID
	userID, _ := strconv.ParseUint(val, 10, 64)

	// 3. 查用户信息 (确保用户没被封号)
	user, err := dao.GetUserByID(fmt.Sprintf("%d", userID), s.DB)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	// 4. 生成新的 Access Token
	return common.GenerateAccessToken(user.ID, user.Username)
}

// Logout 登出
func (s *UserService) Logout(refreshToken string) error {
	return s.RDB.Del(context.Background(), "refresh_token:"+refreshToken).Err()
}

// GetUser 获取单个用户 (带缓存)
func (s *UserService) GetUser(id string) (*models.User, error) {
	cacheKey := "user:" + id
	val, err := s.RDB.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			common.Logger.Info("Cache Hit: " + id)
			return &user, nil
		}
	}

	common.Logger.Info("Cache Miss: " + id)
	user, err := dao.GetUserByID(id, s.DB)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	jsonBytes, _ := json.Marshal(user)
	s.RDB.Set(context.Background(), cacheKey, jsonBytes, 10*time.Minute)

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
	s.RDB.Del(context.Background(), "user:"+id)
	return nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id string, updateData map[string]interface{}) error {
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
	s.RDB.Del(context.Background(), "user:"+id)
	return nil
}

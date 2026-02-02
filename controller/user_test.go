package controller

import (
	"bytes"
	"encoding/json"
	"gin-crud/common"
	"gin-crud/models"
	"gin-crud/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestApp 初始化测试环境
func setupTestApp() (*gin.Engine, *service.UserService) {
	// 1. 初始化配置（JWT 密钥等依赖它）
	common.InitConfig()

	// 2. 初始化测试数据库
	dsn := "root:105822@tcp(127.0.0.1:3306)/go?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	// 3. 组装 Service
	userService := &service.UserService{DB: db}

	// 4. 初始化路由
	r := gin.Default()
	return r, userService
}

func TestUserWorkflow(t *testing.T) {
	r, userService := setupTestApp()
	db := userService.DB // 测试中可能需要直接操作 DB 清理数据

	// 定义路由：现在全部通过 Service 调用
	r.POST("/register", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			common.Fail(400, "参数错误", c)
			return
		}
		if err := userService.Register(&user); err != nil {
			common.Fail(500, err.Error(), c)
			return
		}
		common.Success(user, "注册成功", c)
	})

	r.POST("/login", func(c *gin.Context) {
		Login(c, userService) // 修改后的 Login 接收 userService
	})

	// 受保护路由
	auth := r.Group("/api")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/profile/:id", func(c *gin.Context) {
			// 这里演示如何在 Controller 直接调用 Service
			id := c.Param("id")
			var user models.User
			if err := db.First(&user, id).Error; err != nil {
				common.Fail(404, "用户不存在", c)
				return
			}
			common.Success(user, "获取成功", c)
		})
	}

	// --- 1. 测试注册 ---
	// 先清理旧数据，保证测试可重复
	db.Exec("DELETE FROM users WHERE username = ?", "tester")

	regData := models.User{Username: "tester", Password: "password123", Email: "test@qq.com"}
	regJson, _ := json.Marshal(regData)
	reqReg, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(regJson))
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	assert.Equal(t, 200, wReg.Code)

	// --- 2. 测试登录并获取 Token ---
	loginData := gin.H{"username": "tester", "password": "password123"}
	loginJson, _ := json.Marshal(loginData)
	reqLog, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJson))
	wLog := httptest.NewRecorder()
	r.ServeHTTP(wLog, reqLog)

	assert.Equal(t, 200, wLog.Code)

	var loginResp map[string]interface{}
	json.Unmarshal(wLog.Body.Bytes(), &loginResp)

	// 健壮性检查：防止 data 为 nil 导致断言失败
	data, ok := loginResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("登录响应格式错误: %v", loginResp)
	}
	token := data["token"].(string)

	// --- 3. 测试带 Token 访问 ---
	// 拿到刚才注册的 ID
	var savedUser models.User
	db.Where("username = ?", "tester").First(&savedUser)

	reqGet, _ := http.NewRequest("GET", "/api/profile/1", nil)
	reqGet.Header.Set("Authorization", token) // 关键

	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)

	assert.Equal(t, 200, wGet.Code)
	assert.Contains(t, wGet.Body.String(), "tester")
}

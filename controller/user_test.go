package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-crud/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化测试用的数据库和路由
func setupTestApp() (*gin.Engine, *gorm.DB) {
	dsn := "root:105822@tcp(127.0.0.1:3306)/go?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	r := gin.Default()
	return r, db
}

func TestUserWorkflow(t *testing.T) {
	r, db := setupTestApp()

	// 定义路由
	r.POST("/register", func(c *gin.Context) { CreateUser(c, db) })
	r.POST("/login", func(c *gin.Context) { Login(c, db) })

	// 受保护路由
	auth := r.Group("/api")
	auth.Use(AuthMiddleware()) // 这里会调用你修改后的 AuthMiddleware
	{
		auth.GET("/profile/:id", func(c *gin.Context) { GetUser(c, db) })
	}

	// --- 1. 测试注册 ---
	regData := models.User{Username: "tester", Password: "password123", Email: "test@qq.com"}
	regJson, _ := json.Marshal(regData)
	reqReg, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(regJson))
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	assert.Equal(t, 200, wReg.Code, "注册应该成功")

	// --- 2. 测试登录并获取 Token ---
	loginData := gin.H{"username": "tester", "password": "password123"}
	loginJson, _ := json.Marshal(loginData)
	reqLog, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJson))
	wLog := httptest.NewRecorder()
	r.ServeHTTP(wLog, reqLog)

	assert.Equal(t, 200, wLog.Code, "登录应该成功")

	var loginResp map[string]interface{}
	json.Unmarshal(wLog.Body.Bytes(), &loginResp)
	// 注意这里：从统一返回格式 data 字段里取 token
	data := loginResp["data"].(map[string]interface{})
	token := data["token"].(string)

	// --- 3. 测试带 Token 访问受保护接口 ---
	reqGet, _ := http.NewRequest("GET", "/api/profile/1", nil)
	// 关键：在 Header 中注入 Token
	reqGet.Header.Set("Authorization", token)

	wGet := httptest.NewRecorder()
	r.ServeHTTP(wGet, reqGet)

	assert.Equal(t, 200, wGet.Code, "带 Token 访问应该成功")
	assert.Contains(t, wGet.Body.String(), "tester", "返回内容应包含用户名")
}

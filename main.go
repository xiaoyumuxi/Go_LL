package main

import (
	"gin-crud/common"
	"gin-crud/controller"
	"gin-crud/models"
	"gin-crud/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "gin-crud/docs" // 必须导入生成的 docs 包
)

// @title           Gin CRUD API
// @version         1.0
// @description     这是一个基于 Gin 的 CRUD 示例项目
// @host            localhost:8080
// @BasePath        /
func main() {
	common.InitConfig()
	common.InitLogger()        // 初始化日志
	defer common.Logger.Sync() // 刷新缓冲

	common.InitDB()    // 初始化数据库
	common.InitRedis() // 初始化 Redis

	// 使用自定义的 Logger 和 Recovery
	r := gin.New()
	r.Use(common.GinLogger(), common.GinRecovery(true))

	// 注入 DB 和 Redis
	userService := &service.UserService{
		DB:  common.DB,
		RDB: common.RDB,
	}

	// Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 配置热更新测试接口
	r.GET("/config-test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"jwt_expire": common.Conf.Jwt.Expire,
			"port":       common.Conf.Server.Port,
		})
	})

	// 公开接口
	r.POST("/login", func(c *gin.Context) {
		controller.Login(c, userService)
	})
	r.POST("/refresh", func(c *gin.Context) {
		controller.RefreshToken(c, userService)
	})
	r.POST("/logout", func(c *gin.Context) {
		controller.Logout(c, userService)
	})

	r.POST("/register", func(c *gin.Context) {
		// 你也可以在这里把 dao.CreateUser 重构进 service
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
	// 路由分组1
	userGroup := r.Group("/users")
	{
		userGroup.GET("/:id", func(c *gin.Context) {
			controller.GetUser(c, userService)
		})
		userGroup.PUT("/:id", func(c *gin.Context) {
			controller.UpdateUser(c, userService)
		})
		userGroup.DELETE("/:id", func(c *gin.Context) {
			controller.DeleteUser(c, userService)
		})
	}

	r.Run(":8080")
}

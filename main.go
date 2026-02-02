package main

import (
	"gin-crud/common"
	"gin-crud/controller"
	"gin-crud/dao"
	"gin-crud/models"
	"gin-crud/service"

	"github.com/gin-gonic/gin"
)

func main() {
	common.InitConfig()

	common.InitDB() // 初始化数据库
	r := gin.Default()

	userService := &service.UserService{DB: common.DB}
	// 公开接口
	r.POST("/login", func(c *gin.Context) {
		controller.Login(c, userService)
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
			dao.GetUser(c, common.DB)
		})
		userGroup.PUT("/:id", func(c *gin.Context) {
			dao.UpdateUser(c, common.DB)
		})
		userGroup.DELETE("/:id", func(c *gin.Context) {
			dao.DeleteUser(c, common.DB)
		})
	}

	r.Run(":8080")
}

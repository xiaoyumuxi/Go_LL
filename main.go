package main

import (
	"gin-crud/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB() // 初始化数据库
	r := gin.Default()

	// 路由分组
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", func(c *gin.Context) {
			controller.CreateUser(c, DB)
		})
		userGroup.GET("/:id", func(c *gin.Context) {
			controller.GetUser(c, DB)
		})
		userGroup.PUT("/:id", func(c *gin.Context) {
			controller.UpdateUser(c, DB)
		})
		userGroup.DELETE("/:id", func(c *gin.Context) {
			controller.DeleteUser(c, DB)
		})
	}

	r.Run(":8080")
}

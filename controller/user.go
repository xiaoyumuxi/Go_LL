package controller

import (
	"gin-crud/common"
	"gin-crud/service"

	"github.com/gin-gonic/gin"
)

// GetUser 处理获取用户详情的请求
func GetUser(c *gin.Context, s *service.UserService) {
	id := c.Param("id")

	user, err := s.GetUser(id)
	if err != nil {
		if err.Error() == "用户不存在" {
			common.Fail(404, err.Error(), c)
		} else {
			common.Fail(500, "系统异常: "+err.Error(), c)
		}
		return
	}

	common.Success(user, "获取成功", c)
}

// DeleteUser 处理删除用户的请求
func DeleteUser(c *gin.Context, s *service.UserService) {
	id := c.Param("id")

	err := s.DeleteUser(id)
	if err != nil {
		if err.Error() == "用户不存在" {
			common.Fail(404, err.Error(), c)
		} else {
			common.Fail(500, "删除失败: "+err.Error(), c)
		}
		return
	}

	common.Success(nil, "删除成功", c)
}

// UpdateUser 处理更新用户的请求
func UpdateUser(c *gin.Context, s *service.UserService) {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		common.Fail(400, "无效的 JSON: "+err.Error(), c)
		return
	}

	err := s.UpdateUser(id, updateData)
	if err != nil {
		if err.Error() == "用户不存在" {
			common.Fail(404, err.Error(), c)
		} else {
			common.Fail(500, "更新失败: "+err.Error(), c)
		}
		return
	}

	common.Success(nil, "更新成功", c)
}

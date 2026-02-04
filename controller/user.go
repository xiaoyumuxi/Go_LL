package controller

import (
	"gin-crud/common"
	"gin-crud/service"

	"github.com/gin-gonic/gin"
)

// GetUser 获取用户详情
// @Summary      获取用户详情
// @Description  根据 ID 获取用户信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  common.Response{data=models.User}
// @Failure      404  {object}  common.Response
// @Failure      500  {object}  common.Response
// @Router       /users/{id} [get]
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

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  根据 ID 删除用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  common.Response
// @Failure      404  {object}  common.Response
// @Failure      500  {object}  common.Response
// @Router       /users/{id} [delete]
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

// UpdateUser 更新用户
// @Summary      更新用户
// @Description  根据 ID 更新用户信息
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "User ID"
// @Param        data  body      map[string]interface{} true  "Update Data"
// @Success      200   {object}  common.Response
// @Failure      400   {object}  common.Response
// @Failure      404   {object}  common.Response
// @Failure      500   {object}  common.Response
// @Router       /users/{id} [put]
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

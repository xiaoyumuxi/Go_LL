package controller

import (
	"gin-crud/common"
	"gin-crud/service"

	"github.com/gin-gonic/gin"
)

// Login 登录接口
// @Summary      用户登录
// @Description  使用用户名和密码登录，返回 Access Token 和 Refresh Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      object{username=string,password=string}  true  "Login Data"
// @Success      200   {object}  common.Response{data=service.TokenResponse}
// @Failure      400   {object}  common.Response
// @Failure      401   {object}  common.Response
// @Router       /login [post]
func Login(c *gin.Context, s *service.UserService) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		common.Fail(400, "参数错误", c)
		return
	}

	tokens, err := s.Login(loginData.Username, loginData.Password)
	if err != nil {
		common.Fail(401, err.Error(), c)
		return
	}

	common.Success(tokens, "登录成功", c)
}

// RefreshToken 刷新 Token 接口
// @Summary      刷新 Access Token
// @Description  使用 Refresh Token 换取新的 Access Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      object{refresh_token=string}  true  "Refresh Token"
// @Success      200   {object}  common.Response{data=object{access_token=string}}
// @Failure      400   {object}  common.Response
// @Failure      401   {object}  common.Response
// @Router       /refresh [post]
func RefreshToken(c *gin.Context, s *service.UserService) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(400, "参数错误", c)
		return
	}

	newAccessToken, err := s.RefreshToken(req.RefreshToken)
	if err != nil {
		common.Fail(401, err.Error(), c)
		return
	}

	common.Success(gin.H{"access_token": newAccessToken}, "刷新成功", c)
}

// Logout 登出接口
// @Summary      用户登出
// @Description  使 Refresh Token 失效
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      object{refresh_token=string}  true  "Refresh Token"
// @Success      200   {object}  common.Response
// @Failure      400   {object}  common.Response
// @Router       /logout [post]
func Logout(c *gin.Context, s *service.UserService) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(400, "参数错误", c)
		return
	}

	if err := s.Logout(req.RefreshToken); err != nil {
		// 即使删除失败（比如 key 不存在），通常也返回成功，避免泄露信息
		common.Logger.Error("Logout failed: " + err.Error())
	}

	common.Success(nil, "登出成功", c)
}

// AuthMiddleware 拦截器
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			common.Fail(401, "未登录，请先提供 Token", c)
			c.Abort()
			return
		}

		claims, err := common.ParseToken(token)
		if err != nil {
			common.Fail(401, "Token 无效或已过期", c)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

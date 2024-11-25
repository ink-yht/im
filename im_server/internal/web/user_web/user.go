package user_web

import (
	"github.com/gin-gonic/gin"
)

// 确保 UserHandler 上实现了 Handler 接口
//var _ web.Handler = (*UserHandler)(nil)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// RegisterRoutes 路由注册
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users", u.SignUp)       // 用户注册
	server.POST("/login", u.Login)        // 用户登录
	server.PUT("/users/:id", u.Edit)      // 用户修改个人信息
	server.GET("/users/:id", u.Info)      // 用户信息获取
	server.DELETE("/users/:id", u.Logout) // 用户注销
}

func (u *UserHandler) SignUp(ctx *gin.Context) {

}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Info(ctx *gin.Context) {

}

func (u *UserHandler) Logout(ctx *gin.Context) {

}

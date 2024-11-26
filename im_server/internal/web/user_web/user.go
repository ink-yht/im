package user_web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/service/user_service"
	"github.com/ink-yht/im/internal/web"
	"net/http"
)

// 确保 UserHandler 上实现了 Handler 接口
//var _ web.Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc user_service.UserService
}

func NewUserHandler(svc user_service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
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
	var req user_domain.UserRegisterRequest
	err := ctx.Bind(&req)
	if err != nil {
		// 出错会返回 400 错误
		return
	}

	err = u.svc.Signup(ctx, user_domain.UserRegisterRequest(req))

	if errors.Is(err, user_domain.ErrTheMailboxIsNotInTheRightFormat) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "电子邮件格式无效",
			Data: nil,
		})
		return
	}
	if errors.Is(err, user_domain.ErrThePasswordIsNotInTheRightFormat) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符",
			Data: nil,
		})
		return
	}
	if errors.Is(err, user_domain.ErrThePasswordIsInconsistentTwice) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "两次密码不一致",
			Data: nil,
		})
		return
	}
	if errors.Is(err, user_service.ErrDuplicate) {
		// 邮箱冲突
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "邮箱冲突",
			Data: nil,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "注册成功",
		Data: nil,
	})

}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Info(ctx *gin.Context) {

}

func (u *UserHandler) Logout(ctx *gin.Context) {

}

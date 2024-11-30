package user_web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/service/user_service"
	"github.com/ink-yht/im/internal/web"
	"github.com/ink-yht/im/pkg/logger"
	"net/http"
)

//// 确保 UserHandler 上实现了 Handler 接口
//var _ web.Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc user_service.UserService
	l   logger.Logger
}

func NewUserHandler(svc user_service.UserService, l logger.Logger) *UserHandler {
	return &UserHandler{
		svc: svc,
		l:   l,
	}
}

// RegisterRoutes 路由注册
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp) // 用户注册
	ug.POST("/login", u.Login)   // 用户登录
	ug.POST("/edit", u.Edit)     // 用户修改个人信息
	ug.GET("/info", u.Info)      // 用户信息获取
	ug.GET("/logout", u.Logout)  // 用户注销
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var req user_domain.EmailRegisterRequest
	err := ctx.Bind(&req)
	if err != nil {
		// 出错会返回 400 错误
		return
	}

	err = u.svc.Signup(ctx, req)

	if errors.Is(err, user_domain.ErrTheMailboxIsNotInTheRightFormat) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "电子邮件格式无效",
			Data: nil,
		})
		u.l.Warn("电子邮件格式无效", logger.String("email", req.Email))
		return
	}
	if errors.Is(err, user_domain.ErrThePasswordIsNotInTheRightFormat) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符",
			Data: nil,
		})
		u.l.Warn("密码格式不对", logger.String("email", req.Email))
		return
	}
	if errors.Is(err, user_domain.ErrThePasswordIsInconsistentTwice) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "两次密码不一致",
			Data: nil,
		})
		u.l.Warn("两次密码不一致", logger.String("email", req.Email))
		return
	}
	if errors.Is(err, user_service.ErrDuplicate) {
		// 邮箱冲突
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "邮箱冲突",
			Data: nil,
		})
		u.l.Warn("邮箱冲突", logger.String("email", req.Email))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("系统错误", logger.Error("err", err))
		return
	}

	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "注册成功",
		Data: nil,
	})
	u.l.Info("登录成功", logger.String("email", req.Email))
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var req user_domain.EmailLoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		return
	}

	// 提前提取 User-Agent 头
	userAgent := ctx.GetHeader("User-Agent")

	token, err := u.svc.Login(ctx, &req, userAgent)
	if errors.Is(err, user_service.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "邮箱或密码不存在",
			Data: nil,
		})
		u.l.Warn("邮箱或密码不存在", logger.String("email", req.Email))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		u.l.Error("邮箱或密码不存在", logger.Error("err", err))
		return
	}
	ctx.Header("x-jwt-token", token)
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "登录成功",
		Data: nil,
	})
	u.l.Info("登录成功", logger.String("email", req.Email))
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	userClaims := ctx.MustGet("claims").(*user_service.UserClaims)
	var req user_domain.UpdateInfoRequest
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	req.ID = userClaims.Id
	err = u.svc.Edit(ctx, req)
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
		Msg:  "个人信息修改成功",
		Data: nil,
	})
}

func (u *UserHandler) Info(ctx *gin.Context) {

	userClaims := ctx.MustGet("claims").(*user_service.UserClaims)

	user, err := u.svc.Info(ctx, userClaims.Id)
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
		Msg:  "个人信息获取成功",
		Data: user,
	})
}

func (u *UserHandler) Logout(ctx *gin.Context) {

}

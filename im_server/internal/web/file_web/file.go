package file_web

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ink-yht/im/internal/service/file_service"
	"github.com/ink-yht/im/internal/service/user_service"
	"github.com/ink-yht/im/internal/web"
	"github.com/ink-yht/im/pkg/logger"
	"net/http"
)

type FileHandler struct {
	svc file_service.FileService
	l   logger.Logger
}

func NewFileHandler(svc file_service.FileService, l logger.Logger) *FileHandler {
	return &FileHandler{
		svc: svc,
		l:   l,
	}
}

// RegisterRoutes 路由注册
func (f *FileHandler) RegisterRoutes(server *gin.Engine) {
	fg := server.Group("/files")
	fg.POST("/avatar", f.Avatar)
}

func (f *FileHandler) Avatar(ctx *gin.Context) {

	userClaims := ctx.MustGet("claims").(*user_service.UserClaims)
	imageType := "avatar"

	form, err := ctx.MultipartForm()
	if err != nil {
		f.l.Error("上传图片错误")
		return
	}

	// 获取单个文件，这里假设表单中文件字段名为"image"
	file, ok := form.File["image"]
	if !ok {
		errors.New("上传图片错误，未找到上传文件")
		return
	}

	// 取到单个文件对象
	image := file[0]
	err = f.svc.Avatar(ctx, userClaims.Id, imageType, image)
	if errors.Is(err, file_service.ErrSavePicture) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "保存图片失败",
			Data: nil,
		})
		f.l.Warn("保存图片失败")
		return
	}
	if errors.Is(err, file_service.ErrImageSize) {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 1,
			Msg:  "图片大小超过设定大小，当设定大小为: 2Mb ",
			Data: nil,
		})
		f.l.Warn("图片大小超过设定大小，当设定大小为: 2Mb ")
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, web.Result{
			Code: 2,
			Msg:  "系统错误",
			Data: nil,
		})
		f.l.Error("系统错误")
		return
	}
	ctx.JSON(http.StatusOK, web.Result{
		Code: 0,
		Msg:  "上传头像成功",
		Data: nil,
	})
	f.l.Info("上传头像成功")
}

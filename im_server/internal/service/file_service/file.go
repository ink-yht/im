package file_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/repository/file_repo"
	"github.com/spf13/viper"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	ErrSavePicture = errors.New("保存图片失败")
	ErrImageSize   = errors.New("图片大小超过设定大小，当设定大小为: 2Mb ")
)

// FileService 定义了用户服务的接口
type FileService interface {
	Avatar(ctx context.Context, id int64, imageType string, file *multipart.FileHeader) error
}

// FileServiceImpl 实现了 UserService 接口
type FileServiceImpl struct {
	repo file_repo.FileRepository
}

func NewFileService(repo file_repo.FileRepository) FileService {
	return &FileServiceImpl{
		repo: repo,
	}
}

func (svc FileServiceImpl) Avatar(ctx context.Context, id int64, imageType string, image *multipart.FileHeader) error {
	type Config struct {
		Size int    `yaml:"size"`
		Path string `yaml:"path"`
	}
	var c Config
	err := viper.UnmarshalKey("uploads", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败: %s \n", err))
	}

	// 判断图片路径是否存在，没有就创建
	basePath := c.Path + imageType
	_, err = os.ReadDir(basePath)
	if err != nil {
		err = os.MkdirAll(basePath, os.ModePerm)
		if err != nil {
			return errors.New("图片路径创建失败")
		}
	}

	// 判断大小
	size := float64(image.Size) / float64(1024*1024)
	if size >= float64(c.Size) {
		return ErrImageSize
	}

	// 保存文件
	filePath := filepath.Join(basePath, image.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("创建文件失败: %v", err)
		return errors.New("创建文件失败")
	}
	defer out.Close()

	in, err := image.Open()
	if err != nil {
		log.Printf("打开文件失败: %v", err)
		return errors.New("打开文件失败")
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Printf("保存图片失败: %v", err)
		return ErrSavePicture
	}

	// 更新用户头像
	url := "/" + filepath.ToSlash(filePath)
	data := user_domain.User{
		ID:     id,
		Avatar: url,
	}
	err = svc.repo.Avatar(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

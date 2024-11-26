package user_dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicate      = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

// FindByEmail 查询邮箱
func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ? ", email).First(&user).Error
	return user, err
}

// Insert 注册
func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	// 写入数据库

	avatarPath := "/uploads/avatar/logo.png" // 默认头像路径

	// 毫秒
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	u.Avatar = avatarPath
	u.UserConf.CreateTime = now
	u.UserConf.UpdateTime = now

	err := dao.db.WithContext(ctx).Preload("UserConf").Create(&u).Error

	// 如果错误是MySQL错误类型
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 邮箱冲突
			return ErrDuplicate
		}
	}

	// 系统错误
	return err
}

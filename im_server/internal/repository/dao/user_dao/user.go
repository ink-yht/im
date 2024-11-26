package user_dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrDuplicate = errors.New("邮箱冲突")

type UserDao interface {
	Insert(ctx context.Context, u User) error
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

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

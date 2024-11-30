package user_dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	ErrDuplicate      = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	UpdateInfo(ctx context.Context, u User) error
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDao {
	return &GormUserDAO{db: db}
}

func (dao *GormUserDAO) UpdateInfo(ctx context.Context, u User) error {
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新 User 表
		now := time.Now().UnixMilli()
		userUpdateData := map[string]interface{}{
			"update_time": now,
			"nickname":    u.Nickname,
			"signature":   u.Signature,
			"avatar":      u.Avatar,
			"address":     u.Address,
			"birthday":    u.Birthday,
			"sex":         u.Sex,
		}
		if err := tx.Model(&u).Updates(userUpdateData).Error; err != nil {
			return err
		}

		// 更新 UserConf 表
		userConfUpdateData := map[string]interface{}{
			"update_time":           now,
			"recall_message":        u.UserConf.RecallMessage,
			"friend_online":         u.UserConf.FriendOnline,
			"sound":                 u.UserConf.Sound,
			"secure_link":           u.UserConf.SecureLink,
			"save_pwd":              u.UserConf.SavePwd,
			"search_user":           u.UserConf.SearchUser,
			"verification":          u.UserConf.Verification,
			"verification_question": u.UserConf.VerificationQuestion,
			"online":                u.UserConf.Online,
		}
		if err := tx.Model(&u.UserConf).Updates(userConfUpdateData).Error; err != nil {
			return err
		}

		return nil
	})

	// 错误处理
	if err != nil {
		// 记录日志或处理事务错误
		log.Printf("failed to update user info: %v", err)
		return err
	}

	// 如果事务成功
	log.Println("user info updated successfully")
	return nil

}

func (dao *GormUserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Preload("UserConf").Where("id = ? ", id).First(&user).Error
	return user, err
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

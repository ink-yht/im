package dao

import (
	"github.com/ink-yht/im/internal/repository/dao/chat_dao"
	"github.com/ink-yht/im/internal/repository/dao/group_dao"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&user_dao.User{},          // 用户表
		&user_dao.Friend{},        // 好友表
		&user_dao.FriendRequest{}, // 好友验证表
		&user_dao.UserConf{},      // 用户配置表

		&group_dao.Group{},       // 群信息表
		&group_dao.GroupMember{}, // 群成员表
		&group_dao.GroupVerify{}, // 群验证表

		&chat_dao.Chat{},     // 用户消息表
		&chat_dao.GroupMsg{}, // 群消息表
	)
}

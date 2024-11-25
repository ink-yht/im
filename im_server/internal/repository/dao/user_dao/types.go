package user_dao

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/ink-yht/im/internal/repository/dao/chat_dao"
)

// User 用户表
type User struct {
	ID         int64          `gorm:"primaryKey;autoIncrement"` // ID
	Email      sql.NullString `gorm:"unique"`                   // 邮箱，使用唯一键
	Phone      sql.NullString `gorm:"unique"`                   // 手机号，使用唯一键
	Password   string         `gorm:"size:128"`                 // 加密密码
	Nickname   string         `gorm:"size:32"`                  // 昵称
	Signature  string         `gorm:"size:128"`                 // 个性签名
	Avatar     string         `gorm:"size:255"`                 // 头像
	Address    string         `gorm:"size:255"`                 // 地址
	Birthday   int64          //生日
	Sex        int8           `gorm:"default:0"` // 性别 0 未选择 1 男，2 女
	CreateTime int64          // 创建时间
	UpdateTime int64          // 更新时间

	UserConf         UserConf        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 关联的用户配置
	MessagesSent     []chat_dao.Chat `gorm:"foreignKey:SendUserID"`                                           // 发送的消息
	MessagesReceived []chat_dao.Chat `gorm:"foreignKey:RevUserID"`                                            // 接收的消息
}

// UserConf 用户配置表
type UserConf struct {
	ID                   int64                 `gorm:"primaryKey;autoIncrement"` // ID
	CreateTime           int64                 // 创建时间
	UpdateTime           int64                 // 更新时间
	RecallMessage        *string               `gorm:"size:32"` // 撤回消息的提示内容
	FriendOnline         bool                  // 好友上线提醒
	Sound                bool                  // 声音
	SecureLink           bool                  // 安全链接
	SavePwd              bool                  // 保存密码
	SearchUser           int8                  `gorm:"default:2"` // 别人查找到你的方式 0 不允许别人查找到我， 1  通过用户号找到我 2 可以通过手机号搜索到我
	Verification         int8                  `gorm:"default:1"` // 好友验证 0 不允许任何人添加  1 允许任何人添加  2 需要验证消息 3 需要回答问题  4  需要正确回答问题
	VerificationQuestion *VerificationQuestion // 验证问题  为3和4的时候需要
	Online               bool                  // 是否在线

	UserID int64 `gorm:"uniqueIndex"` // 用户ID，唯一键
}

// Friend 好友表
type Friend struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"` // 唯一 ID
	UserID     int64  `gorm:"not null;index"`           // 用户 ID
	FriendID   int64  `gorm:"not null;index"`           // 好友 ID
	Remark     string `gorm:"size:255"`                 // 用户对好友的备注
	Status     int8   `gorm:"default:1"`                // 好友状态 （1: 正常, 2: 拉黑）
	CreateTime int64  `gorm:"autoCreateTime"`           // 创建时间
	UpdateTime int64  `gorm:"autoUpdateTime"`           // 更新时间

	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // 发起者用户外键
	Friend User `gorm:"foreignKey:FriendID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 目标好友外键
}

// FriendRequest 好友请求表
type FriendRequest struct {
	ID                   int64                 `gorm:"primaryKey;autoIncrement"` // ID
	RequesterID          int64                 `gorm:"not null;index"`           // 发起者用户 ID
	ReceiverID           int64                 `gorm:"not null;index"`           // 接收者用户 ID
	ValidationType       int                   `gorm:"not null"`                 // 验证类型
	ValidationMessage    string                `gorm:"type:text"`                // 验证消息
	ValidationAnswer     string                `gorm:"type:text"`                // 提交的答案
	Status               int                   `gorm:"default:0"`                // 状态（0: 未操作, 1: 同意, 2: 拒绝, 3: 忽略）
	AdditionalMessages   string                `gorm:"size:128"`                 // 附加消息
	VerificationQuestion *VerificationQuestion // 验证问题  为3和4的时候需要
	CreateTime           int64                 // 创建时间
	UpdateTime           int64                 // 更新时间

	Requester User `gorm:"foreignKey:RequesterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 请求者
	Receiver  User `gorm:"foreignKey:ReceiverID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`  // 接收者
}

// VerificationQuestion 验证问题
type VerificationQuestion struct {
	Problem1 *string `json:"problem1"`
	Problem2 *string `json:"problem2"`
	Problem3 *string `json:"problem3"`
	Answer1  *string `json:"answer1"`
	Answer2  *string `json:"answer2"`
	Answer3  *string `json:"answer3"`
}

// Scan 取出来的时候的数据
func (c *VerificationQuestion) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

// Value 入库的数据
func (c *VerificationQuestion) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

// 映射实现

// 设置性别

// GetSexText 响应给前端：在返回数据时，将数值转为文本
func GetSexText(sex int8) string {
	switch sex {
	case 0:
		return "未选择"
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "未知"
	}
}

// GetSexValue 后端接受处理时解析为数值存储
func GetSexValue(sexText string) int8 {
	switch sexText {
	case "未选择":
		return 0
	case "男":
		return 1
	case "女":
		return 2
	default:
		return 99 // 未知或未指定
	}
}

// 别人查找到你的方式

// GetSearchUserText 响应给前端：在返回数据时，将数值转为文本
func GetSearchUserText(searchUser int8) string {
	switch searchUser {
	case 0:
		return "不允许别人查找到我"
	case 1:
		return "通过用户号找到我"
	case 2:
		return "可以通过手机号搜索到我"
	default:
		return "未知"
	}
}

// GetSearchUserValue 后端接受处理时解析为数值存储
func GetSearchUserValue(searchUser string) int8 {
	switch searchUser {
	case "不允许别人查找到我":
		return 0
	case "通过用户号找到我":
		return 1
	case "可以通过手机号搜索到我":
		return 2
	default:
		return 99 // 未知或未指定
	}
}

// 好友验证

// GetVerificationText 响应给前端：在返回数据时，将数值转为文本
func GetVerificationText(verification int8) string {
	switch verification {
	case 0:
		return "不允许任何人添加"
	case 1:
		return "允许任何人添加"
	case 2:
		return "需要验证消息"
	case 3:
		return "需要回答问题"
	case 4:
		return "需要正确回答问题"
	default:
		return "未知"
	}
}

// GetVerificationValue 后端接受处理时解析为数值存储
func GetVerificationValue(verification string) int8 {
	switch verification {
	case "不允许任何人添加":
		return 0
	case "允许任何人添加":
		return 1
	case "需要验证消息":
		return 2
	case "需要回答问题":
		return 3
	case "需要正确回答问题":
		return 4
	default:
		return 99 // 未知或未指定
	}
}

// 好友状态

// GetFriendStatusText 响应给前端：在返回数据时，将数值转为文本
func GetFriendStatusText(status int8) string {
	switch status {
	case 1:
		return "正常"
	case 2:
		return "拉黑"
	default:
		return "未知"
	}
}

// GetFriendStatusValue 后端接受处理时解析为数值存储
func GetFriendStatusValue(searchUser string) int8 {
	switch searchUser {
	case "正常":
		return 1
	case "拉黑":
		return 2
	default:
		return 99 // 未知或未指定
	}
}

// 好友请求状态

// GetRequestStatusText 响应给前端：在返回数据时，将数值转为文本
func GetRequestStatusText(Status int8) string {
	switch Status {
	case 0:
		return "未操作"
	case 1:
		return "同意"
	case 2:
		return "拒绝"
	case 3:
		return "忽略"
	default:
		return "未知"
	}
}

// GetRequestStatusValue 后端接受处理时解析为数值存储
func GetRequestStatusValue(Status string) int8 {
	switch Status {
	case "未操作":
		return 0
	case "同意":
		return 1
	case "拒绝":
		return 2
	case "忽略":
		return 3
	default:
		return 99 // 未知或未指定
	}
}

package group_dao

import (
	"database/sql/driver"
	"encoding/json"
)

// Group 群信息表
type Group struct {
	ID                   int64                 `gorm:"primaryKey;autoIncrement"` // ID
	CreateTime           int64                 // 创建时间
	UpdateTime           int64                 // 更新时间
	Title                string                `gorm:"size:32"`  // 群名
	Abstract             string                `gorm:"size:128"` // 简介
	Avatar               string                `gorm:"size:256"` // 群头像
	IsSearch             bool                  // 是否可以被搜索
	Verification         int8                  // 群验证规则
	VerificationQuestion *VerificationQuestion `gorm:"type:json"` // 验证问题
	IsInvite             bool                  // 是否可邀请好友
	IsTemporarySession   bool                  // 是否开启临时会话
	IsProhibition        bool                  // 是否开启全员禁言
	Size                 int                   // 群规模

	Creator int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 群主ID，指向用户表
}

// GroupVerify 群验证表
type GroupVerify struct {
	ID                   int64                 `gorm:"primaryKey;autoIncrement"` // ID
	CreateTime           int64                 // 创建时间
	UpdateTime           int64                 // 更新时间
	GroupModel           Group                 `gorm:"foreignKey:GroupID"` // 群
	Status               int8                  // 验证状态
	AdditionalMessages   string                `gorm:"size:32"`   // 附加消息
	VerificationQuestion *VerificationQuestion `gorm:"type:json"` // 验证问题
	Type                 int8                  // 验证类型

	GroupID int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 群ID
	UserID  int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 用户ID
}

// GroupMember 群成员表
type GroupMember struct {
	ID              int64  `gorm:"primaryKey;autoIncrement"` // ID
	CreateTime      int64  // 创建时间
	UpdateTime      int64  // 更新时间
	GroupModel      Group  `gorm:"foreignKey:GroupID"` // 群
	MemberNickname  string `gorm:"size:32"`            // 群成员昵称
	Role            int    // 成员角色
	ProhibitionTime int64  // 禁言时间（单位：分钟，0 表示未禁言）

	GroupID int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 群ID
	UserID  int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 用户ID
}

// VerificationQuestion 验证问题
type VerificationQuestion struct {
	Problem1 *string `json:"problem1"` // 验证问题 1
	Problem2 *string `json:"problem2"` // 验证问题 2
	Problem3 *string `json:"problem3"` // 验证问题 3
	Answer1  *string `json:"answer1"`  // 答案 1
	Answer2  *string `json:"answer2"`  // 答案 2
	Answer3  *string `json:"answer3"`  // 答案 3
}

func (c *VerificationQuestion) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}

func (c *VerificationQuestion) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

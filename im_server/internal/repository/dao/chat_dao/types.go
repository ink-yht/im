package chat_dao

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Chat 用户消息表
type Chat struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"` // ID
	CreateTime int64  // 创建时间
	UpdateTime int64  // 更新时间
	MsgType    int8   // 消息类型 1 文本类型  2 图片消息  3 视频消息 4 文件消息 5 语音消息  6 语言通话  7 视频通话  8 撤回消息 9 回复消息 10 引用消息
	MsgPreview string `gorm:"size:64"`   // 消息预览
	Msg        Msg    `gorm:"type:json"` // 消息内容

	SendUserID int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 发送者用户ID
	RevUserID  int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 接收者用户ID
}

// GroupMsg 群消息表
type GroupMsg struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"` // ID
	CreateTime int64  // 创建时间
	UpdateTime int64  // 更新时间
	MsgType    int8   // 消息类型 1 文本类型  2 图片消息  3 视频消息 4 文件消息 5 语音消息  6 语言通话  7 视频通话  8 撤回消息 9回复消息 10 引用消息
	MsgPreview string `gorm:"size:64"`   // 消息预览
	Msg        Msg    `gorm:"type:json"` // 消息内容

	GroupID    int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 群ID
	SendUserID int64 `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 发送者用户ID
}

// Msg 消息表
type Msg struct {
	Type         int8          `json:"type"`         // 消息类型
	Content      *string       `json:"content"`      // 文本消息内容
	ImageMsg     *ImageMsg     `json:"imageMsg"`     // 图片消息
	VideoMsg     *VideoMsg     `json:"videoMsg"`     // 视频消息
	FileMsg      *FileMsg      `json:"fileMsg"`      // 文件消息
	VoiceMsg     *VoiceMsg     `json:"voiceMsg"`     // 语音消息
	VoiceCallMsg *VoiceCallMsg `json:"voiceCallMsg"` // 语言通话
	VideoCallMsg *VideoCallMsg `json:"videoCallMsg"` // 视频通话
	WithdrawMsg  *WithdrawMsg  `json:"withdrawMsg"`  // 撤回消息
	ReplyMsg     *ReplyMsg     `json:"replyMsg"`     // 回复消息
	QuoteMsg     *QuoteMsg     `json:"quoteMsg"`     // 引用消息
	AtMsg        *AtMsg        `json:"atMsg"`        // @消息
}

// Scan 取出来时的数据
func (m *Msg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), m)
}

// Value 入库时的数据
func (m *Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

type ImageMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
}

type VideoMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Time  int    `json:"time"` // 时长（秒）
}

type FileMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Size  int64  `json:"size"` // 文件大小
	Type  string `json:"type"` // 文件类型
}

type VoiceMsg struct {
	Src  string `json:"src"`
	Time int    `json:"time"` // 时长（秒）
}

type VoiceCallMsg struct {
	StartTime time.Time `json:"startTime"` // 开始时间
	EndTime   time.Time `json:"endTime"`   // 结束时间
	EndReason int8      `json:"endReason"` // 结束原因
}

type VideoCallMsg struct {
	StartTime time.Time `json:"startTime"` // 开始时间
	EndTime   time.Time `json:"endTime"`   // 结束时间
	EndReason int8      `json:"endReason"` // 结束原因
}

type WithdrawMsg struct {
	Content   string `json:"content"` // 撤回提示词
	OriginMsg *Msg   `json:"originMsg"`
}

type ReplyMsg struct {
	MsgID   int64  `json:"msgID"`   // 被回复消息ID
	Content string `json:"content"` // 回复文本
	Msg     *Msg   `json:"msg"`
}

type QuoteMsg struct {
	MsgID   int64  `json:"msgID"`   // 引用消息ID
	Content string `json:"content"` // 引用文本
	Msg     *Msg   `json:"msg"`
}

type AtMsg struct {
	UserID  int64  `json:"userID"`  // 被@的用户ID
	Content string `json:"content"` // @消息内容
	Msg     *Msg   `json:"msg"`
}

// 映射实现

// 消息类型

// GetMsgTypeText 响应给前端：在返回数据时，将数值转为文本
func GetMsgTypeText(MsgType int8) string {
	switch MsgType {
	case 1:
		return "文本类型"
	case 2:
		return "图片消息"
	case 3:
		return "视频消息"
	case 4:
		return "文件消息"
	case 5:
		return "语音消息"
	case 6:
		return "语言通话"
	case 7:
		return "视频通话"
	case 8:
		return "撤回消息"
	case 9:
		return "回复消息"
	case 10:
		return "引用消息"
	default:
		return "未知"
	}
}

// GetMsgTypeValue 后端接受处理时解析为数值存储
func GetMsgTypeValue(MsgType string) int8 {
	switch MsgType {
	case "文本类型":
		return 1
	case "图片消息":
		return 2
	case "视频消息":
		return 3
	case "文件消息":
		return 4
	case "语音消息":
		return 5
	case "语言通话":
		return 6
	case "视频通话":
		return 7
	case "撤回消息":
		return 8
	case "回复消息":
		return 9
	case "引用消息":
		return 10
	default:
		return 99 // 未知或未指定
	}
}

package user_domain

import (
	"time"
)

// User 领域对象
type User struct {
	ID         int64     `json:"id"`
	CreateTime time.Time `json:"createTime"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	Nickname   string    `json:"nickname"`
	Signature  string    `json:"signature"`
	Avatar     string    `json:"avatar"`
	Address    string    `json:"address"`
	Birthday   int64     `json:"birthday"`
	Sex        int8      `json:"sex"`
	UserConf   UserConf  `json:"userConf"`
}

type UserConf struct {
	ID                   int64                 `json:"id"`
	CreateTime           time.Time             `json:"createTime"`
	RecallMessage        *string               `json:"recallMessage"`
	FriendOnline         bool                  `json:"friendOnline"`
	Sound                bool                  `json:"sound"`
	SecureLink           bool                  `json:"secureLink"`
	SavePwd              bool                  `json:"savePwd"`
	SearchUser           int8                  `json:"searchUser"`
	Verification         int8                  `json:"verification"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"`
	Online               bool                  `json:"online"`
}

type VerificationQuestion struct {
	Problem1 *string `json:"problem1"`
	Problem2 *string `json:"problem2"`
	Problem3 *string `json:"problem3"`
	Answer1  *string `json:"answer1"`
	Answer2  *string `json:"answer2"`
	Answer3  *string `json:"answer3"`
}

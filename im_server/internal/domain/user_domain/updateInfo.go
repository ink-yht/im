package user_domain

// UpdateInfoRequest 用户邮箱登录请求体
type UpdateInfoRequest struct {
	ID        int64              `json:"id"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Nickname  string             `json:"nickname"`
	Signature string             `json:"signature"`
	Avatar    string             `json:"avatar"`
	Address   string             `json:"address"`
	Birthday  int64              `json:"birthday"`
	Sex       int8               `json:"sex"`
	UserConf  UpdateInfoUserConf `json:"userConf"`
}

type UpdateInfoUserConf struct {
	ID                   int64                           `json:"id"`
	RecallMessage        *string                         `json:"recallMessage"`
	FriendOnline         bool                            `json:"friendOnline"`
	Sound                bool                            `json:"sound"`
	SecureLink           bool                            `json:"secureLink"`
	SavePwd              bool                            `json:"savePwd"`
	SearchUser           int8                            `json:"searchUser"`
	Verification         int8                            `json:"verification"`
	VerificationQuestion *UpdateInfoVerificationQuestion `json:"verificationQuestion"`
	Online               bool                            `json:"online"`
}

type UpdateInfoVerificationQuestion struct {
	Problem1 *string `json:"problem1"`
	Problem2 *string `json:"problem2"`
	Problem3 *string `json:"problem3"`
	Answer1  *string `json:"answer1"`
	Answer2  *string `json:"answer2"`
	Answer3  *string `json:"answer3"`
}

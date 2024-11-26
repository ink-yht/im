package user_domain

// EmailLoginRequest 用户邮箱登录请求体
type EmailLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

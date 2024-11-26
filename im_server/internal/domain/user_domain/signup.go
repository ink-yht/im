package user_domain

import (
	"errors"
	"github.com/dlclark/regexp2"
)

var (
	emailRegex                          = regexp2.MustCompile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, regexp2.None)
	passwordRegex                       = regexp2.MustCompile(`^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[._~!@#$^&*])[A-Za-z0-9._~!@#$^&*]{8,20}$`, regexp2.None)
	ErrTheMailboxIsNotInTheRightFormat  = errors.New("电子邮件格式无效")
	ErrThePasswordIsNotInTheRightFormat = errors.New("密码长度必须为 8-20 个字符，并包含字母、数字和特殊字符")
	ErrThePasswordIsInconsistentTwice   = errors.New("两次密码不一致")
)

// UserRegisterRequest 用户注册请求体
type UserRegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// Validate 校验请求参数
func (req *UserRegisterRequest) Validate() error {
	// 校验邮箱格式
	if match, _ := emailRegex.MatchString(req.Email); !match {
		return ErrTheMailboxIsNotInTheRightFormat
	}

	// 校验密码格式
	if match, _ := passwordRegex.MatchString(req.Password); !match {
		return ErrThePasswordIsNotInTheRightFormat
	}

	// 确认密码是否一致
	if req.Password != req.ConfirmPassword {
		return ErrThePasswordIsInconsistentTwice
	}

	return nil
}

// DefaultUser 初始化默认用户
func DefaultUser(email, password string) *User {
	return &User{
		Email:     email,
		Phone:     "",
		Password:  password,
		Nickname:  "",
		Signature: "",
		Avatar:    "",
		Address:   "",
		Birthday:  0,
		Sex:       0,
		UserConf: UserConf{
			RecallMessage: nil,
			FriendOnline:  false,
			Sound:         false,
			SecureLink:    false,
			SavePwd:       false,
			SearchUser:    1,
			Verification:  1,
			Online:        false,
		},
	}
}

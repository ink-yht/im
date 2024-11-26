package user_service

import (
	"context"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/repository/user_repo"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicate = user_repo.ErrDuplicate

// UserService 定义了用户服务的接口
type UserService interface {
	Signup(ctx context.Context, req user_domain.UserRegisterRequest) error
}

// UserServiceImpl 实现了 UserService 接口
type UserServiceImpl struct {
	repo user_repo.UserRepository
}

func NewUserService(repo user_repo.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (svc *UserServiceImpl) Signup(ctx context.Context, req user_domain.UserRegisterRequest) error {
	// 校验请求
	if err := req.Validate(); err != nil {
		return err
	}

	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// 密码加密失败
		return err
	}

	// 初始化用户
	user := user_domain.DefaultUser(req.Email, string(hash))

	// 插入用户数据
	err = svc.repo.Create(ctx, user_domain.User{
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Address:  user.Address,
		Birthday: user.Birthday,
		Sex:      user.Sex,
		UserConf: user_domain.UserConf{
			RecallMessage: user.UserConf.RecallMessage,
			FriendOnline:  user.UserConf.FriendOnline,
			Sound:         user.UserConf.Sound,
			SecureLink:    user.UserConf.SecureLink,
			SavePwd:       user.UserConf.SavePwd,
			SearchUser:    user.UserConf.SearchUser,
			Verification:  user.UserConf.Verification,
			Online:        user.UserConf.Online,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

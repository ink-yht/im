package user_service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/repository/user_repo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicate      = user_repo.ErrDuplicate
	ErrRecordNotFound = user_repo.ErrRecordNotFound
)

// UserService 定义了用户服务的接口
type UserService interface {
	Signup(ctx context.Context, req user_domain.EmailRegisterRequest) error
	Login(ctx context.Context, req *user_domain.EmailLoginRequest, userAgent string) (string, error)
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

func (svc *UserServiceImpl) Login(ctx context.Context, req *user_domain.EmailLoginRequest, userAgent string) (string, error) {
	// 从数据库中查找用户
	user, err := svc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}

	// 生成 JWT
	token, err := svc.setJWTToken(ctx, user.ID, userAgent)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (svc *UserServiceImpl) Signup(ctx context.Context, req user_domain.EmailRegisterRequest) error {
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

// setJWTToken 生成 token
func (svc *UserServiceImpl) setJWTToken(ctx context.Context, uid int64, userAgent string) (string, error) {
	tokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id:        uid,
		UserAgent: userAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间设置
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
		},
	})
	token, err := tokenStr.SignedString(JWTKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

package user_repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"time"
)

var (
	ErrDuplicate      = user_dao.ErrDuplicate
	ErrRecordNotFound = user_dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user user_domain.User) error
	FindByEmail(ctx context.Context, email string) (user user_domain.User, err error)
	FindByID(ctx context.Context, id int64) (user_domain.User, error)
	UpdateInfo(ctx context.Context, user user_domain.User) error
}

type UserRepositoryImpl struct {
	dao user_dao.UserDao
}

func NewUserRepository(dao user_dao.UserDao) UserRepository {
	return &UserRepositoryImpl{
		dao: dao,
	}
}

func (repo *UserRepositoryImpl) UpdateInfo(ctx context.Context, user user_domain.User) error {
	return repo.dao.UpdateInfo(ctx, repo.domainToEntity(user))
}

func (repo *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (user_domain.User, error) {
	daoUser, err := repo.dao.FindByID(ctx, id)
	if err != nil {
		return user_domain.User{}, err
	}
	return repo.entityToDomain(daoUser), nil
}

func (repo *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (user user_domain.User, err error) {
	daoUser, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return user_domain.User{}, err
	}
	return repo.entityToDomain(daoUser), nil
}

func (repo *UserRepositoryImpl) Create(ctx context.Context, user user_domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *UserRepositoryImpl) domainToEntity(u user_domain.User) user_dao.User {
	var verificationQuestionJSON []byte
	if u.UserConf.VerificationQuestion != nil {
		verificationQuestionJSON, _ = json.Marshal(u.UserConf.VerificationQuestion)
	}
	return user_dao.User{
		ID:         u.ID,
		CreateTime: u.CreateTime.UnixMilli(),
		UpdateTime: u.UpdateTime.UnixMilli(),
		Email: sql.NullString{
			String: u.Email,
			Valid:  len(u.Email) > 0,
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  len(u.Phone) > 0,
		},
		Password:  u.Password,
		Nickname:  u.Nickname,
		Signature: u.Signature,
		Avatar:    u.Avatar,
		Address:   u.Address,
		Birthday:  u.Birthday,
		Sex:       u.Sex,
		UserConf: user_dao.UserConf{
			ID:                   u.UserConf.ID,
			CreateTime:           u.UserConf.CreateTime.UnixMilli(),
			UpdateTime:           u.UserConf.UpdateTime.UnixMilli(),
			RecallMessage:        u.UserConf.RecallMessage,
			FriendOnline:         u.UserConf.FriendOnline,
			Sound:                u.UserConf.Sound,
			SecureLink:           u.UserConf.SecureLink,
			SavePwd:              u.UserConf.SavePwd,
			SearchUser:           u.UserConf.SearchUser,
			Verification:         u.UserConf.Verification,
			VerificationQuestion: string(verificationQuestionJSON),
			Online:               u.UserConf.Online,
		},
	}
}

func (repo *UserRepositoryImpl) entityToDomain(u user_dao.User) user_domain.User {

	var verificationQuestion user_domain.VerificationQuestion
	if u.UserConf.VerificationQuestion != "" {
		_ = json.Unmarshal([]byte(u.UserConf.VerificationQuestion), &verificationQuestion)
	}
	return user_domain.User{
		ID:         u.ID,
		CreateTime: time.UnixMilli(u.CreateTime),
		Email:      nullStringToString(u.Email),
		Phone:      nullStringToString(u.Phone),
		Password:   u.Password,
		Nickname:   u.Nickname,
		Signature:  u.Signature,
		Avatar:     u.Avatar,
		Address:    u.Address,
		Birthday:   u.Birthday,
		Sex:        u.Sex,
		UserConf: user_domain.UserConf{
			ID:                   u.UserConf.ID,
			CreateTime:           time.UnixMilli(u.UserConf.CreateTime),
			RecallMessage:        u.UserConf.RecallMessage,
			FriendOnline:         u.UserConf.FriendOnline,
			Sound:                u.UserConf.Sound,
			SecureLink:           u.UserConf.SecureLink,
			SavePwd:              u.UserConf.SavePwd,
			SearchUser:           u.UserConf.SearchUser,
			Verification:         u.UserConf.Verification,
			VerificationQuestion: &verificationQuestion,
			Online:               u.UserConf.Online,
		},
	}
}

// 辅助函数：将 sql.NullString 转换为 string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

package user_repo

import (
	"context"
	"database/sql"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"time"
)

var ErrDuplicate = user_dao.ErrDuplicate

type UserRepository interface {
	Create(ctx context.Context, user user_domain.User) error
}

type UserRepositoryImpl struct {
	dao user_dao.UserDao
}

func NewUserRepository(dao user_dao.UserDao) UserRepository {
	return &UserRepositoryImpl{
		dao: dao,
	}
}

func (repo *UserRepositoryImpl) Create(ctx context.Context, user user_domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *UserRepositoryImpl) domainToEntity(u user_domain.User) user_dao.User {
	return user_dao.User{
		ID:         u.ID,
		CreateTime: u.CreateTime.UnixMilli(),
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
			RecallMessage:        u.UserConf.RecallMessage,
			FriendOnline:         u.UserConf.FriendOnline,
			Sound:                u.UserConf.Sound,
			SecureLink:           u.UserConf.SecureLink,
			SavePwd:              u.UserConf.SavePwd,
			SearchUser:           u.UserConf.SearchUser,
			Verification:         u.UserConf.Verification,
			VerificationQuestion: domainVerificationToDAO(u.UserConf.VerificationQuestion),
			Online:               u.UserConf.Online,
		},
	}
}

// 辅助函数：领域模型 VerificationQuestion 转换为 DAO 模型
func domainVerificationToDAO(vq *user_domain.VerificationQuestion) *user_dao.VerificationQuestion {
	if vq == nil {
		return nil
	}
	return &user_dao.VerificationQuestion{
		Problem1: vq.Problem1,
		Problem2: vq.Problem2,
		Problem3: vq.Problem3,
		Answer1:  vq.Answer1,
		Answer2:  vq.Answer2,
		Answer3:  vq.Answer3,
	}
}

func (repo *UserRepositoryImpl) entityToDomain(u user_dao.User) user_domain.User {
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
			VerificationQuestion: daoVerificationToDomain(u.UserConf.VerificationQuestion),
			Online:               u.UserConf.Online,
		},
	}
}

// 辅助函数：DAO 模型 VerificationQuestion 转换为领域模型
func daoVerificationToDomain(vq *user_dao.VerificationQuestion) *user_domain.VerificationQuestion {
	if vq == nil {
		return nil
	}
	return &user_domain.VerificationQuestion{
		Problem1: vq.Problem1,
		Problem2: vq.Problem2,
		Problem3: vq.Problem3,
		Answer1:  vq.Answer1,
		Answer2:  vq.Answer2,
		Answer3:  vq.Answer3,
	}
}

// 辅助函数：将 sql.NullString 转换为 string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

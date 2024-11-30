package file_repo

import (
	"database/sql"
	"github.com/goccy/go-json"
	"github.com/ink-yht/im/internal/domain/user_domain"
	"github.com/ink-yht/im/internal/repository/dao/file_dao"
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"golang.org/x/net/context"
	"time"
)

type FileRepository interface {
	Avatar(ctx context.Context, user user_domain.User) error
}

type FileRepositoryImpl struct {
	dao file_dao.FileDao
}

func NewFileRepository(dao file_dao.FileDao) FileRepository {
	return &FileRepositoryImpl{
		dao: dao,
	}
}

func (repo *FileRepositoryImpl) Avatar(ctx context.Context, user user_domain.User) error {
	return repo.dao.Avatar(ctx, repo.domainToEntity(user))
}

func (repo *FileRepositoryImpl) domainToEntity(u user_domain.User) user_dao.User {
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

func (repo *FileRepositoryImpl) entityToDomain(u user_dao.User) user_domain.User {

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

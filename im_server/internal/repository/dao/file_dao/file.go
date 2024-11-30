package file_dao

import (
	"github.com/ink-yht/im/internal/repository/dao/user_dao"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type FileDao interface {
	Avatar(ctx context.Context, u user_dao.User) error
}

type GormFileDAO struct {
	db *gorm.DB
}

func NewFileDAO(db *gorm.DB) FileDao {
	return &GormFileDAO{db: db}
}

func (dao GormFileDAO) Avatar(ctx context.Context, u user_dao.User) error {
	return dao.db.WithContext(ctx).Model(&u).Where("id = ?", u.ID).Update("Avatar", u.Avatar).Error
}

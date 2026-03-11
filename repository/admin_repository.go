package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *model.Admin) error
	FindByID(ctx context.Context, id uint) (*model.Admin, error)
	Update(ctx context.Context, admin *model.Admin) error
	Delete(ctx context.Context, id uint) error

	Paginate(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Admin, int64, error)

	FindByUsername(ctx context.Context, username string) (*model.Admin, error)
	FindByEmail(ctx context.Context, email string) (*model.Admin, error)
	FindByPhone(ctx context.Context, phone string) (*model.Admin, error)

	BatchDelete(ctx context.Context, ids []uint) error
	BatchUpdateStatus(ctx context.Context, ids []uint, status int8) error

	UpdateLastLoginAt(ctx context.Context, id uint) error
}

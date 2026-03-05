package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
)

type LoginLogFilters struct {
	UserID *uint
}

type AdminLoginLogRepository interface {
	Create(ctx context.Context, log *model.AdminLoginLog) error
	Paginate(ctx context.Context, page, pageSize int, filters LoginLogFilters) ([]model.AdminLoginLog, int64, error)
	FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AdminLoginLog, error)
}

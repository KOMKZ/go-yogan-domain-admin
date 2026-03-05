package service

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
)

type AdminLoginLogService struct {
	repo repository.AdminLoginLogRepository
}

func NewAdminLoginLogService(repo repository.AdminLoginLogRepository) *AdminLoginLogService {
	return &AdminLoginLogService{repo: repo}
}

func (s *AdminLoginLogService) Create(ctx context.Context, log *model.AdminLoginLog) error {
	return s.repo.Create(ctx, log)
}

func (s *AdminLoginLogService) Paginate(ctx context.Context, page, pageSize int, filters repository.LoginLogFilters) ([]model.AdminLoginLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.Paginate(ctx, page, pageSize, filters)
}

func (s *AdminLoginLogService) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AdminLoginLog, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.FindByUserID(ctx, userID, limit)
}

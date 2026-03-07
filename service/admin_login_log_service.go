package service

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
)

type PaginateLoginLogInput struct {
	UserID *uint
}

type CreateLoginLogInput struct {
	UserID    uint
	Username  string
	IP        string
	UserAgent string
}

type AdminLoginLogService struct {
	repo repository.AdminLoginLogRepository
}

func NewAdminLoginLogService(repo repository.AdminLoginLogRepository) *AdminLoginLogService {
	return &AdminLoginLogService{repo: repo}
}

func (s *AdminLoginLogService) Create(ctx context.Context, input CreateLoginLogInput) error {
	log := &model.AdminLoginLog{
		UserID:    input.UserID,
		Username:  input.Username,
		IP:        input.IP,
		UserAgent: input.UserAgent,
	}
	return s.repo.Create(ctx, log)
}

func (s *AdminLoginLogService) Paginate(ctx context.Context, page, pageSize int, input PaginateLoginLogInput) ([]model.AdminLoginLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	filters := repository.LoginLogFilters{
		UserID: input.UserID,
	}
	return s.repo.Paginate(ctx, page, pageSize, filters)
}

func (s *AdminLoginLogService) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AdminLoginLog, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.FindByUserID(ctx, userID, limit)
}

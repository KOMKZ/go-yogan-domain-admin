package repository

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"gorm.io/gorm"
)

type AdminLoginLogMySQLRepository struct {
	db *gorm.DB
}

func NewAdminLoginLogMySQLRepository(db *gorm.DB) *AdminLoginLogMySQLRepository {
	return &AdminLoginLogMySQLRepository{db: db}
}

func (r *AdminLoginLogMySQLRepository) Create(ctx context.Context, log *model.AdminLoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AdminLoginLogMySQLRepository) Paginate(ctx context.Context, page, pageSize int, filters LoginLogFilters) ([]model.AdminLoginLog, int64, error) {
	var logs []model.AdminLoginLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AdminLoginLog{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *AdminLoginLogMySQLRepository) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AdminLoginLog, error) {
	var logs []model.AdminLoginLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("id DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

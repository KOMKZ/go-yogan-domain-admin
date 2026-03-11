package repository

import (
	"context"
	"errors"
	"time"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-framework/database"
	"gorm.io/gorm"
)

type AdminMySQLRepository struct {
	base *database.BaseRepository[model.Admin]
	db   *gorm.DB
}

func NewAdminMySQLRepository(db *gorm.DB) *AdminMySQLRepository {
	return &AdminMySQLRepository{
		base: database.NewBaseRepository[model.Admin](db),
		db:   db,
	}
}

func (r *AdminMySQLRepository) Create(ctx context.Context, admin *model.Admin) error {
	return r.base.Create(ctx, admin)
}

func (r *AdminMySQLRepository) FindByID(ctx context.Context, id uint) (*model.Admin, error) {
	result, err := r.base.FindByID(ctx, id)
	if errors.Is(err, database.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}

func (r *AdminMySQLRepository) Update(ctx context.Context, admin *model.Admin) error {
	return r.base.Update(ctx, admin)
}

func (r *AdminMySQLRepository) Delete(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *AdminMySQLRepository) Paginate(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Admin, int64, error) {
	var admins []model.Admin
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Admin{})

	if username, ok := filters["username"]; ok {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&admins).Error; err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

func (r *AdminMySQLRepository) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminMySQLRepository) FindByEmail(ctx context.Context, email string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminMySQLRepository) FindByPhone(ctx context.Context, phone string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminMySQLRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Admin{}).Error
}

func (r *AdminMySQLRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status int8) error {
	return r.db.WithContext(ctx).Model(&model.Admin{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *AdminMySQLRepository) UpdateLastLoginAt(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.Admin{}).Where("id = ?", id).Update("last_login_at", now).Error
}

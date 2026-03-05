package service

import (
	"context"

	domainerrors "github.com/KOMKZ/go-yogan-domain-admin/errors"
	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

type CreateAdminInput struct {
	Username string
	Password string
	RealName string
	Email    string
	Phone    string
	Role     int8
	Status   int8
}

type UpdateAdminInput struct {
	Username *string
	RealName *string
	Email    *string
	Phone    *string
	Status   *int8
}

type PageResult struct {
	Records []model.Admin `json:"records"`
	Total   int64         `json:"total"`
	Size    int           `json:"size"`
	Current int           `json:"current"`
	Pages   int           `json:"pages"`
}

func (s *AdminService) Create(ctx context.Context, input CreateAdminInput) (*model.Admin, error) {
	existing, err := s.repo.FindByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domainerrors.ErrUsernameExists
	}

	if input.Email != "" {
		existingEmail, err := s.repo.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, err
		}
		if existingEmail != nil {
			return nil, domainerrors.ErrEmailExists
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := input.Role
	if role == 0 {
		role = 2
	}
	status := input.Status
	if status == 0 {
		status = 1
	}

	admin := &model.Admin{
		Username: input.Username,
		Password: string(passwordHash),
		RealName: input.RealName,
		Email:    input.Email,
		Phone:    input.Phone,
		Role:     role,
		Status:   status,
	}

	if err := s.repo.Create(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *AdminService) GetByID(ctx context.Context, id uint) (*model.Admin, error) {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, domainerrors.ErrAdminNotFound
	}
	return admin, nil
}

func (s *AdminService) Update(ctx context.Context, id uint, input UpdateAdminInput) (*model.Admin, error) {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, domainerrors.ErrAdminNotFound
	}

	if input.Username != nil && *input.Username != admin.Username {
		existing, err := s.repo.FindByUsername(ctx, *input.Username)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, domainerrors.ErrUsernameExists
		}
		admin.Username = *input.Username
	}
	if input.RealName != nil {
		admin.RealName = *input.RealName
	}
	if input.Email != nil && *input.Email != admin.Email {
		existing, err := s.repo.FindByEmail(ctx, *input.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, domainerrors.ErrEmailExists
		}
		admin.Email = *input.Email
	}
	if input.Phone != nil {
		admin.Phone = *input.Phone
	}
	if input.Status != nil {
		admin.Status = *input.Status
	}

	if err := s.repo.Update(ctx, admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *AdminService) Delete(ctx context.Context, id uint) error {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if admin == nil {
		return domainerrors.ErrAdminNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *AdminService) ListPage(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*PageResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	admins, total, err := s.repo.Paginate(ctx, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}

	return &PageResult{
		Records: admins,
		Total:   total,
		Size:    pageSize,
		Current: page,
		Pages:   pages,
	}, nil
}

func (s *AdminService) BatchDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BatchDelete(ctx, ids)
}

func (s *AdminService) BatchUpdateStatus(ctx context.Context, ids []uint, status int8) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BatchUpdateStatus(ctx, ids, status)
}

func (s *AdminService) ResetPassword(ctx context.Context, id uint, newPassword string) error {
	admin, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if admin == nil {
		return domainerrors.ErrAdminNotFound
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin.Password = string(passwordHash)
	return s.repo.Update(ctx, admin)
}

func (s *AdminService) UpdateLastLoginAt(ctx context.Context, id uint) error {
	return s.repo.UpdateLastLoginAt(ctx, id)
}

// GetPermissionsByRole 根据角色获取权限列表（admin 域特有的业务逻辑）
func (s *AdminService) GetPermissionsByRole(role int8) []string {
	if role == 1 {
		return []string{
			"*:*",
			"dashboard:view",
			"user:view", "user:write", "user:delete",
			"admin:view", "admin:write", "admin:delete",
			"role:view", "role:write", "role:delete",
			"permission:view", "permission:write",
			"system:view", "system:write",
			"log:view",
		}
	}
	return []string{
		"dashboard:view",
		"user:view",
	}
}

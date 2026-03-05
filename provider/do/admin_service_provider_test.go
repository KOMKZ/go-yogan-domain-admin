package providerdo

import (
	"context"
	"testing"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
	"github.com/samber/do/v2"
)

type stubAdminRepo struct{}

func (stubAdminRepo) Create(ctx context.Context, admin *model.Admin) error { return nil }
func (stubAdminRepo) FindByID(ctx context.Context, id uint) (*model.Admin, error) { return nil, nil }
func (stubAdminRepo) Update(ctx context.Context, admin *model.Admin) error   { return nil }
func (stubAdminRepo) Delete(ctx context.Context, id uint) error             { return nil }
func (stubAdminRepo) Paginate(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Admin, int64, error) {
	return nil, 0, nil
}
func (stubAdminRepo) FindByUsername(ctx context.Context, username string) (*model.Admin, error) { return nil, nil }
func (stubAdminRepo) FindByEmail(ctx context.Context, email string) (*model.Admin, error)       { return nil, nil }
func (stubAdminRepo) BatchDelete(ctx context.Context, ids []uint) error                         { return nil }
func (stubAdminRepo) BatchUpdateStatus(ctx context.Context, ids []uint, status int8) error      { return nil }
func (stubAdminRepo) UpdateLastLoginAt(ctx context.Context, id uint) error                       { return nil }

type stubLoginLogRepo struct{}

func (stubLoginLogRepo) Create(ctx context.Context, log *model.AdminLoginLog) error { return nil }
func (stubLoginLogRepo) Paginate(ctx context.Context, page, pageSize int, filters repository.LoginLogFilters) ([]model.AdminLoginLog, int64, error) {
	return nil, 0, nil
}
func (stubLoginLogRepo) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AdminLoginLog, error) {
	return nil, nil
}

func TestProvideAdminService(t *testing.T) {
	injector := do.New()
	do.Provide(injector, func(i do.Injector) (repository.AdminRepository, error) {
		return stubAdminRepo{}, nil
	})

	svc, err := ProvideAdminService(injector)
	if err != nil {
		t.Fatalf("ProvideAdminService() err = %v", err)
	}
	if svc == nil {
		t.Fatal("ProvideAdminService() returned nil service")
	}
}

func TestProvideAdminLoginLogService(t *testing.T) {
	injector := do.New()
	do.Provide(injector, func(i do.Injector) (repository.AdminLoginLogRepository, error) {
		return stubLoginLogRepo{}, nil
	})

	svc, err := ProvideAdminLoginLogService(injector)
	if err != nil {
		t.Fatalf("ProvideAdminLoginLogService() err = %v", err)
	}
	if svc == nil {
		t.Fatal("ProvideAdminLoginLogService() returned nil service")
	}
}

func TestProvideAdminService_InvokeFails(t *testing.T) {
	// Empty injector - no AdminRepository provided
	injector := do.New()

	_, err := ProvideAdminService(injector)
	if err == nil {
		t.Fatal("ProvideAdminService() expected error when repo not in injector")
	}
}

func TestProvideAdminLoginLogService_InvokeFails(t *testing.T) {
	// Empty injector - no AdminLoginLogRepository provided
	injector := do.New()

	_, err := ProvideAdminLoginLogService(injector)
	if err == nil {
		t.Fatal("ProvideAdminLoginLogService() expected error when repo not in injector")
	}
}

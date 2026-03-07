package do

import (
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
	"github.com/KOMKZ/go-yogan-domain-admin/service"
	frameworkauth "github.com/KOMKZ/go-yogan-framework/auth"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

// ---- Repository Providers ----

func ProvideAdminRepository(i do.Injector) (repository.AdminRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewAdminMySQLRepository(db), nil
}

func ProvideAdminLoginLogRepository(i do.Injector) (repository.AdminLoginLogRepository, error) {
	db, err := do.Invoke[*gorm.DB](i)
	if err != nil {
		return nil, err
	}
	return repository.NewAdminLoginLogMySQLRepository(db), nil
}

// ---- Service Providers ----

func ProvideAdminService(i do.Injector) (*service.AdminService, error) {
	repo, err := do.Invoke[repository.AdminRepository](i)
	if err != nil {
		return nil, err
	}
	passwordSvc, err := do.Invoke[*frameworkauth.PasswordService](i)
	if err != nil {
		return nil, err
	}
	log := do.MustInvokeNamed[*logger.CtxZapLogger](i, "admin")
	return service.NewAdminService(repo, passwordSvc, log), nil
}

func ProvideAdminLoginLogService(i do.Injector) (*service.AdminLoginLogService, error) {
	repo, err := do.Invoke[repository.AdminLoginLogRepository](i)
	if err != nil {
		return nil, err
	}
	return service.NewAdminLoginLogService(repo), nil
}

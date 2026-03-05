package providerdo

import (
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
	"github.com/KOMKZ/go-yogan-domain-admin/service"
	"github.com/samber/do/v2"
)

func ProvideAdminService(i do.Injector) (*service.AdminService, error) {
	repo, err := do.Invoke[repository.AdminRepository](i)
	if err != nil {
		return nil, err
	}
	return service.NewAdminService(repo), nil
}

func ProvideAdminLoginLogService(i do.Injector) (*service.AdminLoginLogService, error) {
	repo, err := do.Invoke[repository.AdminLoginLogRepository](i)
	if err != nil {
		return nil, err
	}
	return service.NewAdminLoginLogService(repo), nil
}

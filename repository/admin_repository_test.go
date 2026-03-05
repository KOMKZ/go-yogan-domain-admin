package repository

import "testing"

func TestAdminRepositoryInterfaceCompiles(t *testing.T) {
	var _ AdminRepository = nil
	_ = t
}

func TestAdminLoginLogRepositoryInterfaceCompiles(t *testing.T) {
	var _ AdminLoginLogRepository = nil
	_ = t
}

func TestLoginLogFiltersZeroValue(t *testing.T) {
	f := LoginLogFilters{}
	if f.UserID != nil {
		t.Fatal("zero value UserID should be nil")
	}
}

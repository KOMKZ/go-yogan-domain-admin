package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
)

type mockLoginLogRepo struct {
	mu          sync.RWMutex
	logs        []*model.AdminLoginLog
	nextID      uint
	createErr   error
	paginateErr error
	findByErr   error
}

func newMockLoginLogRepo() *mockLoginLogRepo {
	return &mockLoginLogRepo{
		logs:   make([]*model.AdminLoginLog, 0),
		nextID: 1,
	}
}

func (m *mockLoginLogRepo) Create(ctx context.Context, log *model.AdminLoginLog) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if log.ID == 0 {
		log.ID = m.nextID
		m.nextID++
	}
	cp := *log
	m.logs = append(m.logs, &cp)
	return nil
}

func (m *mockLoginLogRepo) Paginate(ctx context.Context, page, pageSize int, filters repository.LoginLogFilters) ([]model.AdminLoginLog, int64, error) {
	if m.paginateErr != nil {
		return nil, 0, m.paginateErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []model.AdminLoginLog
	for _, l := range m.logs {
		if filters.UserID != nil && l.UserID != *filters.UserID {
			continue
		}
		filtered = append(filtered, *l)
	}
	total := int64(len(filtered))
	offset := (page - 1) * pageSize
	if offset >= len(filtered) {
		return []model.AdminLoginLog{}, total, nil
	}
	end := offset + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (m *mockLoginLogRepo) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AdminLoginLog, error) {
	if m.findByErr != nil {
		return nil, m.findByErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []model.AdminLoginLog
	for i := len(m.logs) - 1; i >= 0 && len(result) < limit; i-- {
		if m.logs[i].UserID == userID {
			result = append(result, *m.logs[i])
		}
	}
	return result, nil
}

var _ repository.AdminLoginLogRepository = (*mockLoginLogRepo)(nil)

func TestAdminLoginLogService_Create_Success(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	log := &model.AdminLoginLog{
		UserID:   1,
		Username: "admin",
		IP:       "127.0.0.1",
		UserAgent: "test",
	}
	if err := svc.Create(ctx, log); err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if log.ID == 0 {
		t.Error("Create should set ID")
	}
}

func TestAdminLoginLogService_Create_RepoErr(t *testing.T) {
	repo := newMockLoginLogRepo()
	repo.createErr = errors.New("db error")
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	err := svc.Create(ctx, &model.AdminLoginLog{UserID: 1, Username: "a", IP: "x"})
	if err == nil || err.Error() != "db error" {
		t.Errorf("err = %v, want db error", err)
	}
}

func TestAdminLoginLogService_Paginate_Success(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_ = svc.Create(ctx, &model.AdminLoginLog{UserID: 1, Username: "a", IP: "1"})
	_ = svc.Create(ctx, &model.AdminLoginLog{UserID: 1, Username: "a", IP: "2"})

	logs, total, err := svc.Paginate(ctx, 1, 10, repository.LoginLogFilters{})
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if len(logs) != 2 {
		t.Errorf("len(logs) = %d, want 2", len(logs))
	}
}

func TestAdminLoginLogService_Paginate_NormalizePage(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_, _, err := svc.Paginate(ctx, 0, 10, repository.LoginLogFilters{})
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	// Repo receives page=1 after normalization
}

func TestAdminLoginLogService_Paginate_NormalizePageSize(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_, _, err := svc.Paginate(ctx, 1, 0, repository.LoginLogFilters{})
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
}

func TestAdminLoginLogService_Paginate_PageSizeOver100(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_, _, err := svc.Paginate(ctx, 1, 200, repository.LoginLogFilters{})
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
}

func TestAdminLoginLogService_Paginate_WithUserIDFilter(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_ = svc.Create(ctx, &model.AdminLoginLog{UserID: 1, Username: "a", IP: "1"})
	_ = svc.Create(ctx, &model.AdminLoginLog{UserID: 2, Username: "b", IP: "2"})

	uid := uint(1)
	logs, total, err := svc.Paginate(ctx, 1, 10, repository.LoginLogFilters{UserID: &uid})
	if err != nil {
		t.Fatalf("Paginate() err = %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(logs) != 1 || logs[0].UserID != 1 {
		t.Errorf("logs = %v", logs)
	}
}

func TestAdminLoginLogService_Paginate_RepoErr(t *testing.T) {
	repo := newMockLoginLogRepo()
	repo.paginateErr = errors.New("paginate failed")
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_, _, err := svc.Paginate(ctx, 1, 10, repository.LoginLogFilters{})
	if err == nil || err.Error() != "paginate failed" {
		t.Errorf("err = %v, want paginate failed", err)
	}
}

func TestAdminLoginLogService_FindByUserID_Success(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_ = svc.Create(ctx, &model.AdminLoginLog{UserID: 1, Username: "a", IP: "1"})
	_ = svc.Create(ctx, &model.AdminLoginLog{UserID: 1, Username: "a", IP: "2"})

	logs, err := svc.FindByUserID(ctx, 1, 5)
	if err != nil {
		t.Fatalf("FindByUserID() err = %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("len(logs) = %d, want 2", len(logs))
	}
}

func TestAdminLoginLogService_FindByUserID_NormalizeLimit(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	// limit 0 normalized to 20 - call should succeed without error
	logs, err := svc.FindByUserID(ctx, 1, 0)
	if err != nil {
		t.Fatalf("FindByUserID() err = %v", err)
	}
	// No logs for user 1, expect empty result
	if len(logs) != 0 {
		t.Errorf("len(logs) = %d, want 0", len(logs))
	}
}

func TestAdminLoginLogService_FindByUserID_NegativeLimit(t *testing.T) {
	repo := newMockLoginLogRepo()
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	// limit -1 normalized to 20 - call should succeed
	logs, err := svc.FindByUserID(ctx, 1, -1)
	if err != nil {
		t.Fatalf("FindByUserID() err = %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("len(logs) = %d, want 0 for non-existent user", len(logs))
	}
}

func TestAdminLoginLogService_FindByUserID_RepoErr(t *testing.T) {
	repo := newMockLoginLogRepo()
	repo.findByErr = errors.New("find failed")
	svc := NewAdminLoginLogService(repo)
	ctx := context.Background()

	_, err := svc.FindByUserID(ctx, 1, 10)
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

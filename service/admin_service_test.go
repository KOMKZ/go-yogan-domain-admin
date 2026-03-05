package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	domainerrors "github.com/KOMKZ/go-yogan-domain-admin/errors"
	"github.com/KOMKZ/go-yogan-domain-admin/model"
	"github.com/KOMKZ/go-yogan-domain-admin/repository"
)

type mockAdminRepo struct {
	mu               sync.RWMutex
	admins           map[uint]*model.Admin
	byUsername       map[string]*model.Admin
	byEmail          map[string]*model.Admin
	nextID           uint
	createErr        error
	updateErr        error
	deleteErr        error
	findByIDErr      error
	findByUsernameErr error
	findByEmailErr   error
	paginateErr      error
	batchDeleteErr   error
	batchStatusErr   error
	updateLoginErr   error
}

func newMockAdminRepo() *mockAdminRepo {
	return &mockAdminRepo{
		admins:     make(map[uint]*model.Admin),
		byUsername: make(map[string]*model.Admin),
		byEmail:    make(map[string]*model.Admin),
		nextID:     1,
	}
}

func (m *mockAdminRepo) Create(ctx context.Context, admin *model.Admin) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if admin.ID == 0 {
		admin.ID = m.nextID
		m.nextID++
	}
	m.admins[admin.ID] = copyAdmin(admin)
	m.byUsername[admin.Username] = m.admins[admin.ID]
	if admin.Email != "" {
		m.byEmail[admin.Email] = m.admins[admin.ID]
	}
	return nil
}

func (m *mockAdminRepo) FindByID(ctx context.Context, id uint) (*model.Admin, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.admins[id]
	if !ok {
		return nil, nil
	}
	return copyAdmin(a), nil
}

func (m *mockAdminRepo) Update(ctx context.Context, admin *model.Admin) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	old := m.admins[admin.ID]
	if old != nil {
		delete(m.byUsername, old.Username)
		if old.Email != "" {
			delete(m.byEmail, old.Email)
		}
	}
	m.admins[admin.ID] = copyAdmin(admin)
	m.byUsername[admin.Username] = m.admins[admin.ID]
	if admin.Email != "" {
		m.byEmail[admin.Email] = m.admins[admin.ID]
	}
	return nil
}

func (m *mockAdminRepo) Delete(ctx context.Context, id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	a := m.admins[id]
	if a != nil {
		delete(m.byUsername, a.Username)
		if a.Email != "" {
			delete(m.byEmail, a.Email)
		}
		delete(m.admins, id)
	}
	return nil
}

func (m *mockAdminRepo) Paginate(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]model.Admin, int64, error) {
	if m.paginateErr != nil {
		return nil, 0, m.paginateErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []model.Admin
	for _, a := range m.admins {
		list = append(list, *copyAdmin(a))
	}
	total := int64(len(list))
	offset := (page - 1) * pageSize
	if offset >= len(list) {
		return []model.Admin{}, total, nil
	}
	end := offset + pageSize
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], total, nil
}

func (m *mockAdminRepo) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	if m.findByUsernameErr != nil {
		return nil, m.findByUsernameErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.byUsername[username]
	if !ok {
		return nil, nil
	}
	return copyAdmin(a), nil
}

func (m *mockAdminRepo) FindByEmail(ctx context.Context, email string) (*model.Admin, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.byEmail[email]
	if !ok {
		return nil, nil
	}
	return copyAdmin(a), nil
}

func (m *mockAdminRepo) BatchDelete(ctx context.Context, ids []uint) error {
	if m.batchDeleteErr != nil {
		return m.batchDeleteErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range ids {
		a := m.admins[id]
		if a != nil {
			delete(m.byUsername, a.Username)
			if a.Email != "" {
				delete(m.byEmail, a.Email)
			}
			delete(m.admins, id)
		}
	}
	return nil
}

func (m *mockAdminRepo) BatchUpdateStatus(ctx context.Context, ids []uint, status int8) error {
	if m.batchStatusErr != nil {
		return m.batchStatusErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range ids {
		if a := m.admins[id]; a != nil {
			a.Status = status
		}
	}
	return nil
}

func (m *mockAdminRepo) UpdateLastLoginAt(ctx context.Context, id uint) error {
	if m.updateLoginErr != nil {
		return m.updateLoginErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	// In mock we don't persist time; just no-op
	return nil
}

func copyAdmin(a *model.Admin) *model.Admin {
	if a == nil {
		return nil
	}
	cp := *a
	return &cp
}

func int8Ptr(v int8) *int8 { return &v }

func TestAdminService_Create_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	admin, err := svc.Create(ctx, CreateAdminInput{
		Username: "user1",
		Password: "secret123",
		RealName: "Real",
		Email:    "a@b.com",
		Phone:    "123",
		Role:     int8Ptr(2),
		Status:   int8Ptr(1),
	})
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if admin == nil {
		t.Fatal("Create() returned nil admin")
	}
	if admin.Username != "user1" {
		t.Errorf("Username = %q, want user1", admin.Username)
	}
	if admin.Password == "secret123" {
		t.Error("Password should be hashed, not plain")
	}
	if admin.RealName != "Real" {
		t.Errorf("RealName = %q, want Real", admin.RealName)
	}
	if admin.Email != "a@b.com" {
		t.Errorf("Email = %q, want a@b.com", admin.Email)
	}
	if admin.Phone != "123" {
		t.Errorf("Phone = %q, want 123", admin.Phone)
	}
	if admin.Role != 2 {
		t.Errorf("Role = %d, want 2", admin.Role)
	}
	if admin.Status != 1 {
		t.Errorf("Status = %d, want 1", admin.Status)
	}
}

func TestAdminService_Create_DefaultRoleAndStatus(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	admin, err := svc.Create(ctx, CreateAdminInput{
		Username: "user2",
		Password: "pw",
	})
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if admin.Role != 2 {
		t.Errorf("Role = %d, want 2 (default when nil)", admin.Role)
	}
	if admin.Status != 1 {
		t.Errorf("Status = %d, want 1 (default when nil)", admin.Status)
	}
}

func TestAdminService_Create_ExplicitZeroStatus(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	admin, err := svc.Create(ctx, CreateAdminInput{
		Username: "user_disabled",
		Password: "pw",
		Status:   int8Ptr(0),
	})
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if admin.Status != 0 {
		t.Errorf("Status = %d, want 0 (explicit zero should not be overridden)", admin.Status)
	}
}

func TestAdminService_Create_ExplicitRole(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	admin, err := svc.Create(ctx, CreateAdminInput{
		Username: "user_superadmin",
		Password: "pw",
		Role:     int8Ptr(1),
	})
	if err != nil {
		t.Fatalf("Create() err = %v", err)
	}
	if admin.Role != 1 {
		t.Errorf("Role = %d, want 1 (explicit role)", admin.Role)
	}
}

func TestAdminService_Create_UsernameExists(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateAdminInput{Username: "dup", Password: "pw"})
	_, err := svc.Create(ctx, CreateAdminInput{Username: "dup", Password: "pw2"})
	if err != domainerrors.ErrUsernameExists {
		t.Errorf("err = %v, want ErrUsernameExists", err)
	}
}

func TestAdminService_Create_EmailExists(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw", Email: "same@x.com"})
	_, err := svc.Create(ctx, CreateAdminInput{Username: "u2", Password: "pw2", Email: "same@x.com"})
	if err != domainerrors.ErrEmailExists {
		t.Errorf("err = %v, want ErrEmailExists", err)
	}
}

func TestAdminService_Create_EmptyEmailNoConflict(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err1 := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw", Email: ""})
	_, err2 := svc.Create(ctx, CreateAdminInput{Username: "u2", Password: "pw2", Email: ""})
	if err1 != nil || err2 != nil {
		t.Errorf("Create() err1=%v err2=%v", err1, err2)
	}
}

func TestAdminService_Create_RepoCreateErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.createErr = errors.New("db error")
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	if err == nil || err.Error() != "db error" {
		t.Errorf("err = %v, want db error", err)
	}
}

func TestAdminService_Create_FindByUsernameErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.findByUsernameErr = errors.New("db lookup failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	if err == nil || err.Error() != "db lookup failed" {
		t.Errorf("err = %v, want db lookup failed", err)
	}
}

func TestAdminService_Create_FindByEmailErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.findByEmailErr = errors.New("email lookup failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw", Email: "a@b.com"})
	if err == nil || err.Error() != "email lookup failed" {
		t.Errorf("err = %v, want email lookup failed", err)
	}
}

func TestAdminService_GetByID_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() err = %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %d, want %d", got.ID, created.ID)
	}
}

func TestAdminService_GetByID_NotFound(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 999)
	if err != domainerrors.ErrAdminNotFound {
		t.Errorf("err = %v, want ErrAdminNotFound", err)
	}
}

func TestAdminService_GetByID_RepoErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.findByIDErr = errors.New("db error")
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 1)
	if err == nil || err.Error() != "db error" {
		t.Errorf("err = %v, want db error", err)
	}
}

func TestAdminService_Update_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	created, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	rn := "NewReal"
	email := "new@email.com"
	updated, err := svc.Update(ctx, created.ID, UpdateAdminInput{
		RealName: &rn,
		Email:    &email,
	})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.RealName != "NewReal" {
		t.Errorf("RealName = %q, want NewReal", updated.RealName)
	}
	if updated.Email != "new@email.com" {
		t.Errorf("Email = %q, want new@email.com", updated.Email)
	}
}

func TestAdminService_Update_NotFound(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	rn := "x"
	_, err := svc.Update(ctx, 999, UpdateAdminInput{RealName: &rn})
	if err != domainerrors.ErrAdminNotFound {
		t.Errorf("err = %v, want ErrAdminNotFound", err)
	}
}

func TestAdminService_Update_FindByIDErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.findByIDErr = errors.New("find failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	rn := "x"
	_, err := svc.Update(ctx, 1, UpdateAdminInput{RealName: &rn})
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestAdminService_Update_UsernameConflict(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	u2, _ := svc.Create(ctx, CreateAdminInput{Username: "u2", Password: "pw"})
	newUsername := "u1"
	_, err := svc.Update(ctx, u2.ID, UpdateAdminInput{Username: &newUsername})
	if err != domainerrors.ErrUsernameExists {
		t.Errorf("err = %v, want ErrUsernameExists", err)
	}
}

func TestAdminService_Update_UsernameSameNoConflict(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	same := "u1"
	updated, err := svc.Update(ctx, u.ID, UpdateAdminInput{Username: &same})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.Username != "u1" {
		t.Errorf("Username = %q, want u1", updated.Username)
	}
}

func TestAdminService_Update_EmailConflict(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw", Email: "a@b.com"})
	u2, _ := svc.Create(ctx, CreateAdminInput{Username: "u2", Password: "pw", Email: "b@b.com"})
	newEmail := "a@b.com"
	_, err := svc.Update(ctx, u2.ID, UpdateAdminInput{Email: &newEmail})
	if err != domainerrors.ErrEmailExists {
		t.Errorf("err = %v, want ErrEmailExists", err)
	}
}

func TestAdminService_Update_RepoUpdateErr(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	repo.updateErr = errors.New("update failed")
	rn := "x"
	_, err := svc.Update(ctx, u.ID, UpdateAdminInput{RealName: &rn})
	if err == nil || err.Error() != "update failed" {
		t.Errorf("err = %v, want update failed", err)
	}
}

func TestAdminService_Update_FindByUsernameErr(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	repo.findByUsernameErr = errors.New("username lookup failed")
	newUser := "u2"
	_, err := svc.Update(ctx, u.ID, UpdateAdminInput{Username: &newUser})
	if err == nil || err.Error() != "username lookup failed" {
		t.Errorf("err = %v, want username lookup failed", err)
	}
}

func TestAdminService_Update_FindByEmailErr(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw", Email: "a@b.com"})
	repo.findByEmailErr = errors.New("email lookup failed")
	newEmail := "b@c.com"
	_, err := svc.Update(ctx, u.ID, UpdateAdminInput{Email: &newEmail})
	if err == nil || err.Error() != "email lookup failed" {
		t.Errorf("err = %v, want email lookup failed", err)
	}
}

func TestAdminService_Update_EmailSameNoConflict(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw", Email: "a@b.com"})
	sameEmail := "a@b.com"
	updated, err := svc.Update(ctx, u.ID, UpdateAdminInput{Email: &sameEmail})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.Email != "a@b.com" {
		t.Errorf("Email = %q, want a@b.com", updated.Email)
	}
}

func TestAdminService_Update_OnlyPhone(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	phone := "123456"
	updated, err := svc.Update(ctx, u.ID, UpdateAdminInput{Phone: &phone})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.Phone != "123456" {
		t.Errorf("Phone = %q, want 123456", updated.Phone)
	}
}

func TestAdminService_Update_OnlyStatus(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	newStatus := int8(0)
	updated, err := svc.Update(ctx, u.ID, UpdateAdminInput{Status: &newStatus})
	if err != nil {
		t.Fatalf("Update() err = %v", err)
	}
	if updated.Status != 0 {
		t.Errorf("Status = %d, want 0", updated.Status)
	}
}

func TestAdminService_Delete_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	if err := svc.Delete(ctx, u.ID); err != nil {
		t.Fatalf("Delete() err = %v", err)
	}
	_, err := svc.GetByID(ctx, u.ID)
	if err != domainerrors.ErrAdminNotFound {
		t.Errorf("GetByID after delete: err = %v, want ErrAdminNotFound", err)
	}
}

func TestAdminService_Delete_NotFound(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	err := svc.Delete(ctx, 999)
	if err != domainerrors.ErrAdminNotFound {
		t.Errorf("err = %v, want ErrAdminNotFound", err)
	}
}

func TestAdminService_Delete_RepoErr(t *testing.T) {
	repo := newMockAdminRepo()
	_ = repo.Create(context.Background(), &model.Admin{Username: "u", Password: "x"})
	repo.deleteErr = errors.New("delete failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	got, _ := svc.GetByID(context.Background(), 1)
	if got == nil {
		t.Fatal("need admin with ID 1")
	}
	err := svc.Delete(ctx, got.ID)
	if err == nil || err.Error() != "delete failed" {
		t.Errorf("err = %v, want delete failed", err)
	}
}

func TestAdminService_Delete_FindByIDErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.findByIDErr = errors.New("find failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	err := svc.Delete(ctx, 1)
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestAdminService_ListPage_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	_, _ = svc.Create(ctx, CreateAdminInput{Username: "u2", Password: "pw"})
	res, err := svc.ListPage(ctx, 1, 10, nil)
	if err != nil {
		t.Fatalf("ListPage() err = %v", err)
	}
	if res.Total != 2 {
		t.Errorf("Total = %d, want 2", res.Total)
	}
	if res.Pages != 1 {
		t.Errorf("Pages = %d, want 1", res.Pages)
	}
	if res.Current != 1 {
		t.Errorf("Current = %d, want 1", res.Current)
	}
	if res.Size != 10 {
		t.Errorf("Size = %d, want 10", res.Size)
	}
}

func TestAdminService_ListPage_NormalizePage(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	res, err := svc.ListPage(ctx, 0, 10, nil)
	if err != nil {
		t.Fatalf("ListPage() err = %v", err)
	}
	if res.Current != 1 {
		t.Errorf("Current = %d, want 1 (page 0 normalized to 1)", res.Current)
	}
}

func TestAdminService_ListPage_NormalizePageSize(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	res, err := svc.ListPage(ctx, 1, 0, nil)
	if err != nil {
		t.Fatalf("ListPage() err = %v", err)
	}
	if res.Size != 10 {
		t.Errorf("Size = %d, want 10 (pageSize 0 normalized)", res.Size)
	}
}

func TestAdminService_ListPage_PageSizeOver100(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	res, err := svc.ListPage(ctx, 1, 200, nil)
	if err != nil {
		t.Fatalf("ListPage() err = %v", err)
	}
	if res.Size != 10 {
		t.Errorf("Size = %d, want 10 (pageSize >100 capped)", res.Size)
	}
}

func TestAdminService_ListPage_PagesCalc(t *testing.T) {
	repo := newMockAdminRepo()
	for i := 0; i < 25; i++ {
		_ = repo.Create(context.Background(), &model.Admin{Username: fmt.Sprintf("u%d", i), Password: "x"})
	}
	svc := NewAdminService(repo)
	ctx := context.Background()

	res, err := svc.ListPage(ctx, 1, 10, nil)
	if err != nil {
		t.Fatalf("ListPage() err = %v", err)
	}
	if res.Pages != 3 {
		t.Errorf("Pages = %d, want 3 (25/10+1)", res.Pages)
	}
}

func TestAdminService_ListPage_RepoErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.paginateErr = errors.New("paginate failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	_, err := svc.ListPage(ctx, 1, 10, nil)
	if err == nil || err.Error() != "paginate failed" {
		t.Errorf("err = %v, want paginate failed", err)
	}
}

func TestAdminService_BatchDelete_Empty(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	if err := svc.BatchDelete(ctx, nil); err != nil {
		t.Errorf("BatchDelete(nil) err = %v", err)
	}
	if err := svc.BatchDelete(ctx, []uint{}); err != nil {
		t.Errorf("BatchDelete([]) err = %v", err)
	}
}

func TestAdminService_BatchDelete_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u1, _ := svc.Create(ctx, CreateAdminInput{Username: "u1", Password: "pw"})
	u2, _ := svc.Create(ctx, CreateAdminInput{Username: "u2", Password: "pw"})
	if err := svc.BatchDelete(ctx, []uint{u1.ID, u2.ID}); err != nil {
		t.Fatalf("BatchDelete() err = %v", err)
	}
	_, err := svc.GetByID(ctx, u1.ID)
	if err != domainerrors.ErrAdminNotFound {
		t.Errorf("GetByID u1: err = %v", err)
	}
}

func TestAdminService_BatchUpdateStatus_Empty(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	if err := svc.BatchUpdateStatus(ctx, nil, 1); err != nil {
		t.Errorf("BatchUpdateStatus(nil) err = %v", err)
	}
}

func TestAdminService_BatchUpdateStatus_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "pw"})
	if err := svc.BatchUpdateStatus(ctx, []uint{u.ID}, 0); err != nil {
		t.Fatalf("BatchUpdateStatus() err = %v", err)
	}
	got, _ := svc.GetByID(ctx, u.ID)
	if got.Status != 0 {
		t.Errorf("Status = %d, want 0", got.Status)
	}
}

func TestAdminService_ResetPassword_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "old"})
	if err := svc.ResetPassword(ctx, u.ID, "newpass"); err != nil {
		t.Fatalf("ResetPassword() err = %v", err)
	}
	got, _ := svc.GetByID(ctx, u.ID)
	if got.Password == "newpass" {
		t.Error("Password should be hashed")
	}
	if got.Password == "old" {
		t.Error("Password should have changed")
	}
}

func TestAdminService_ResetPassword_NotFound(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	err := svc.ResetPassword(ctx, 999, "new")
	if err != domainerrors.ErrAdminNotFound {
		t.Errorf("err = %v, want ErrAdminNotFound", err)
	}
}

func TestAdminService_ResetPassword_RepoUpdateErr(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	u, _ := svc.Create(ctx, CreateAdminInput{Username: "u", Password: "old"})
	repo.updateErr = errors.New("update failed")
	err := svc.ResetPassword(ctx, u.ID, "newpass")
	if err == nil || err.Error() != "update failed" {
		t.Errorf("err = %v, want update failed", err)
	}
}

func TestAdminService_ResetPassword_FindByIDErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.findByIDErr = errors.New("find failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	err := svc.ResetPassword(ctx, 1, "newpass")
	if err == nil || err.Error() != "find failed" {
		t.Errorf("err = %v, want find failed", err)
	}
}

func TestAdminService_UpdateLastLoginAt_Success(t *testing.T) {
	repo := newMockAdminRepo()
	svc := NewAdminService(repo)
	ctx := context.Background()

	if err := svc.UpdateLastLoginAt(ctx, 1); err != nil {
		t.Errorf("UpdateLastLoginAt() err = %v", err)
	}
}

func TestAdminService_UpdateLastLoginAt_RepoErr(t *testing.T) {
	repo := newMockAdminRepo()
	repo.updateLoginErr = errors.New("login update failed")
	svc := NewAdminService(repo)
	ctx := context.Background()

	err := svc.UpdateLastLoginAt(ctx, 1)
	if err == nil || err.Error() != "login update failed" {
		t.Errorf("err = %v, want login update failed", err)
	}
}

func TestAdminService_GetPermissionsByRole_SuperAdmin(t *testing.T) {
	svc := NewAdminService((*mockAdminRepo)(nil))
	perms := svc.GetPermissionsByRole(1)
	if len(perms) < 2 {
		t.Errorf("len(perms) = %d, want >=2", len(perms))
	}
	found := false
	for _, p := range perms {
		if p == "*:*" {
			found = true
			break
		}
	}
	if !found {
		t.Error("SuperAdmin should have *:* permission")
	}
}

func TestAdminService_GetPermissionsByRole_Regular(t *testing.T) {
	svc := NewAdminService((*mockAdminRepo)(nil))
	perms := svc.GetPermissionsByRole(2)
	if len(perms) != 2 {
		t.Errorf("len(perms) = %d, want 2", len(perms))
	}
	hasDashboard := false
	hasUserView := false
	for _, p := range perms {
		if p == "dashboard:view" {
			hasDashboard = true
		}
		if p == "user:view" {
			hasUserView = true
		}
	}
	if !hasDashboard || !hasUserView {
		t.Errorf("perms = %v", perms)
	}
}

// Ensure mockAdminRepo implements repository.AdminRepository
var _ repository.AdminRepository = (*mockAdminRepo)(nil)

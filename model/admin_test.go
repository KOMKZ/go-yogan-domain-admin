package model

import (
	"testing"
)

func TestAdmin_TableName(t *testing.T) {
	var a Admin
	if got := a.TableName(); got != "admins" {
		t.Errorf("TableName() = %q, want %q", got, "admins")
	}
}

func TestAdmin_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status int8
		want   bool
	}{
		{"active when status 1", 1, true},
		{"inactive when status 0", 0, false},
		{"inactive when status 2", 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Admin{Status: tt.status}
			if got := a.IsActive(); got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdmin_IsSuperAdmin(t *testing.T) {
	tests := []struct {
		name string
		role int8
		want bool
	}{
		{"super admin when role 1", 1, true},
		{"not super admin when role 2", 2, false},
		{"not super admin when role 0", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Admin{Role: tt.role}
			if got := a.IsSuperAdmin(); got != tt.want {
				t.Errorf("IsSuperAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdmin_GetID(t *testing.T) {
	want := uint(42)
	a := &Admin{ID: want}
	if got := a.GetID(); got != want {
		t.Errorf("GetID() = %v, want %v", got, want)
	}
}

func TestAdmin_GetEmail(t *testing.T) {
	want := "admin@example.com"
	a := &Admin{Email: want}
	if got := a.GetEmail(); got != want {
		t.Errorf("GetEmail() = %q, want %q", got, want)
	}
}

func TestAdmin_GetPasswordHash(t *testing.T) {
	want := "$2a$10$hashed"
	a := &Admin{Password: want}
	if got := a.GetPasswordHash(); got != want {
		t.Errorf("GetPasswordHash() = %q, want %q", got, want)
	}
}

func TestAdmin_GetID_Zero(t *testing.T) {
	a := &Admin{}
	if got := a.GetID(); got != 0 {
		t.Errorf("GetID() = %v, want 0", got)
	}
}

func TestAdmin_GetEmail_Empty(t *testing.T) {
	a := &Admin{}
	if got := a.GetEmail(); got != "" {
		t.Errorf("GetEmail() = %q, want empty", got)
	}
}

func TestAdmin_GetPasswordHash_Empty(t *testing.T) {
	a := &Admin{}
	if got := a.GetPasswordHash(); got != "" {
		t.Errorf("GetPasswordHash() = %q, want empty", got)
	}
}

func TestAdmin_GetAvatarURL(t *testing.T) {
	baseURL := "/api/admin/files/storage"
	tests := []struct {
		name            string
		avatar          string
		avatarStorageID string
		want            string
	}{
		{"explicit avatar URL", "https://cdn.example.com/avatar.jpg", "", "https://cdn.example.com/avatar.jpg"},
		{"derived from storage ID", "", "local:avatar@abc123.jpg", "/api/admin/files/storage/local:avatar@abc123.jpg"},
		{"avatar takes priority over storage ID", "https://cdn.example.com/a.jpg", "local:avatar@b.jpg", "https://cdn.example.com/a.jpg"},
		{"both empty returns empty", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Admin{Avatar: tt.avatar, AvatarStorageID: tt.avatarStorageID}
			if got := a.GetAvatarURL(baseURL); got != tt.want {
				t.Errorf("GetAvatarURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

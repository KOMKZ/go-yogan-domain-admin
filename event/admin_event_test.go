package event

import (
	"testing"
)

func TestAdminLoginEvent_Name(t *testing.T) {
	e := NewAdminLoginEvent(1, "admin", "127.0.0.1", "Mozilla")
	if got := e.Name(); got != EventAdminLogin {
		t.Errorf("Name() = %q, want %q", got, EventAdminLogin)
	}
}

func TestAdminLoginEvent_Fields(t *testing.T) {
	e := NewAdminLoginEvent(42, "user1", "10.0.0.1", "Chrome")
	if e.AdminID != 42 {
		t.Errorf("AdminID = %d, want 42", e.AdminID)
	}
	if e.Username != "user1" {
		t.Errorf("Username = %q, want user1", e.Username)
	}
	if e.IP != "10.0.0.1" {
		t.Errorf("IP = %q, want 10.0.0.1", e.IP)
	}
	if e.UserAgent != "Chrome" {
		t.Errorf("UserAgent = %q, want Chrome", e.UserAgent)
	}
}

func TestAdminLogoutEvent_Name(t *testing.T) {
	e := NewAdminLogoutEvent(1, "admin")
	if got := e.Name(); got != EventAdminLogout {
		t.Errorf("Name() = %q, want %q", got, EventAdminLogout)
	}
}

func TestAdminLogoutEvent_Fields(t *testing.T) {
	e := NewAdminLogoutEvent(99, "user2")
	if e.AdminID != 99 {
		t.Errorf("AdminID = %d, want 99", e.AdminID)
	}
	if e.Username != "user2" {
		t.Errorf("Username = %q, want user2", e.Username)
	}
}

func TestEventConstants(t *testing.T) {
	if EventAdminLogin != "admin:login" {
		t.Errorf("EventAdminLogin = %q, want admin:login", EventAdminLogin)
	}
	if EventAdminLogout != "admin:logout" {
		t.Errorf("EventAdminLogout = %q, want admin:logout", EventAdminLogout)
	}
}

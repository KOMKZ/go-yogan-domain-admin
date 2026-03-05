package errors

import (
	"errors"
	"testing"
)

func TestErrAdminNotFound(t *testing.T) {
	if ErrAdminNotFound == nil {
		t.Fatal("ErrAdminNotFound should not be nil")
	}
	if !errors.Is(ErrAdminNotFound, ErrAdminNotFound) {
		t.Error("errors.Is should match ErrAdminNotFound")
	}
	if ErrAdminNotFound.Error() != "admin not found" {
		t.Errorf("Error() = %q, want admin not found", ErrAdminNotFound.Error())
	}
}

func TestErrUsernameExists(t *testing.T) {
	if ErrUsernameExists == nil {
		t.Fatal("ErrUsernameExists should not be nil")
	}
	if ErrUsernameExists.Error() != "username already exists" {
		t.Errorf("Error() = %q, want username already exists", ErrUsernameExists.Error())
	}
}

func TestErrEmailExists(t *testing.T) {
	if ErrEmailExists == nil {
		t.Fatal("ErrEmailExists should not be nil")
	}
	if ErrEmailExists.Error() != "email already exists" {
		t.Errorf("Error() = %q, want email already exists", ErrEmailExists.Error())
	}
}

package errors

import "errors"

var (
	ErrAdminNotFound  = errors.New("admin not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
)

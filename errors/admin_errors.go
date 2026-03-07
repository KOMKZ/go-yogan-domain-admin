package errors

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

const ModuleAdmin = 24

var (
	ErrAdminNotFound = errcode.Register(errcode.New(
		ModuleAdmin, 1001, "admin",
		"error.admin.not_found", "管理员不存在",
		http.StatusNotFound,
	))
	ErrUsernameExists = errcode.Register(errcode.New(
		ModuleAdmin, 1002, "admin",
		"error.admin.username_exists", "用户名已存在",
		http.StatusConflict,
	))
	ErrEmailExists = errcode.Register(errcode.New(
		ModuleAdmin, 1003, "admin",
		"error.admin.email_exists", "邮箱已存在",
		http.StatusConflict,
	))
)

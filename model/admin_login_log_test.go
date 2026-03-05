package model

import (
	"testing"
)

func TestAdminLoginLog_TableName(t *testing.T) {
	var l AdminLoginLog
	if got := l.TableName(); got != "admin_login_logs" {
		t.Errorf("TableName() = %q, want %q", got, "admin_login_logs")
	}
}

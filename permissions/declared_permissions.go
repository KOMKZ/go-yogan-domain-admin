package permissions

type DeclaredPermission struct {
	PermissionCode string
	PermissionName string
	PermissionType string
	ResourceCode   string
	GroupCode      string
	Description    string
}

func DeclaredPermissions() []DeclaredPermission {
	return []DeclaredPermission{
		{
			PermissionCode: "admin:read",
			PermissionName: "查看管理员",
			PermissionType: "READ",
			ResourceCode:   "admin",
			GroupCode:      "SYSTEM",
			Description:    "管理员列表与详情查看",
		},
		{
			PermissionCode: "admin:write",
			PermissionName: "管理管理员",
			PermissionType: "WRITE",
			ResourceCode:   "admin",
			GroupCode:      "SYSTEM",
			Description:    "管理员新增、编辑、删除",
		},
		{
			PermissionCode: "login_log:read",
			PermissionName: "查看管理员登录日志",
			PermissionType: "READ",
			ResourceCode:   "login_log",
			GroupCode:      "SYSTEM",
			Description:    "管理员登录日志查看",
		},
	}
}

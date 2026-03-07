package event

const (
	EventAdminLogin  = "admin:login"
	EventAdminLogout = "admin:logout"
)

type AdminLoginEvent struct {
	name       string
	AdminID    uint
	Username   string
	AdminEmail string
	IP         string
	UserAgent  string
}

func NewAdminLoginEvent(adminID uint, username, email, ip, userAgent string) *AdminLoginEvent {
	return &AdminLoginEvent{
		name:       EventAdminLogin,
		AdminID:    adminID,
		Username:   username,
		AdminEmail: email,
		IP:         ip,
		UserAgent:  userAgent,
	}
}

func (e *AdminLoginEvent) Name() string {
	return e.name
}

type AdminLogoutEvent struct {
	name     string
	AdminID  uint
	Username string
}

func NewAdminLogoutEvent(adminID uint, username string) *AdminLogoutEvent {
	return &AdminLogoutEvent{
		name:     EventAdminLogout,
		AdminID:  adminID,
		Username: username,
	}
}

func (e *AdminLogoutEvent) Name() string {
	return e.name
}

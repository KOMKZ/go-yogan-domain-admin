package model

import (
	"time"
)

type Admin struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Username    string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	RealName    string     `gorm:"size:50" json:"real_name"`
	Email       string     `gorm:"size:100;uniqueIndex" json:"email"`
	Phone       string     `gorm:"size:20" json:"phone"`
	Role        int8       `gorm:"default:2" json:"role"`   // 1=超级管理员, 2=普通管理员
	Status      int8       `gorm:"default:1" json:"status"` // 0=禁用, 1=启用
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Admin) TableName() string {
	return "admins"
}

func (a *Admin) IsActive() bool {
	return a.Status == 1
}

func (a *Admin) IsSuperAdmin() bool {
	return a.Role == 1
}

// Authenticatable interface implementation

func (a *Admin) GetID() uint {
	return a.ID
}

func (a *Admin) GetEmail() string {
	return a.Email
}

func (a *Admin) GetPasswordHash() string {
	return a.Password
}

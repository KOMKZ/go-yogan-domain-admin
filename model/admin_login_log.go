package model

import (
	"time"
)

type AdminLoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	IP        string    `gorm:"column:ip;size:50;not null" json:"ip"`
	UserAgent string    `gorm:"column:user_agent;size:500" json:"user_agent"`
	DeviceID  string    `gorm:"column:device_id;size:100" json:"device_id"`
	City      string    `gorm:"size:100" json:"city"`
	Country   string    `gorm:"size:100" json:"country"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AdminLoginLog) TableName() string {
	return "admin_login_logs"
}

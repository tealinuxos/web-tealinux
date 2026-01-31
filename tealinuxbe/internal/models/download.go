package models

import "time"

type Download struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Edition   string    `gorm:"size:50;not null;index" json:"edition"` // COSMIC, PLASMA
	IPAddress string    `gorm:"size:45" json:"ip_address"`             // IPv4 or IPv6
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	UserID    *uint     `json:"user_id"` // Optional, if authenticated
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

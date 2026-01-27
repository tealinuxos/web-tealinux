package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:120" json:"name"`
	Email        string    `gorm:"uniqueIndex;size:120" json:"email"`
	Password     string    `gorm:"size:255" json:"-"`
	Provider     string    `gorm:"size:50" json:"provider"` // local, google, github
	Avatar       string    `gorm:"size:255" json:"avatar"`
	Role         string    `gorm:"default:user" json:"role"`
	RefreshToken string    `gorm:"size:500" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

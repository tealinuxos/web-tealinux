package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:120"`
	Email        string `gorm:"uniqueIndex;size:120"`
	Password     string `gorm:"size:255"`
	Provider     string `gorm:"size:50"` // local, google, github
	Avatar       string `gorm:"size:255"`
	Role         string `gorm:"default:user"`
	RefreshToken string `gorm:"size:500"`
	CreatedAt    time.Time
}

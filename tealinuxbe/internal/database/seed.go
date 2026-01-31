package database

import (
	"log"
	"tealinux-api/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), 14)
		if err != nil {
			log.Println("Failed to hash admin password:", err)
			return
		}

		admin := models.User{
			Name:     "TeaLinux Admin",
			Email:    "admin@tealinux.org",
			Password: string(hash),
			Role:     "admin",
			Provider: "local",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Failed to seed admin:", err)
		} else {
			log.Println("Admin user seeded: admin@tealinux.org / admin123")
		}
	}
}

package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"tealinux-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("DB ERROR:", err)
	}

	// Test connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("DB PING ERROR:", err)
	}

	// Auto migrate - User and Download models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Download{},
	); err != nil {
		log.Fatal("MIGRATION ERROR:", err)
	}

	DB = db
	log.Println("Database connected & migrated")

	// Seed initial admin data
	SeedAdmin()
}

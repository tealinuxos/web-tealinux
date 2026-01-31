package handlers

import (
	"time"

	"tealinux-api/internal/database"
	"tealinux-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

type TrackDownloadRequest struct {
	Edition string `json:"edition"`
}

func TrackDownload(c *fiber.Ctx) error {
	var body TrackDownloadRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.Edition == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Edition is required"})
	}

	// Validate edition
	validEditions := map[string]bool{"COSMIC": true, "PLASMA": true}
	if !validEditions[body.Edition] {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid edition. Must be COSMIC or PLASMA"})
	}

	// Get optional user ID from context (if authenticated)
	var userID *uint
	if id := c.Locals("user_id"); id != nil {
		switch v := id.(type) {
		case float64:
			uid := uint(v)
			userID = &uid
		case uint:
			userID = &v
		case int:
			uid := uint(v)
			userID = &uid
		}
	}

	download := models.Download{
		Edition:   body.Edition,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&download).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to record download"})
	}

	return c.JSON(fiber.Map{
		"message": "Download tracked successfully",
		"id":      download.ID,
	})
}

type DownloadStats struct {
	TotalDownloads     int64             `json:"total_downloads"`
	DownloadsByEdition map[string]int64  `json:"downloads_by_edition"`
	DownloadsToday     int64             `json:"downloads_today"`
	DownloadsThisWeek  int64             `json:"downloads_this_week"`
	DownloadsThisMonth int64             `json:"downloads_this_month"`
	RecentDownloads    []models.Download `json:"recent_downloads"`
}

func GetDownloadStats(c *fiber.Ctx) error {
	var stats DownloadStats

	// Total downloads
	database.DB.Model(&models.Download{}).Count(&stats.TotalDownloads)

	// Downloads by edition
	stats.DownloadsByEdition = make(map[string]int64)
	var editionCounts []struct {
		Edition string
		Count   int64
	}
	database.DB.Model(&models.Download{}).
		Select("edition, count(*) as count").
		Group("edition").
		Scan(&editionCounts)

	for _, ec := range editionCounts {
		stats.DownloadsByEdition[ec.Edition] = ec.Count
	}

	// Time-based stats
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := startOfDay.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	database.DB.Model(&models.Download{}).Where("created_at >= ?", startOfDay).Count(&stats.DownloadsToday)
	database.DB.Model(&models.Download{}).Where("created_at >= ?", startOfWeek).Count(&stats.DownloadsThisWeek)
	database.DB.Model(&models.Download{}).Where("created_at >= ?", startOfMonth).Count(&stats.DownloadsThisMonth)

	// Recent downloads (last 10)
	database.DB.Order("created_at desc").Limit(10).Find(&stats.RecentDownloads)

	return c.JSON(stats)
}

type DailyDownload struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

func GetDownloadHistory(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30) // Default to 30 days
	if days > 365 {
		days = 365
	}

	startDate := time.Now().AddDate(0, 0, -days)

	var dailyDownloads []DailyDownload
	database.DB.Model(&models.Download{}).
		Select("DATE(created_at) as date, count(*) as count").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date asc").
		Scan(&dailyDownloads)

	return c.JSON(dailyDownloads)
}

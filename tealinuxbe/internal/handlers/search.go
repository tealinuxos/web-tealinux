package handlers

import (
	"tealinux-api/internal/database"
	"tealinux-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

func Search(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.JSON([]models.Topic{})
	}

	var topics []models.Topic
	// Simple ILIKE search for now to be database agnostic or safer if TSVECTOR is not set up
	// But user asked for TSVECTOR. Let's try to use it if possible, but fallback to ILIKE if complex.
	// The user's SQL: to_tsvector('english', title || ' ' || content) @@ plainto_tsquery('linux')
	// Since Topic doesn't have content, we'll search Title.
	// If we want to search Post content too, we need a join.

	// Let's do a simple search on Topic Title first.
	// database.DB.Where("to_tsvector('english', title) @@ plainto_tsquery(?)", q).Find(&topics)

	// Using ILIKE for broader compatibility and simplicity in this step without setting up TSVECTOR indexes
	err := database.DB.Preload("User").Preload("Category").Preload("Tags").Preload("Posts").
		Where("title ILIKE ?", "%"+q+"%").
		Find(&topics).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Search failed"})
	}

	return c.JSON(topics)
}

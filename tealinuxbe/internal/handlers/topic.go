package handlers

import (
	"tealinux-api/internal/database"
	"tealinux-api/internal/models"
	"tealinux-api/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetTopics(c *fiber.Ctx) error {
	var topics []models.Topic
	query := database.DB.Preload("User").Preload("Category").Preload("Tags").Order("created_at desc")

	if catID := c.Query("category_id"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}

	if err := query.Find(&topics).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch topics"})
	}
	return c.JSON(topics)
}

func GetTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	var topic models.Topic
	if err := database.DB.Preload("User").Preload("Category").Preload("Tags").Preload("Posts.User").Where("id = ? OR slug = ?", id, id).First(&topic).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Topic not found"})
	}

	// Increment view count
	go func() {
		database.DB.Model(&topic).UpdateColumn("views", topic.Views+1)
	}()

	return c.JSON(topic)
}

func CreateTopic(c *fiber.Ctx) error {
	type Req struct {
		Title      string   `json:"title"`
		CategoryID string   `json:"category_id"`
		Content    string   `json:"content"`
		Tags       []string `json:"tags"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	userID := c.Locals("user_id")
	// Ensure userID is uint
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	case int:
		uid = uint(v)
	default:
		// Fallback or error if type is unexpected
		// Assuming middleware ensures it's a valid number
	}

	catUUID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	topic := models.Topic{
		Title:      body.Title,
		Slug:       utils.MakeSlug(body.Title),
		CategoryID: catUUID,
		UserID:     uid,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&topic).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create topic"})
	}

	// Handle Tags
	if len(body.Tags) > 0 {
		var tags []models.Tag
		for _, tagName := range body.Tags {
			var tag models.Tag
			// Find or create tag
			if err := tx.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				continue
			}
			tags = append(tags, tag)
		}
		if len(tags) > 0 {
			if err := tx.Model(&topic).Association("Tags").Replace(tags); err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to save tags"})
			}
		}
	}

	post := models.Post{
		TopicID:   topic.ID,
		UserID:    uid,
		Content:   body.Content,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create initial post"})
	}

	tx.Commit()

	// Reload topic to include tags
	database.DB.Preload("Tags").First(&topic, topic.ID)

	return c.JSON(fiber.Map{"topic": topic, "post": post})
}

func UpdateTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	var topic models.Topic
	if err := database.DB.Where("id = ?", id).First(&topic).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Topic not found"})
	}

	// Check ownership or admin
	userID := c.Locals("user_id")
	role := c.Locals("role")

	// Convert userID to uint safely
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	case int:
		uid = uint(v)
	}

	if topic.UserID != uid && role != "admin" && role != "moderator" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Req struct {
		Title string `json:"title"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	if body.Title != "" {
		topic.Title = body.Title
		topic.Slug = utils.MakeSlug(body.Title)
	}

	database.DB.Save(&topic)
	return c.JSON(topic)
}

func DeleteTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	var topic models.Topic
	if err := database.DB.Where("id = ?", id).First(&topic).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Topic not found"})
	}

	userID := c.Locals("user_id")
	role := c.Locals("role")

	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	case int:
		uid = uint(v)
	}

	if topic.UserID != uid && role != "admin" && role != "moderator" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Delete topic (and cascading posts if configured in DB, otherwise manual)
	// GORM usually handles cascade if defined, but here we didn't define constraint explicitly in struct tags for cascade delete
	// Let's manually delete posts first or rely on DB. For now, manual.
	database.DB.Where("topic_id = ?", topic.ID).Delete(&models.Post{})
	database.DB.Delete(&topic)

	return c.JSON(fiber.Map{"message": "Topic deleted"})
}

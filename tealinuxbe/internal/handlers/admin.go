package handlers

import (
	"tealinux-api/internal/database"
	"tealinux-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Dashboard Stats
func AdminDashboardStats(c *fiber.Ctx) error {
	var userCount int64
	var topicCount int64
	var categoryCount int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Topic{}).Count(&topicCount)
	database.DB.Model(&models.Category{}).Count(&categoryCount)

	return c.JSON(fiber.Map{
		"users":      userCount,
		"topics":     topicCount,
		"categories": categoryCount,
	})
}

// User Management
func AdminListUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Find(&users)
	return c.JSON(users)
}

func AdminUpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	type Req struct {
		Role string `json:"role"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "bad request"})
	}

	if err := database.DB.Model(&models.User{}).Where("id = ?", id).Update("role", body.Role).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update role"})
	}

	return c.JSON(fiber.Map{"message": "role updated"})
}

func AdminDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.User{}, id)
	return c.JSON(fiber.Map{"message": "user deleted"})
}

// Thread & Post Control
func AdminLockTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Model(&models.Topic{}).Where("id = ?", id).Update("is_locked", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to lock topic"})
	}
	return c.JSON(fiber.Map{"message": "topic locked"})
}

func AdminPinTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Model(&models.Topic{}).Where("id = ?", id).Update("is_pinned", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to pin topic"})
	}
	return c.JSON(fiber.Map{"message": "topic pinned"})
}

func AdminDeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Delete(&models.Post{}, "id = ?", id)
	return c.JSON(fiber.Map{"message": "post deleted"})
}

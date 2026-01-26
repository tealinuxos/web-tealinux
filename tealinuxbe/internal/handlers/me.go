package handlers

import (
	"tealinux-api/internal/database"
	"tealinux-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

func Me(c *fiber.Ctx) error {
	id := c.Locals("user_id")

	var user models.User
	database.DB.First(&user, id)

	return c.JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"role":   user.Role,
		"avatar": user.Avatar,
	})
}

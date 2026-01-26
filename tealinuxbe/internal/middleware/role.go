package middleware

import "github.com/gofiber/fiber/v2"

func RoleOnly(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")

		if userRole != role {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

		return c.Next()
	}
}

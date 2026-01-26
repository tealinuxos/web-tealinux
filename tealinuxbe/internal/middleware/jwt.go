package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		if auth == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token format"})
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["id"])
		c.Locals("id", claims["id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

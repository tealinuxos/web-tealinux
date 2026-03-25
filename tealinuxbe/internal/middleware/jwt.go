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
		tokenStr := ""

		if auth != "" {
			parts := strings.Split(auth, " ")
			if len(parts) == 2 {
				tokenStr = parts[1]
			}
		}

		// Fallback to cookie if Authorization header is missing or invalid
		// This is important for web clients that use HttpOnly cookies
		if tokenStr == "" {
			tokenStr = c.Cookies("tealinux_access_token")
		}

		if tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

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

package routes

import (
	"tealinux-api/internal/handlers"
	"tealinux-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/health", handlers.Health)

	// Auth Routes
	app.Post("/auth/register", handlers.Register)
	app.Post("/auth/login", handlers.Login)

	app.Get("/auth/google", handlers.GoogleLogin)
	app.Get("/auth/google/callback", handlers.GoogleCallback)

	app.Get("/auth/github", handlers.GithubLogin)
	app.Get("/auth/github/callback", handlers.GithubCallback)

	app.Post("/auth/refresh", handlers.RefreshToken)

	// Download Tracking (public - allows anonymous tracking)
	app.Post("/api/downloads/track", handlers.TrackDownload)

	// Protected Routes
	api := app.Group("/api", middleware.JWTProtected())
	api.Get("/me", handlers.Me)
	api.Post("/logout", handlers.Logout)

	// Admin Routes
	admin := api.Group("/admin", middleware.AdminOnly())
	admin.Get("/downloads/stats", handlers.GetDownloadStats)
	admin.Get("/downloads/history", handlers.GetDownloadHistory)
}

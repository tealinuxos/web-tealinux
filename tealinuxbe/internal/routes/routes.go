package routes

import (
	"tealinux-api/internal/handlers"
	"tealinux-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/health", handlers.Health)

	app.Post("/auth/register", handlers.Register)
	app.Post("/auth/login", handlers.Login)

	app.Get("/auth/google", handlers.GoogleLogin)
	app.Get("/auth/google/callback", handlers.GoogleCallback)

	app.Get("/auth/github", handlers.GithubLogin)
	app.Get("/auth/github/callback", handlers.GithubCallback)

	app.Post("/auth/refresh", handlers.RefreshToken)

	api := app.Group("/api", middleware.JWTProtected())
	api.Get("/me", handlers.Me)
	api.Post("/logout", handlers.Logout)

	// Forum Routes
	// Public
	app.Get("/categories", handlers.GetCategories)
	app.Get("/categories/:id", handlers.GetCategory)
	app.Get("/topics", handlers.GetTopics)
	app.Get("/topics/:id", handlers.GetTopic)
	app.Get("/topics/:id/posts", handlers.GetPostsByTopic)
	app.Get("/search", handlers.Search)

	// Protected (User)
	api.Post("/topics", handlers.CreateTopic)
	api.Put("/topics/:id", handlers.UpdateTopic)
	api.Delete("/topics/:id", handlers.DeleteTopic)

	api.Post("/topics/:id/posts", handlers.CreatePost)
	api.Put("/posts/:id", handlers.UpdatePost)
	api.Delete("/posts/:id", handlers.DeletePost)
	api.Post("/posts/:id/like", handlers.LikePost)
	api.Delete("/posts/:id/like", handlers.UnlikePost)

	admin := api.Group("/admin", middleware.AdminOnly())
	admin.Get("/dashboard", handlers.AdminDashboardStats)

	// Admin - User Management
	admin.Get("/users", handlers.AdminListUsers)
	admin.Patch("/users/:id/role", handlers.AdminUpdateUserRole)
	admin.Delete("/users/:id", handlers.AdminDeleteUser)

	// Admin - Topic & Post Control
	admin.Patch("/topics/:id/lock", handlers.AdminLockTopic)
	admin.Patch("/topics/:id/pin", handlers.AdminPinTopic)
	admin.Delete("/posts/:id", handlers.AdminDeletePost)

	// Admin - Category Management
	admin.Post("/categories", handlers.CreateCategory)
	admin.Put("/categories/:id", handlers.UpdateCategory)
	admin.Delete("/categories/:id", handlers.DeleteCategory)
}

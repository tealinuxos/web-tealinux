package handlers

import (
	"tealinux-api/internal/database"
	"tealinux-api/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetPostsByTopic(c *fiber.Ctx) error {
	topicID := c.Params("id")
	var posts []models.Post
	if err := database.DB.Preload("User").Preload("Likes").Where("topic_id = ?", topicID).Order("created_at asc").Find(&posts).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch posts"})
	}
	return c.JSON(posts)
}

func CreatePost(c *fiber.Ctx) error {
	topicID := c.Params("id")
	type Req struct {
		Content   string `json:"content"`
		ReplyToID string `json:"reply_to_id"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	tUUID, err := uuid.Parse(topicID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid topic ID"})
	}

	userID := c.Locals("user_id")
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	case int:
		uid = uint(v)
	}

	post := models.Post{
		TopicID:   tUUID,
		UserID:    uid,
		Content:   body.Content,
		CreatedAt: time.Now(),
	}

	if body.ReplyToID != "" {
		rUUID, err := uuid.Parse(body.ReplyToID)
		if err == nil {
			post.ReplyToID = &rUUID
		}
	}

	if err := database.DB.Create(&post).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create post"})
	}

	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post
	if err := database.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
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

	if post.UserID != uid && role != "admin" && role != "moderator" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type Req struct {
		Content string `json:"content"`
	}
	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	post.Content = body.Content
	database.DB.Save(&post)
	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	id := c.Params("id")
	var post models.Post
	if err := database.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Post not found"})
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

	if post.UserID != uid && role != "admin" && role != "moderator" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	database.DB.Delete(&post)
	return c.JSON(fiber.Map{"message": "Post deleted"})
}

func LikePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	pUUID, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	userID := c.Locals("user_id")
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	case int:
		uid = uint(v)
	}

	// Check if already liked
	var existingLike models.Like
	if err := database.DB.Where("post_id = ? AND user_id = ?", pUUID, uid).First(&existingLike).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Already liked"})
	}

	like := models.Like{
		PostID:    pUUID,
		UserID:    uid,
		CreatedAt: time.Now(),
	}

	database.DB.Create(&like)
	return c.JSON(fiber.Map{"message": "Liked"})
}

func UnlikePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	pUUID, err := uuid.Parse(postID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid post ID"})
	}

	userID := c.Locals("user_id")
	var uid uint
	switch v := userID.(type) {
	case float64:
		uid = uint(v)
	case uint:
		uid = v
	case int:
		uid = uint(v)
	}

	database.DB.Where("post_id = ? AND user_id = ?", pUUID, uid).Delete(&models.Like{})
	return c.JSON(fiber.Map{"message": "Unliked"})
}

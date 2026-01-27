package handlers

import (
	"tealinux-api/internal/database"
	"tealinux-api/internal/models"
	"tealinux-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Order("\"order\" asc").Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}
	return c.JSON(categories)
}

func GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.Category

	query := database.DB
	if _, err := uuid.Parse(id); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("slug = ?", id)
	}

	if err := query.First(&category).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}
	return c.JSON(category)
}

func CreateCategory(c *fiber.Ctx) error {
	type Req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Order       int    `json:"order"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	category := models.Category{
		Name:        body.Name,
		Slug:        utils.MakeSlug(body.Name),
		Description: body.Description,
		Order:       body.Order,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create category"})
	}

	return c.JSON(category)
}

func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category models.Category
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	type Req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Order       int    `json:"order"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad request"})
	}

	if body.Name != "" {
		category.Name = body.Name
		category.Slug = utils.MakeSlug(body.Name)
	}
	if body.Description != "" {
		category.Description = body.Description
	}
	category.Order = body.Order

	database.DB.Save(&category)
	return c.JSON(category)
}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	// Check if valid UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := database.DB.Delete(&models.Category{}, uid).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete category"})
	}
	return c.JSON(fiber.Map{"message": "Category deleted"})
}

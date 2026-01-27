package database

import (
	"log"
	"tealinux-api/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), 14)
		if err != nil {
			log.Println("Failed to hash admin password:", err)
			return
		}

		admin := models.User{
			Name:     "TeaLinux Admin",
			Email:    "admin@tealinux.org",
			Password: string(hash),
			Role:     "admin",
			Provider: "local",
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Failed to seed admin:", err)
		} else {
			log.Println("Admin user seeded: admin@tealinux.org / admin123")
		}
	}
}

func SeedForumCategories() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)

	if count > 0 {
		log.Println("Forum categories already seeded, skipping...")
		return
	}

	log.Println("Seeding forum categories...")

	categories := []models.Category{
		{
			Name:        "Bug Reports",
			Slug:        "bug-reports",
			Description: "Report bugs and issues you've encountered with TeaLinuxOS",
			Order:       1,
		},
		{
			Name:        "Questions & Help",
			Slug:        "questions-help",
			Description: "Ask questions and get help from the community",
			Order:       2,
		},
		{
			Name:        "General Discussion",
			Slug:        "general-discussion",
			Description: "Discuss anything related to TeaLinuxOS and Linux in general",
			Order:       3,
		},
	}

	for i := range categories {
		if err := DB.Create(&categories[i]).Error; err != nil {
			log.Printf("Failed to create category %s: %v\n", categories[i].Name, err)
		} else {
			log.Printf("Created category: %s\n", categories[i].Name)
		}
	}

	log.Println("Forum categories seeded successfully!")
}

func SeedForumData() {
	// Check if we already have topics
	var topicCount int64
	DB.Model(&models.Topic{}).Count(&topicCount)

	if topicCount > 0 {
		log.Println("Forum topics already seeded, skipping...")
		return
	}

	log.Println("Seeding forum sample data...")

	// Get categories
	var categories []models.Category
	DB.Find(&categories)
	if len(categories) == 0 {
		log.Println("No categories found, please seed categories first")
		return
	}

	// Get or create test user
	var testUser models.User
	DB.Where("email = ?", "testuser@example.com").First(&testUser)
	if testUser.ID == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 14)
		testUser = models.User{
			Name:     "Test User",
			Email:    "testuser@example.com",
			Password: string(hash),
			Role:     "user",
			Provider: "local",
		}
		DB.Create(&testUser)
	}

	// Sample topics
	topics := []struct {
		Title      string
		Slug       string
		Content    string
		CategoryID string
		Tags       []string
	}{
		{
			Title:      "Bug: System crashes when opening Settings",
			Slug:       "bug-system-crashes-when-opening-settings",
			Content:    "I found a critical bug. When I try to open System Settings, the entire system freezes.\n\nSteps to reproduce:\n1. Click on Settings icon\n2. System freezes immediately\n3. Must force restart\n\nExpected: Settings opens normally\nActual: Complete system freeze\n\nSystem: TeaLinuxOS COSMIC, 8GB RAM, Intel HD Graphics",
			CategoryID: categories[0].ID.String(), // Bug Reports
			Tags:       []string{"bug", "critical", "settings", "cosmic"},
		},
		{
			Title:      "How to install TeaLinuxOS on UEFI system?",
			Slug:       "how-to-install-tealinux-on-uefi",
			Content:    "I'm trying to install TeaLinuxOS on my laptop with UEFI. The installer doesn't seem to detect my hard drive. Any help would be appreciated!",
			CategoryID: categories[1].ID.String(), // Questions & Help
			Tags:       []string{"installation", "uefi", "help"},
		},
		{
			Title:      "TeaLinuxOS Performance vs Other Distros",
			Slug:       "tealinux-performance-comparison",
			Content:    "I've been using TeaLinuxOS for a month now. Compared to Ubuntu and Fedora, I'm getting much better performance especially on my older laptop. What's your experience?",
			CategoryID: categories[2].ID.String(), // General Discussion
			Tags:       []string{"performance", "discussion", "comparison"},
		},
	}

	for _, t := range topics {
		// Parse category UUID
		categoryUUID, err := uuid.Parse(t.CategoryID)
		if err != nil {
			log.Printf("Invalid category UUID for topic %s: %v\n", t.Title, err)
			continue
		}

		// Create topic
		topic := models.Topic{
			Title:      t.Title,
			Slug:       t.Slug,
			UserID:     testUser.ID,
			CategoryID: categoryUUID,
			Views:      0,
			IsPinned:   false,
			IsLocked:   false,
		}

		if err := DB.Create(&topic).Error; err != nil {
			log.Printf("Failed to create topic %s: %v\n", t.Title, err)
			continue
		}

		// Create initial post (first post is the topic content)
		post := models.Post{
			TopicID: topic.ID,
			UserID:  testUser.ID,
			Content: t.Content,
		}

		if err := DB.Create(&post).Error; err != nil {
			log.Printf("Failed to create post for topic %s: %v\n", t.Title, err)
			continue
		}

		// Create tags
		for _, tagName := range t.Tags {
			tag := models.Tag{Name: tagName}
			DB.FirstOrCreate(&tag, models.Tag{Name: tagName})

			// Associate tag with topic
			DB.Exec("INSERT INTO topic_tags (topic_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
				topic.ID.String(), tag.ID.String())
		}

		log.Printf("Created topic: %s\n", t.Title)
	}

	log.Println("Forum sample data seeded successfully!")
}

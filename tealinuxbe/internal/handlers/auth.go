package handlers

import (
	"context"
	"encoding/json"
	"os"

	"tealinux-api/internal/database"
	"tealinux-api/internal/models"
	"tealinux-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type UserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Avatar string `json:"avatar"`
}

func Register(c *fiber.Ctx) error {
	type Req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "email and password required"})
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", body.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "email already registered"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
		Provider: "local",
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Avatar: user.Avatar,
	})
}

func Login(c *fiber.Ctx) error {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	if err := database.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "user not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "wrong password"})
	}

	access, refresh, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate token"})
	}

	database.DB.Model(&user).Update("refresh_token", refresh)

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
		"user": UserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Role:   user.Role,
			Avatar: user.Avatar,
		},
	})
}

func getGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func getGithubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

func GoogleLogin(c *fiber.Ctx) error {
	url := getGoogleConfig().AuthCodeURL("state")
	return c.Redirect(url)
}

func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := getGoogleConfig().Exchange(context.Background(), code)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	client := getGoogleConfig().Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to get user info"})
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to parse user info"})
	}

	email, _ := data["email"].(string)
	name, _ := data["name"].(string)
	picture, _ := data["picture"].(string)

	var user models.User
	err = database.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		user = models.User{
			Name:     name,
			Email:    email,
			Provider: "google",
			Avatar:   picture,
			Role:     "user",
		}
		database.DB.Create(&user)
	}

	access, refresh, _ := utils.GenerateToken(user.ID, user.Role)

	database.DB.Model(&user).Update("refresh_token", refresh)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4321"
	}

	redirectURL := frontendURL + "/auth/callback?" +
		"access_token=" + access +
		"&refresh_token=" + refresh +
		"&id=" + utils.UintToString(user.ID) +
		"&name=" + user.Name +
		"&email=" + user.Email +
		"&role=" + user.Role +
		"&avatar=" + user.Avatar

	return c.Redirect(redirectURL)
}

func GithubLogin(c *fiber.Ctx) error {
	url := getGithubConfig().AuthCodeURL("state")
	return c.Redirect(url)
}

func GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	token, err := getGithubConfig().Exchange(context.Background(), code)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	client := getGithubConfig().Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to get user info"})
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "failed to parse user info"})
	}

	email := ""
	if v, ok := data["email"]; ok && v != nil {
		email = v.(string)
	} else {
		// If email is private, we might need to fetch it from /user/emails
		// For now, fallback to login@github.com as per prompt
		if login, ok := data["login"].(string); ok {
			email = login + "@github.com"
		}
	}

	name, _ := data["name"].(string)
	avatar, _ := data["avatar_url"].(string)

	var user models.User
	err = database.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		user = models.User{
			Name:     name,
			Email:    email,
			Provider: "github",
			Avatar:   avatar,
			Role:     "user",
		}
		database.DB.Create(&user)
	}

	access, refresh, _ := utils.GenerateToken(user.ID, user.Role)

	database.DB.Model(&user).Update("refresh_token", refresh)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:4321"
	}

	redirectURL := frontendURL + "/auth/callback?" +
		"access_token=" + access +
		"&refresh_token=" + refresh +
		"&id=" + utils.UintToString(user.ID) +
		"&name=" + user.Name +
		"&email=" + user.Email +
		"&role=" + user.Role +
		"&avatar=" + user.Avatar

	return c.Redirect(redirectURL)
}

func RefreshToken(c *fiber.Ctx) error {
	type Req struct {
		RefreshToken string `json:"refresh_token"`
	}

	var body Req
	c.BodyParser(&body)

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "user not found"})
	}

	if user.RefreshToken == "" {
		return c.Status(401).JSON(fiber.Map{"error": "logged out"})
	}

	if user.RefreshToken != body.RefreshToken {
		return c.Status(401).JSON(fiber.Map{"error": "token revoked"})
	}

	access, refresh, _ := utils.GenerateToken(user.ID, user.Role)

	database.DB.Model(&user).Update("refresh_token", refresh)

	return c.JSON(fiber.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func Logout(c *fiber.Ctx) error {
	id := c.Locals("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "user not found"})
	}

	database.DB.Model(&user).Update("refresh_token", "")

	return c.JSON(fiber.Map{
		"message": "logout success",
	})
}

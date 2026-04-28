package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"golang.org/x/crypto/bcrypt"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

const (
	minUsernameLen = 3
	maxUsernameLen = 30
	minPasswordLen = 8
	maxPasswordLen = 72 // bcrypt input limit
)

func normalizeRegisterRequest(req *models.RegisterRequest) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.FullName = strings.TrimSpace(req.FullName)
}

func validateUsername(username string) string {
	if len(username) < minUsernameLen || len(username) > maxUsernameLen {
		return "Username must be between 3 and 30 characters"
	}
	if !usernamePattern.MatchString(username) {
		return "Username can only contain letters, numbers, underscore, dash, and dot"
	}
	return ""
}

func validatePasswordComplexity(password string) string {
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return "Password must be between 8 and 72 characters"
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return "Password must include uppercase, lowercase, number, and special character"
	}
	return ""
}

// POST /api/v1/auth/register
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	normalizeRegisterRequest(&req)

	if req.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Full name is required"})
		return
	}
	if msg := validateUsername(req.Username); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if msg := validatePasswordComplexity(req.Password); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or username already taken"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("register: bcrypt hash generation failed for email=%q username=%q: %v", req.Email, req.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed. Please try again."})
		return
	}
	user.Password = string(hash)

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("register: user create failed for email=%q username=%q: %v", req.Email, req.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed. Please try again."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// POST /api/v1/auth/login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.ToLower(req.Email)

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}

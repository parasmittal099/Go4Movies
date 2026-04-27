package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"golang.org/x/crypto/bcrypt"
)

// POST /api/v1/auth/register
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.ToLower(req.Email)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hash)

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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

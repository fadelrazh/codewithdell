package handlers

import (
	"net/http"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
	Location string `json:"location"`
}

// GetProfile handles getting user profile
func GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Remove sensitive information
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles updating user profile
func UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Find user
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if username is already taken (if changing)
	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := db.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Reload user data
	db.First(&user, userID)
	user.Password = ""

	c.JSON(http.StatusOK, user)
} 
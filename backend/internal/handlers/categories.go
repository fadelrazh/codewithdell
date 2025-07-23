package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCategoryRequest represents category creation request
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=200"`
	Color       string `json:"color" binding:"max=7"`
	Icon        string `json:"icon" binding:"max=50"`
}

// UpdateCategoryRequest represents category update request
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=50"`
	Description string `json:"description" binding:"max=200"`
	Color       string `json:"color" binding:"max=7"`
	Icon        string `json:"icon" binding:"max=50"`
}

// GetCategories handles getting all categories
func GetCategories(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var categories []models.Category
	if err := db.Order("name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      len(categories),
	})
}

// GetCategoryBySlug handles getting a single category by slug
func GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category
	if err := db.Preload("Posts").Preload("Projects").
		Where("slug = ?", slug).
		First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// CreateCategory handles creating a new category (admin only)
func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Generate slug from name
	slug := generateSlug(req.Name)

	// Check if category already exists
	var existingCategory models.Category
	if err := db.Where("name = ? OR slug = ?", req.Name, slug).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	}

	// Create category
	category := models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": category,
	})
}

// UpdateCategory handles updating a category (admin only)
func UpdateCategory(c *gin.Context) {
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category
	if err := db.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if new name conflicts with existing category
	if req.Name != "" && req.Name != category.Name {
		var existingCategory models.Category
		if err := db.Where("name = ? AND id != ?", req.Name, categoryID).First(&existingCategory).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Category name already exists"})
			return
		}
		category.Name = req.Name
		category.Slug = generateSlug(req.Name)
	}

	// Update other fields
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}

	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category updated successfully",
		"category": category,
	})
}

// DeleteCategory handles deleting a category (admin only)
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category
	if err := db.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category has associated posts or projects
	var postCount int64
	db.Model(&models.Post{}).Joins("JOIN post_categories ON posts.id = post_categories.post_id").
		Where("post_categories.category_id = ?", categoryID).Count(&postCount)

	var projectCount int64
	db.Model(&models.Project{}).Joins("JOIN project_categories ON projects.id = project_categories.project_id").
		Where("project_categories.category_id = ?", categoryID).Count(&projectCount)

	if postCount > 0 || projectCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cannot delete category with associated posts or projects",
			"post_count":    postCount,
			"project_count": projectCount,
		})
		return
	}

	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetCategoryPosts handles getting posts by category
func GetCategoryPosts(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category
	if err := db.Where("slug = ?", slug).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var posts []models.Post
	if err := db.Preload("Author").Preload("Tags").
		Joins("JOIN post_categories ON posts.id = post_categories.post_id").
		Where("post_categories.category_id = ? AND posts.status = ?", category.ID, models.PostStatusPublished).
		Order("posts.created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"posts":    posts,
		"total":    len(posts),
	})
}

// GetCategoryProjects handles getting projects by category
func GetCategoryProjects(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var category models.Category
	if err := db.Where("slug = ?", slug).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var projects []models.Project
	if err := db.Preload("Technologies").Preload("Tags").
		Joins("JOIN project_categories ON projects.id = project_categories.project_id").
		Where("project_categories.category_id = ? AND projects.status = ?", category.ID, models.ProjectStatusPublished).
		Order("projects.created_at DESC").
		Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"projects": projects,
		"total":    len(projects),
	})
}

// generateSlug generates a URL-friendly slug from a string
func generateSlug(input string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(input)
	slug = strings.ReplaceAll(slug, " ", "-")
	
	// Remove special characters except hyphens
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	
	// Remove multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	
	return slug
} 
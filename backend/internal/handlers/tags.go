package handlers

import (
	"net/http"
	"strings"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateTagRequest represents tag creation request
type CreateTagRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=30"`
	Color string `json:"color" binding:"max=7"`
}

// UpdateTagRequest represents tag update request
type UpdateTagRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=30"`
	Color string `json:"color" binding:"max=7"`
}

// GetTags handles getting all tags
func GetTags(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var tags []models.Tag
	if err := db.Order("name ASC").Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"total": len(tags),
	})
}

// GetTagBySlug handles getting a single tag by slug
func GetTagBySlug(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var tag models.Tag
	if err := db.Preload("Posts").Preload("Projects").
		Where("slug = ?", slug).
		First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

// CreateTag handles creating a new tag (admin only)
func CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Generate slug from name
	slug := generateSlug(req.Name)

	// Check if tag already exists
	var existingTag models.Tag
	if err := db.Where("name = ? OR slug = ?", req.Name, slug).First(&existingTag).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tag already exists"})
		return
	}

	// Create tag
	tag := models.Tag{
		Name:  req.Name,
		Slug:  slug,
		Color: req.Color,
	}

	if err := db.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tag created successfully",
		"tag":     tag,
	})
}

// UpdateTag handles updating a tag (admin only)
func UpdateTag(c *gin.Context) {
	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var tag models.Tag
	if err := db.First(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Check if new name conflicts with existing tag
	if req.Name != "" && req.Name != tag.Name {
		var existingTag models.Tag
		if err := db.Where("name = ? AND id != ?", req.Name, tagID).First(&existingTag).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Tag name already exists"})
			return
		}
		tag.Name = req.Name
		tag.Slug = generateSlug(req.Name)
	}

	// Update color
	if req.Color != "" {
		tag.Color = req.Color
	}

	if err := db.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag updated successfully",
		"tag":     tag,
	})
}

// DeleteTag handles deleting a tag (admin only)
func DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var tag models.Tag
	if err := db.First(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Check if tag has associated posts or projects
	var postCount int64
	db.Model(&models.Post{}).Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ?", tagID).Count(&postCount)

	var projectCount int64
	db.Model(&models.Project{}).Joins("JOIN project_tags ON projects.id = project_tags.project_id").
		Where("project_tags.tag_id = ?", tagID).Count(&projectCount)

	if postCount > 0 || projectCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Cannot delete tag with associated posts or projects",
			"post_count":    postCount,
			"project_count": projectCount,
		})
		return
	}

	if err := db.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

// GetTagPosts handles getting posts by tag
func GetTagPosts(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var tag models.Tag
	if err := db.Where("slug = ?", slug).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var posts []models.Post
	if err := db.Preload("Author").Preload("Tags").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND posts.status = ?", tag.ID, models.PostStatusPublished).
		Order("posts.created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tag posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag":   tag,
		"posts": posts,
		"total": len(posts),
	})
}

// GetTagProjects handles getting projects by tag
func GetTagProjects(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var tag models.Tag
	if err := db.Where("slug = ?", slug).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var projects []models.Project
	if err := db.Preload("Technologies").Preload("Tags").
		Joins("JOIN project_tags ON projects.id = project_tags.project_id").
		Where("project_tags.tag_id = ? AND projects.status = ?", tag.ID, models.ProjectStatusPublished).
		Order("projects.created_at DESC").
		Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tag projects"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tag":      tag,
		"projects": projects,
		"total":    len(projects),
	})
}

// GetPopularTags handles getting popular tags
func GetPopularTags(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	limit := 10 // Default limit

	var tags []models.Tag
	if err := db.Select("tags.*, COUNT(post_tags.post_id) + COUNT(project_tags.project_id) as usage_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN project_tags ON tags.id = project_tags.tag_id").
		Group("tags.id").
		Order("usage_count DESC").
		Limit(limit).
		Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch popular tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags":  tags,
		"total": len(tags),
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
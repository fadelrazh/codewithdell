package handlers

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePostRequest represents post creation request
type CreatePostRequest struct {
	Title       string   `json:"title" binding:"required,min=3"`
	Content     string   `json:"content" binding:"required,min=10"`
	Excerpt     string   `json:"excerpt"`
	Slug        string   `json:"slug"`
	Status      string   `json:"status" binding:"required,oneof=draft published archived"`
	TagIDs      []string `json:"tag_ids"`
}

// UpdatePostRequest represents post update request
type UpdatePostRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Excerpt     string   `json:"excerpt"`
	Slug        string   `json:"slug"`
	Status      string   `json:"status" binding:"omitempty,oneof=draft published archived"`
	CategoryID  string   `json:"category_id"`
	TagIDs      []string `json:"tag_ids"`
	IsPublished bool     `json:"is_published"`
}

// GetPosts handles getting all posts with pagination and filters
func GetPosts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Simple query without any filters
	var posts []models.Post
	if err := db.Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": len(posts),
	})
}

// GetPostBySlug handles getting a single post by slug
func GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	db := c.MustGet("db").(*gorm.DB)

	var post models.Post
	if err := db.Preload("Author").Preload("Tags").
		Where("slug = ? AND status = ?", slug, models.PostStatusPublished).
		First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Increment view count
	db.Model(&post).UpdateColumn("view_count", post.ViewCount+1)

	c.JSON(http.StatusOK, post)
}

// CreatePost handles creating a new post (admin only)
func CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	// Generate slug if not provided
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	// Check if slug already exists
	var existingPost models.Post
	if err := db.Where("slug = ?", req.Slug).First(&existingPost).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Post with this slug already exists"})
		return
	}

	// Convert userID to uint
	authorID, _ := strconv.ParseUint(userID, 10, 32)
	
	// Convert status string to PostStatus
	var status models.PostStatus
	switch req.Status {
	case "draft":
		status = models.PostStatusDraft
	case "published":
		status = models.PostStatusPublished
	case "archived":
		status = models.PostStatusArchived
	default:
		status = models.PostStatusDraft
	}

	// Create post
	post := models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Excerpt:  req.Excerpt,
		Slug:     req.Slug,
		Status:   status,
		AuthorID: uint(authorID),
	}

	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Add tags if provided
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := db.Where("id IN ?", req.TagIDs).Find(&tags).Error; err == nil {
			db.Model(&post).Association("Tags").Append(tags)
		}
	}

	// Load relationships
	db.Preload("Author").Preload("Tags").First(&post, post.ID)

	c.JSON(http.StatusCreated, post)
}

// UpdatePost handles updating a post (admin only)
func UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Find post
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Excerpt != "" {
		updates["excerpt"] = req.Excerpt
	}
	if req.Slug != "" {
		updates["slug"] = req.Slug
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.CategoryID != "" {
		updates["category_id"] = req.CategoryID
	}
	updates["is_published"] = req.IsPublished

	if err := db.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Update tags if provided
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := db.Where("id IN ?", req.TagIDs).Find(&tags).Error; err == nil {
			db.Model(&post).Association("Tags").Replace(tags)
		}
	}

	// Load relationships
	db.Preload("Author").Preload("Tags").First(&post, post.ID)

	c.JSON(http.StatusOK, post)
}

// DeletePost handles deleting a post (admin only)
func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	// Find post
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Soft delete
	if err := db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// generateSlug generates a URL-friendly slug from title
func generateSlug(title string) string {
	// Simple slug generation - in production, use a proper slug library
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple dashes
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")
	// Remove leading/trailing dashes
	slug = strings.Trim(slug, "-")
	return slug
} 
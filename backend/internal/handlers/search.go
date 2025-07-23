package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchRequest represents search request parameters
type SearchRequest struct {
	Query     string   `json:"query" form:"q"`
	Type      string   `json:"type" form:"type"` // "posts", "projects", "all"
	Category  string   `json:"category" form:"category"`
	Tags      []string `json:"tags" form:"tags"`
	Author    string   `json:"author" form:"author"`
	Status    string   `json:"status" form:"status"`
	SortBy    string   `json:"sort_by" form:"sort_by"` // "relevance", "date", "views", "likes"
	SortOrder string   `json:"sort_order" form:"sort_order"` // "asc", "desc"
	Page      int      `json:"page" form:"page"`
	Limit     int      `json:"limit" form:"limit"`
}

// SearchResponse represents search response
type SearchResponse struct {
	Query   string      `json:"query"`
	Type    string      `json:"type"`
	Results interface{} `json:"results"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Pages   int         `json:"pages"`
}

// Search handles advanced search functionality
func Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "relevance"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	db := c.MustGet("db").(*gorm.DB)

	switch req.Type {
	case "posts":
		results, total, err := searchPosts(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
			return
		}
		c.JSON(http.StatusOK, createSearchResponse(req, results, total))
	case "projects":
		results, total, err := searchProjects(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search projects"})
			return
		}
		c.JSON(http.StatusOK, createSearchResponse(req, results, total))
	case "all", "":
		// Search both posts and projects
		postResults, postTotal, err := searchPosts(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
			return
		}

		projectResults, projectTotal, err := searchProjects(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search projects"})
			return
		}

		combinedResults := gin.H{
			"posts":    postResults,
			"projects": projectResults,
		}

		c.JSON(http.StatusOK, createSearchResponse(req, combinedResults, postTotal+projectTotal))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search type"})
	}
}

// searchPosts performs search on posts
func searchPosts(db *gorm.DB, req SearchRequest) ([]models.Post, int64, error) {
	query := db.Preload("Author").Preload("Tags").Preload("Categories")

	// Apply search query
	if req.Query != "" {
		searchQuery := "%" + strings.ToLower(req.Query) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ?",
			searchQuery, searchQuery, searchQuery)
	}

	// Apply filters
	if req.Category != "" {
		query = query.Joins("JOIN post_categories ON posts.id = post_categories.post_id").
			Joins("JOIN categories ON post_categories.category_id = categories.id").
			Where("categories.slug = ?", req.Category)
	}

	if len(req.Tags) > 0 {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Joins("JOIN tags ON post_tags.tag_id = tags.id").
			Where("tags.slug IN ?", req.Tags)
	}

	if req.Author != "" {
		query = query.Joins("JOIN users ON posts.author_id = users.id").
			Where("users.username = ?", req.Author)
	}

	if req.Status != "" {
		query = query.Where("posts.status = ?", req.Status)
	} else {
		// Default to published posts for public search
		query = query.Where("posts.status = ?", models.PostStatusPublished)
	}

	// Apply sorting
	switch req.SortBy {
	case "date":
		query = query.Order("posts.created_at " + req.SortOrder)
	case "views":
		query = query.Order("posts.view_count " + req.SortOrder)
	case "likes":
		query = query.Order("posts.like_count " + req.SortOrder)
	case "relevance":
		fallthrough
	default:
		// For relevance, we'll sort by a combination of factors
		if req.Query != "" {
			// If there's a search query, prioritize posts that match the query
			query = query.Order("CASE WHEN LOWER(title) LIKE ? THEN 1 ELSE 2 END", "%"+strings.ToLower(req.Query)+"%")
		}
		query = query.Order("posts.view_count DESC, posts.created_at DESC")
	}

	// Get total count
	var total int64
	query.Model(&models.Post{}).Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	var posts []models.Post
	err := query.Find(&posts).Error
	return posts, total, err
}

// searchProjects performs search on projects
func searchProjects(db *gorm.DB, req SearchRequest) ([]models.Project, int64, error) {
	query := db.Preload("Technologies").Preload("Tags").Preload("Categories")

	// Apply search query
	if req.Query != "" {
		searchQuery := "%" + strings.ToLower(req.Query) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(content) LIKE ?",
			searchQuery, searchQuery, searchQuery)
	}

	// Apply filters
	if req.Category != "" {
		query = query.Joins("JOIN project_categories ON projects.id = project_categories.project_id").
			Joins("JOIN categories ON project_categories.category_id = categories.id").
			Where("categories.slug = ?", req.Category)
	}

	if len(req.Tags) > 0 {
		query = query.Joins("JOIN project_tags ON projects.id = project_tags.project_id").
			Joins("JOIN tags ON project_tags.tag_id = tags.id").
			Where("tags.slug IN ?", req.Tags)
	}

	if req.Status != "" {
		query = query.Where("projects.status = ?", req.Status)
	} else {
		// Default to published projects for public search
		query = query.Where("projects.status = ?", models.ProjectStatusPublished)
	}

	// Apply sorting
	switch req.SortBy {
	case "date":
		query = query.Order("projects.created_at " + req.SortOrder)
	case "views":
		query = query.Order("projects.view_count " + req.SortOrder)
	case "likes":
		query = query.Order("projects.like_count " + req.SortOrder)
	case "relevance":
		fallthrough
	default:
		// For relevance, we'll sort by a combination of factors
		if req.Query != "" {
			// If there's a search query, prioritize projects that match the query
			query = query.Order("CASE WHEN LOWER(title) LIKE ? THEN 1 ELSE 2 END", "%"+strings.ToLower(req.Query)+"%")
		}
		query = query.Order("projects.view_count DESC, projects.created_at DESC")
	}

	// Get total count
	var total int64
	query.Model(&models.Project{}).Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	query = query.Offset(offset).Limit(req.Limit)

	var projects []models.Project
	err := query.Find(&projects).Error
	return projects, total, err
}

// createSearchResponse creates a standardized search response
func createSearchResponse(req SearchRequest, results interface{}, total int64) SearchResponse {
	pages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	if pages <= 0 {
		pages = 1
	}

	return SearchResponse{
		Query:   req.Query,
		Type:    req.Type,
		Results: results,
		Total:   total,
		Page:    req.Page,
		Limit:   req.Limit,
		Pages:   pages,
	}
}

// GetSearchSuggestions handles getting search suggestions
func GetSearchSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{"suggestions": []string{}})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	searchQuery := "%" + strings.ToLower(query) + "%"

	var suggestions []string

	// Get post titles
	var postTitles []string
	db.Model(&models.Post{}).
		Where("LOWER(title) LIKE ? AND status = ?", searchQuery, models.PostStatusPublished).
		Limit(5).
		Pluck("title", &postTitles)
	suggestions = append(suggestions, postTitles...)

	// Get project titles
	var projectTitles []string
	db.Model(&models.Project{}).
		Where("LOWER(title) LIKE ? AND status = ?", searchQuery, models.ProjectStatusPublished).
		Limit(5).
		Pluck("title", &projectTitles)
	suggestions = append(suggestions, projectTitles...)

	// Get tag names
	var tagNames []string
	db.Model(&models.Tag{}).
		Where("LOWER(name) LIKE ?", searchQuery).
		Limit(5).
		Pluck("name", &tagNames)
	suggestions = append(suggestions, tagNames...)

	// Get category names
	var categoryNames []string
	db.Model(&models.Category{}).
		Where("LOWER(name) LIKE ?", searchQuery).
		Limit(5).
		Pluck("name", &categoryNames)
	suggestions = append(suggestions, categoryNames...)

	// Remove duplicates and limit results
	uniqueSuggestions := make([]string, 0)
	seen := make(map[string]bool)
	for _, suggestion := range suggestions {
		if !seen[suggestion] && len(uniqueSuggestions) < 10 {
			uniqueSuggestions = append(uniqueSuggestions, suggestion)
			seen[suggestion] = true
		}
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": uniqueSuggestions})
}

// GetSearchStats handles getting search statistics
func GetSearchStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Count total posts
	var postCount int64
	db.Model(&models.Post{}).Where("status = ?", models.PostStatusPublished).Count(&postCount)

	// Count total projects
	var projectCount int64
	db.Model(&models.Project{}).Where("status = ?", models.ProjectStatusPublished).Count(&projectCount)

	// Count total tags
	var tagCount int64
	db.Model(&models.Tag{}).Count(&tagCount)

	// Count total categories
	var categoryCount int64
	db.Model(&models.Category{}).Count(&categoryCount)

	// Get popular tags
	var popularTags []models.Tag
	db.Select("tags.*, COUNT(post_tags.post_id) + COUNT(project_tags.project_id) as usage_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN project_tags ON tags.id = project_tags.tag_id").
		Group("tags.id").
		Order("usage_count DESC").
		Limit(5).
		Find(&popularTags)

	c.JSON(http.StatusOK, gin.H{
		"total_posts":     postCount,
		"total_projects":  projectCount,
		"total_tags":      tagCount,
		"total_categories": categoryCount,
		"popular_tags":    popularTags,
	})
} 
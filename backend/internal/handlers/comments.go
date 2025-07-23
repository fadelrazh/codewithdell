package handlers

import (
	"net/http"
	"strconv"
	"time"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Content   string `json:"content" binding:"required,min=1,max=1000"`
	PostID    *uint  `json:"post_id"`
	ProjectID *uint  `json:"project_id"`
	ParentID  *uint  `json:"parent_id"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CommentResponse represents comment response
type CommentResponse struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	Content   string    `json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      struct {
		ID        uint   `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		Avatar    string `json:"avatar"`
	} `json:"user"`
	Children []CommentResponse `json:"children,omitempty"`
}

// GetComments handles getting comments for a post or project
func GetComments(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	postID := c.Query("post_id")
	projectID := c.Query("project_id")
	
	if postID == "" && projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either post_id or project_id is required"})
		return
	}

	var comments []models.Comment
	query := db.Preload("User").Preload("Children.User").Where("parent_id IS NULL")
	
	if postID != "" {
		postIDUint, _ := strconv.ParseUint(postID, 10, 32)
		query = query.Where("post_id = ?", uint(postIDUint))
	}
	
	if projectID != "" {
		projectIDUint, _ := strconv.ParseUint(projectID, 10, 32)
		query = query.Where("project_id = ?", uint(projectIDUint))
	}
	
	// Only show approved comments for public
	query = query.Where("status = ?", models.CommentStatusApproved)
	
	if err := query.Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	// Convert to response format
	var responses []CommentResponse
	for _, comment := range comments {
		responses = append(responses, convertCommentToResponse(comment))
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": responses,
		"total":    len(responses),
	})
}

// CreateComment handles creating a new comment
func CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that either post_id or project_id is provided
	if req.PostID == nil && req.ProjectID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either post_id or project_id is required"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	// Check if parent comment exists if provided
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := db.First(&parentComment, *req.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
			return
		}
	}

	// Check if post/project exists
	if req.PostID != nil {
		var post models.Post
		if err := db.First(&post, *req.PostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
	}

	if req.ProjectID != nil {
		var project models.Project
		if err := db.First(&project, *req.ProjectID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
	}

	// Create comment
	comment := models.Comment{
		Content:   req.Content,
		UserID:    uint(userIDUint),
		PostID:    req.PostID,
		ProjectID: req.ProjectID,
		ParentID:  req.ParentID,
		Status:    models.CommentStatusPending, // Default to pending for moderation
	}

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load user data
	db.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully and pending approval",
		"comment": convertCommentToResponse(comment),
	})
}

// UpdateComment handles updating a comment
func UpdateComment(c *gin.Context) {
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	commentID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if user owns the comment or is admin
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)
	userRole := c.MustGet("user_role").(string)
	
	if comment.UserID != uint(userIDUint) && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this comment"})
		return
	}

	// Update comment
	comment.Content = req.Content
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// Load user data
	db.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": convertCommentToResponse(comment),
	})
}

// DeleteComment handles deleting a comment
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)
	userRole := c.MustGet("user_role").(string)

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if user owns the comment or is admin
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)
	if comment.UserID != uint(userIDUint) && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this comment"})
		return
	}

	// Delete comment and its children
	if err := db.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// ApproveComment handles approving a comment (admin only)
func ApproveComment(c *gin.Context) {
	commentID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	comment.Status = models.CommentStatusApproved
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment approved successfully"})
}

// RejectComment handles rejecting a comment (admin only)
func RejectComment(c *gin.Context) {
	commentID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	comment.Status = models.CommentStatusSpam
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment rejected successfully"})
}

// GetPendingComments handles getting pending comments for moderation (admin only)
func GetPendingComments(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var comments []models.Comment
	if err := db.Preload("User").Preload("Post").Preload("Project").
		Where("status = ?", models.CommentStatusPending).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending comments"})
		return
	}

	var responses []CommentResponse
	for _, comment := range comments {
		responses = append(responses, convertCommentToResponse(comment))
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": responses,
		"total":    len(responses),
	})
}

// convertCommentToResponse converts a comment model to response format
func convertCommentToResponse(comment models.Comment) CommentResponse {
	response := CommentResponse{
		ID:        comment.ID,
		UUID:      comment.UUID,
		Content:   comment.Content,
		Status:    string(comment.Status),
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User: struct {
			ID        uint   `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Avatar    string `json:"avatar"`
		}{
			ID:        comment.User.ID,
			FirstName: comment.User.FirstName,
			LastName:  comment.User.LastName,
			Username:  comment.User.Username,
			Avatar:    comment.User.Avatar,
		},
	}

	// Add children comments
	for _, child := range comment.Children {
		response.Children = append(response.Children, convertCommentToResponse(child))
	}

	return response
} 
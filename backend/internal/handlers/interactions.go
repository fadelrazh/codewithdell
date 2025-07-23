package handlers

import (
	"net/http"
	"strconv"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LikePost handles liking a post
func LikePost(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	postIDUint, _ := strconv.ParseUint(postID, 10, 32)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	// Check if post exists
	var post models.Post
	if err := db.First(&post, postIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user already liked the post
	var existingLike models.Like
	if err := db.Where("user_id = ? AND post_id = ?", userIDUint, postIDUint).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Post already liked"})
		return
	}

	// Create like
	like := models.Like{
		UserID: uint(userIDUint),
		PostID: &post.ID,
	}

	if err := db.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	// Increment like count
	db.Model(&post).UpdateColumn("like_count", post.LikeCount+1)

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

// UnlikePost handles unliking a post
func UnlikePost(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	postIDUint, _ := strconv.ParseUint(postID, 10, 32)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	// Check if like exists
	var like models.Like
	if err := db.Where("user_id = ? AND post_id = ?", userIDUint, postIDUint).First(&like).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	// Delete like
	if err := db.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}

	// Decrement like count
	var post models.Post
	db.First(&post, postIDUint)
	if post.LikeCount > 0 {
		db.Model(&post).UpdateColumn("like_count", post.LikeCount-1)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

// BookmarkPost handles bookmarking a post
func BookmarkPost(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	postIDUint, _ := strconv.ParseUint(postID, 10, 32)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	// Check if post exists
	var post models.Post
	if err := db.First(&post, postIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user already bookmarked the post
	var existingBookmark models.Bookmark
	if err := db.Where("user_id = ? AND post_id = ?", userIDUint, postIDUint).First(&existingBookmark).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Post already bookmarked"})
		return
	}

	// Create bookmark
	bookmark := models.Bookmark{
		UserID: uint(userIDUint),
		PostID: &post.ID,
	}

	if err := db.Create(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bookmark post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post bookmarked successfully"})
}

// RemoveBookmark handles removing a bookmark
func RemoveBookmark(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	postIDUint, _ := strconv.ParseUint(postID, 10, 32)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	// Check if bookmark exists
	var bookmark models.Bookmark
	if err := db.Where("user_id = ? AND post_id = ?", userIDUint, postIDUint).First(&bookmark).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bookmark not found"})
		return
	}

	// Delete bookmark
	if err := db.Delete(&bookmark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove bookmark"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed successfully"})
}

// GetUserLikes handles getting user's liked posts
func GetUserLikes(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	var likes []models.Like
	if err := db.Preload("Post").Preload("Post.Author").
		Where("user_id = ?", userIDUint).
		Order("created_at DESC").
		Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user likes"})
		return
	}

	var posts []models.Post
	for _, like := range likes {
		if like.Post != nil {
			posts = append(posts, *like.Post)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": len(posts),
	})
}

// GetUserBookmarks handles getting user's bookmarked posts
func GetUserBookmarks(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	var bookmarks []models.Bookmark
	if err := db.Preload("Post").Preload("Post.Author").
		Where("user_id = ?", userIDUint).
		Order("created_at DESC").
		Find(&bookmarks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user bookmarks"})
		return
	}

	var posts []models.Post
	for _, bookmark := range bookmarks {
		if bookmark.Post != nil {
			posts = append(posts, *bookmark.Post)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": len(posts),
	})
}

// CheckUserInteraction handles checking if user has interacted with a post
func CheckUserInteraction(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	postIDUint, _ := strconv.ParseUint(postID, 10, 32)
	userIDUint, _ := strconv.ParseUint(userID, 10, 32)

	// Check if user liked the post
	var like models.Like
	isLiked := db.Where("user_id = ? AND post_id = ?", userIDUint, postIDUint).First(&like).Error == nil

	// Check if user bookmarked the post
	var bookmark models.Bookmark
	isBookmarked := db.Where("user_id = ? AND post_id = ?", userIDUint, postIDUint).First(&bookmark).Error == nil

	c.JSON(http.StatusOK, gin.H{
		"is_liked":     isLiked,
		"is_bookmarked": isBookmarked,
	})
}

// GetPostStats handles getting post interaction statistics
func GetPostStats(c *gin.Context) {
	postID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)

	postIDUint, _ := strconv.ParseUint(postID, 10, 32)

	var post models.Post
	if err := db.First(&post, postIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Count likes
	var likeCount int64
	db.Model(&models.Like{}).Where("post_id = ?", postIDUint).Count(&likeCount)

	// Count bookmarks
	var bookmarkCount int64
	db.Model(&models.Bookmark{}).Where("post_id = ?", postIDUint).Count(&bookmarkCount)

	// Count comments
	var commentCount int64
	db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postIDUint, models.CommentStatusApproved).Count(&commentCount)

	c.JSON(http.StatusOK, gin.H{
		"post_id":        post.ID,
		"view_count":     post.ViewCount,
		"like_count":     likeCount,
		"bookmark_count": bookmarkCount,
		"comment_count":  commentCount,
	})
} 
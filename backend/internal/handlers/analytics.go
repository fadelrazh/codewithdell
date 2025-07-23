package handlers

import (
	"net/http"
	"time"

	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AnalyticsResponse represents analytics response
type AnalyticsResponse struct {
	Overview    map[string]interface{} `json:"overview"`
	Trends      map[string]interface{} `json:"trends"`
	Popular     map[string]interface{} `json:"popular"`
	UserStats   map[string]interface{} `json:"user_stats"`
	Engagement  map[string]interface{} `json:"engagement"`
}

// GetAnalytics handles getting comprehensive analytics
func GetAnalytics(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get overview statistics
	overview := getOverviewStats(db)
	
	// Get trends
	trends := getTrendsStats(db)
	
	// Get popular content
	popular := getPopularContent(db)
	
	// Get user statistics
	userStats := getUserStats(db)
	
	// Get engagement statistics
	engagement := getEngagementStats(db)

	response := AnalyticsResponse{
		Overview:   overview,
		Trends:     trends,
		Popular:    popular,
		UserStats:  userStats,
		Engagement: engagement,
	}

	c.JSON(http.StatusOK, response)
}

// getOverviewStats gets overview statistics
func getOverviewStats(db *gorm.DB) map[string]interface{} {
	// Total posts
	var totalPosts int64
	db.Model(&models.Post{}).Where("status = ?", models.PostStatusPublished).Count(&totalPosts)

	// Total projects
	var totalProjects int64
	db.Model(&models.Project{}).Where("status = ?", models.ProjectStatusPublished).Count(&totalProjects)

	// Total users
	var totalUsers int64
	db.Model(&models.User{}).Where("status = ?", models.StatusActive).Count(&totalUsers)

	// Total comments
	var totalComments int64
	db.Model(&models.Comment{}).Where("status = ?", models.CommentStatusApproved).Count(&totalComments)

	// Total likes
	var totalLikes int64
	db.Model(&models.Like{}).Count(&totalLikes)

	// Total bookmarks
	var totalBookmarks int64
	db.Model(&models.Bookmark{}).Count(&totalBookmarks)

	return map[string]interface{}{
		"total_posts":     totalPosts,
		"total_projects":  totalProjects,
		"total_users":     totalUsers,
		"total_comments":  totalComments,
		"total_likes":     totalLikes,
		"total_bookmarks": totalBookmarks,
	}
}

// getTrendsStats gets trends statistics
func getTrendsStats(db *gorm.DB) map[string]interface{} {
	// Posts created in last 30 days
	var recentPosts int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	db.Model(&models.Post{}).Where("created_at >= ?", thirtyDaysAgo).Count(&recentPosts)

	// Projects created in last 30 days
	var recentProjects int64
	db.Model(&models.Project{}).Where("created_at >= ?", thirtyDaysAgo).Count(&recentProjects)

	// New users in last 30 days
	var newUsers int64
	db.Model(&models.User{}).Where("created_at >= ?", thirtyDaysAgo).Count(&newUsers)

	// Comments in last 30 days
	var recentComments int64
	db.Model(&models.Comment{}).Where("created_at >= ?", thirtyDaysAgo).Count(&recentComments)

	// Monthly growth
	var lastMonthPosts int64
	sixtyDaysAgo := time.Now().AddDate(0, 0, -60)
	db.Model(&models.Post{}).Where("created_at >= ? AND created_at < ?", sixtyDaysAgo, thirtyDaysAgo).Count(&lastMonthPosts)

	postGrowth := float64(0)
	if lastMonthPosts > 0 {
		postGrowth = float64(recentPosts-lastMonthPosts) / float64(lastMonthPosts) * 100
	}

	return map[string]interface{}{
		"recent_posts":     recentPosts,
		"recent_projects":  recentProjects,
		"new_users":        newUsers,
		"recent_comments":  recentComments,
		"post_growth_rate": postGrowth,
	}
}

// getPopularContent gets popular content statistics
func getPopularContent(db *gorm.DB) map[string]interface{} {
	// Most viewed posts
	var popularPosts []models.Post
	db.Preload("Author").Preload("Tags").
		Where("status = ?", models.PostStatusPublished).
		Order("view_count DESC").
		Limit(5).
		Find(&popularPosts)

	// Most liked posts
	var mostLikedPosts []models.Post
	db.Preload("Author").Preload("Tags").
		Where("status = ?", models.PostStatusPublished).
		Order("like_count DESC").
		Limit(5).
		Find(&mostLikedPosts)

	// Most viewed projects
	var popularProjects []models.Project
	db.Preload("Technologies").Preload("Tags").
		Where("status = ?", models.ProjectStatusPublished).
		Order("view_count DESC").
		Limit(5).
		Find(&popularProjects)

	// Popular tags
	var popularTags []models.Tag
	db.Select("tags.*, COUNT(post_tags.post_id) + COUNT(project_tags.project_id) as usage_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN project_tags ON tags.id = project_tags.tag_id").
		Group("tags.id").
		Order("usage_count DESC").
		Limit(10).
		Find(&popularTags)

	return map[string]interface{}{
		"popular_posts":    popularPosts,
		"most_liked_posts": mostLikedPosts,
		"popular_projects": popularProjects,
		"popular_tags":     popularTags,
	}
}

// getUserStats gets user statistics
func getUserStats(db *gorm.DB) map[string]interface{} {
	// Active users (users with activity in last 30 days)
	var activeUsers int64
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	db.Model(&models.User{}).
		Joins("LEFT JOIN posts ON users.id = posts.author_id").
		Joins("LEFT JOIN comments ON users.id = comments.user_id").
		Joins("LEFT JOIN likes ON users.id = likes.user_id").
		Where("posts.created_at >= ? OR comments.created_at >= ? OR likes.created_at >= ?", 
			thirtyDaysAgo, thirtyDaysAgo, thirtyDaysAgo).
		Distinct("users.id").
		Count(&activeUsers)

	// Top contributors (users with most posts)
	var topContributors []models.User
	db.Select("users.*, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON users.id = posts.author_id").
		Where("posts.status = ?", models.PostStatusPublished).
		Group("users.id").
		Order("post_count DESC").
		Limit(5).
		Find(&topContributors)

	// Most engaged users (users with most comments)
	var mostEngagedUsers []models.User
	db.Select("users.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON users.id = comments.user_id").
		Where("comments.status = ?", models.CommentStatusApproved).
		Group("users.id").
		Order("comment_count DESC").
		Limit(5).
		Find(&mostEngagedUsers)

	return map[string]interface{}{
		"active_users":        activeUsers,
		"top_contributors":    topContributors,
		"most_engaged_users":  mostEngagedUsers,
	}
}

// getEngagementStats gets engagement statistics
func getEngagementStats(db *gorm.DB) map[string]interface{} {
	// Average views per post
	var avgViews float64
	db.Model(&models.Post{}).
		Where("status = ?", models.PostStatusPublished).
		Select("AVG(view_count)").
		Scan(&avgViews)

	// Average likes per post
	var avgLikes float64
	db.Model(&models.Post{}).
		Where("status = ?", models.PostStatusPublished).
		Select("AVG(like_count)").
		Scan(&avgLikes)

	// Average comments per post
	var avgComments float64
	db.Model(&models.Post{}).
		Where("status = ?", models.PostStatusPublished).
		Select("AVG(comment_count)").
		Scan(&avgComments)

	// Engagement rate (likes + comments / views)
	var totalViews int64
	var totalLikes int64
	var totalComments int64
	
	db.Model(&models.Post{}).
		Where("status = ?", models.PostStatusPublished).
		Select("SUM(view_count), SUM(like_count), SUM(comment_count)").
		Scan(&totalViews, &totalLikes, &totalComments)

	engagementRate := float64(0)
	if totalViews > 0 {
		engagementRate = float64(totalLikes+totalComments) / float64(totalViews) * 100
	}

	return map[string]interface{}{
		"avg_views_per_post":    avgViews,
		"avg_likes_per_post":    avgLikes,
		"avg_comments_per_post": avgComments,
		"engagement_rate":       engagementRate,
		"total_views":           totalViews,
		"total_likes":           totalLikes,
		"total_comments":        totalComments,
	}
}

// GetPostStats handles getting post statistics
func GetPostStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get post ID from URL parameter
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Count likes
	var likeCount int64
	db.Model(&models.Like{}).Where("post_id = ?", post.ID).Count(&likeCount)

	// Count bookmarks
	var bookmarkCount int64
	db.Model(&models.Bookmark{}).Where("post_id = ?", post.ID).Count(&bookmarkCount)

	// Count comments
	var commentCount int64
	db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", post.ID, models.CommentStatusApproved).Count(&commentCount)

	// Get recent activity
	var recentComments []models.Comment
	db.Preload("User").
		Where("post_id = ? AND status = ?", post.ID, models.CommentStatusApproved).
		Order("created_at DESC").
		Limit(5).
		Find(&recentComments)

	c.JSON(http.StatusOK, gin.H{
		"post_id":         post.ID,
		"title":           post.Title,
		"view_count":      post.ViewCount,
		"like_count":      likeCount,
		"bookmark_count":  bookmarkCount,
		"comment_count":   commentCount,
		"recent_comments": recentComments,
	})
}

// GetUserStats handles getting user statistics
func GetUserStats(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Count user's posts
	var postCount int64
	db.Model(&models.Post{}).Where("author_id = ? AND status = ?", user.ID, models.PostStatusPublished).Count(&postCount)

	// Count user's comments
	var commentCount int64
	db.Model(&models.Comment{}).Where("user_id = ? AND status = ?", user.ID, models.CommentStatusApproved).Count(&commentCount)

	// Count user's likes
	var likeCount int64
	db.Model(&models.Like{}).Where("user_id = ?", user.ID).Count(&likeCount)

	// Count user's bookmarks
	var bookmarkCount int64
	db.Model(&models.Bookmark{}).Where("user_id = ?", user.ID).Count(&bookmarkCount)

	// Get user's recent posts
	var recentPosts []models.Post
	db.Where("author_id = ? AND status = ?", user.ID, models.PostStatusPublished).
		Order("created_at DESC").
		Limit(5).
		Find(&recentPosts)

	// Get user's recent comments
	var recentComments []models.Comment
	db.Preload("Post").
		Where("user_id = ? AND status = ?", user.ID, models.CommentStatusApproved).
		Order("created_at DESC").
		Limit(5).
		Find(&recentComments)

	c.JSON(http.StatusOK, gin.H{
		"user_id":         user.ID,
		"username":        user.Username,
		"post_count":      postCount,
		"comment_count":   commentCount,
		"like_count":      likeCount,
		"bookmark_count":  bookmarkCount,
		"recent_posts":    recentPosts,
		"recent_comments": recentComments,
	})
} 
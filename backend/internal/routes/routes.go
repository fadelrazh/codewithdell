package routes

import (
	"codewithdell/backend/internal/config"
	"codewithdell/backend/internal/handlers"
	"codewithdell/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

// Setup configures all routes for the application
func Setup(router *gin.Engine, cfg *config.Config) {
	// Health check
	router.GET("/health", handlers.HealthCheck)
	
	// Static file serving for uploads
	router.Static("/uploads", "./uploads")
	
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("")
		{
			// Auth routes
			auth := public.Group("/auth")
			{
				auth.POST("/register", handlers.Register)
				auth.POST("/login", handlers.Login)
				auth.POST("/refresh", handlers.RefreshToken)
			}

			// Public content routes
			posts := public.Group("/posts")
			{
				posts.GET("", handlers.GetPosts)
				posts.GET("/:slug", handlers.GetPostBySlug)
			}

			// Categories routes
			categories := public.Group("/categories")
			{
				categories.GET("", handlers.GetCategories)
				categories.GET("/:slug", handlers.GetCategoryBySlug)
				categories.GET("/:slug/posts", handlers.GetCategoryPosts)
				categories.GET("/:slug/projects", handlers.GetCategoryProjects)
			}

			// Tags routes
			tags := public.Group("/tags")
			{
				tags.GET("", handlers.GetTags)
				tags.GET("/popular", handlers.GetPopularTags)
				tags.GET("/:slug", handlers.GetTagBySlug)
				tags.GET("/:slug/posts", handlers.GetTagPosts)
				tags.GET("/:slug/projects", handlers.GetTagProjects)
			}

			// Comments routes (public read)
			comments := public.Group("/comments")
			{
				comments.GET("", handlers.GetComments)
			}

			// Search routes
			search := public.Group("/search")
			{
				search.GET("", handlers.Search)
				search.GET("/suggestions", handlers.GetSearchSuggestions)
				search.GET("/stats", handlers.GetSearchStats)
			}

			// Analytics routes (public read)
			analytics := public.Group("/analytics")
			{
				analytics.GET("", handlers.GetAnalytics)
				analytics.GET("/posts/:id", handlers.GetPostStats)
				analytics.GET("/users/:id", handlers.GetUserStats)
			}

			// Simple test endpoint
			public.GET("/test", handlers.TestEndpoint)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.Auth(cfg.JWT.Secret))
		{
			// User profile
			profile := protected.Group("/profile")
			{
				profile.GET("", handlers.GetProfile)
				profile.PUT("", handlers.UpdateProfile)
			}

			// User interactions
			interactions := protected.Group("/interactions")
			{
				interactions.POST("/posts/:id/like", handlers.LikePost)
				interactions.DELETE("/posts/:id/like", handlers.UnlikePost)
				interactions.POST("/posts/:id/bookmark", handlers.BookmarkPost)
				interactions.DELETE("/posts/:id/bookmark", handlers.RemoveBookmark)
				interactions.GET("/posts/:id/check", handlers.CheckUserInteraction)
				interactions.GET("/likes", handlers.GetUserLikes)
				interactions.GET("/bookmarks", handlers.GetUserBookmarks)
			}

			// Comments routes (authenticated write)
			comments := protected.Group("/comments")
			{
				comments.POST("", handlers.CreateComment)
				comments.PUT("/:id", handlers.UpdateComment)
				comments.DELETE("/:id", handlers.DeleteComment)
			}

			// Upload routes
			upload := protected.Group("/upload")
			{
				upload.POST("/image", handlers.UploadImage)
				upload.POST("/file", handlers.UploadFile)
				upload.DELETE("/:filename", handlers.DeleteFile)
				upload.GET("/stats", handlers.GetUploadStats)
			}
		}

		// Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth(cfg.JWT.Secret), middleware.RequireRole("admin"))
		{
			// Content management
			posts := admin.Group("/posts")
			{
				posts.POST("", handlers.CreatePost)
				posts.PUT("/:id", handlers.UpdatePost)
				posts.DELETE("/:id", handlers.DeletePost)
			}

			// Categories management
			categories := admin.Group("/categories")
			{
				categories.POST("", handlers.CreateCategory)
				categories.PUT("/:id", handlers.UpdateCategory)
				categories.DELETE("/:id", handlers.DeleteCategory)
			}

			// Tags management
			tags := admin.Group("/tags")
			{
				tags.POST("", handlers.CreateTag)
				tags.PUT("/:id", handlers.UpdateTag)
				tags.DELETE("/:id", handlers.DeleteTag)
			}

			// Comments moderation
			comments := admin.Group("/comments")
			{
				comments.GET("/pending", handlers.GetPendingComments)
				comments.POST("/:id/approve", handlers.ApproveComment)
				comments.POST("/:id/reject", handlers.RejectComment)
			}
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus metrics
	router.GET("/metrics", middleware.PrometheusHandler())
} 
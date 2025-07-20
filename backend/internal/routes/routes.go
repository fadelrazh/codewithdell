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
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus metrics
	router.GET("/metrics", middleware.PrometheusHandler())
} 
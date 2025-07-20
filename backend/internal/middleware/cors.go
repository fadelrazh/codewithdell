package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware for handling cross-origin requests
func CORS(origin string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{origin}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"X-Requested-With",
		"X-CSRF-Token",
		"X-API-Key",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"X-Total-Count",
		"X-Page-Count",
		"X-Current-Page",
		"X-Per-Page",
	}
	
	return cors.New(config)
} 
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "codewithdell-backend",
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

// TestEndpoint handles test requests
func TestEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Backend is working!",
		"timestamp": time.Now(),
		"service": "codewithdell-backend",
	})
} 
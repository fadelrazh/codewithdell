package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	Requests int           // Number of requests allowed
	Window   time.Duration // Time window for the limit
}

// RateLimiter represents a rate limiter instance
type RateLimiter struct {
	redisClient *redis.Client
	config      RateLimitConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		config:      config,
	}
}

// RateLimit creates a rate limiting middleware
func RateLimit(redisClient *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(redisClient, config)
	return limiter.Handle
}

// Handle processes the rate limiting logic
func (rl *RateLimiter) Handle(c *gin.Context) {
	// Get client identifier (IP address or user ID)
	clientID := rl.getClientID(c)
	
	// Check if request is allowed
	allowed, remaining, resetTime, err := rl.isAllowed(c.Request.Context(), clientID)
	if err != nil {
		// If Redis is unavailable, allow the request
		c.Next()
		return
	}

	// Set rate limit headers
	c.Header("X-RateLimit-Limit", strconv.Itoa(rl.config.Requests))
	c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate limit exceeded",
			"retry_after": resetTime - time.Now().Unix(),
		})
		c.Abort()
		return
	}

	c.Next()
}

// isAllowed checks if the request is allowed based on rate limiting rules
func (rl *RateLimiter) isAllowed(ctx context.Context, clientID string) (bool, int, int64, error) {
	key := "rate_limit:" + clientID
	now := time.Now().Unix()
	windowStart := now - int64(rl.config.Window.Seconds())

	// Use Redis pipeline for atomic operations
	pipe := rl.redisClient.Pipeline()
	
	// Remove old entries
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
	
	// Count current requests
	countCmd := pipe.ZCard(ctx, key)
	
	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: now,
	})
	
	// Set expiration
	pipe.Expire(ctx, key, rl.config.Window)
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return true, rl.config.Requests, now + int64(rl.config.Window.Seconds()), err
	}

	currentCount := countCmd.Val()
	allowed := currentCount < int64(rl.config.Requests)
	remaining := rl.config.Requests - int(currentCount)
	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, now + int64(rl.config.Window.Seconds()), nil
}

// getClientID returns a unique identifier for the client
func (rl *RateLimiter) getClientID(c *gin.Context) string {
	// Try to get user ID from context first
	if userID, exists := c.Get("user_id"); exists {
		return "user:" + userID.(string)
	}
	
	// Fall back to IP address
	return "ip:" + c.ClientIP()
}

// PerEndpointRateLimit creates rate limiting middleware for specific endpoints
func PerEndpointRateLimit(redisClient *redis.Client, configs map[string]RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		config, exists := configs[endpoint]
		if !exists {
			// Use default config
			config = RateLimitConfig{
				Requests: 100,
				Window:   time.Minute,
			}
		}

		limiter := NewRateLimiter(redisClient, config)
		limiter.Handle(c)
	}
}

// UserRateLimit creates rate limiting middleware for authenticated users
func UserRateLimit(redisClient *redis.Client, config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply rate limiting to authenticated users
		if userID, exists := c.Get("user_id"); exists {
			limiter := NewRateLimiter(redisClient, config)
			limiter.Handle(c)
		} else {
			// For unauthenticated users, use IP-based rate limiting with stricter limits
			ipConfig := RateLimitConfig{
				Requests: config.Requests / 2, // Half the requests for unauthenticated users
				Window:   config.Window,
			}
			limiter := NewRateLimiter(redisClient, ipConfig)
			limiter.Handle(c)
		}
	}
}

// BurstRateLimit creates rate limiting middleware that allows bursts
func BurstRateLimit(redisClient *redis.Client, config RateLimitConfig, burstSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := "ip:" + c.ClientIP()
		key := "burst_rate_limit:" + clientID
		now := time.Now().Unix()

		// Check current burst count
		burstCount, err := redisClient.Get(c.Request.Context(), key).Int()
		if err != nil && err != redis.Nil {
			c.Next()
			return
		}

		if burstCount >= burstSize {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Burst rate limit exceeded",
			})
			c.Abort()
			return
		}

		// Increment burst count
		pipe := redisClient.Pipeline()
		pipe.Incr(c.Request.Context(), key)
		pipe.Expire(c.Request.Context(), key, config.Window)
		pipe.Exec(c.Request.Context())

		c.Next()
	}
} 
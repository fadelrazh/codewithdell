package middleware

import (
	"codewithdell/backend/internal/logger"

	"github.com/gin-gonic/gin"
)

// Logger middleware for HTTP request logging
func Logger(log *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.HTTPRequest(
			param.Method,
			param.Path,
			param.ClientIP,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
		)
		return ""
	})
} 
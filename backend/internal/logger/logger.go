package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger wraps zerolog logger
type Logger struct {
	logger zerolog.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	// Set up pretty console logging for development
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Use pretty console logging in development
	if os.Getenv("ENVIRONMENT") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	logger := log.With().Str("service", "codewithdell-backend").Logger()

	return &Logger{
		logger: logger,
	}
}

// GetLevel returns the current log level
func (l *Logger) GetLevel() zerolog.Level {
	return l.logger.GetLevel()
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level zerolog.Level) {
	l.logger = l.logger.Level(level)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.logger.Debug()
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}
	event.Msg(msg)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.logger.Info()
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.logger.Warn()
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Error().Err(err)
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Fatal().Err(err)
	for _, field := range fields {
		for key, value := range field {
			event = event.Interface(key, value)
		}
	}
	event.Msg(msg)
}

// WithContext returns a logger with context fields
func (l *Logger) WithContext(fields map[string]interface{}) *Logger {
	event := l.logger.With()
	for key, value := range fields {
		event = event.Interface(key, value)
	}
	
	return &Logger{
		logger: event.Logger(),
	}
}

// HTTPRequest logs HTTP request information
func (l *Logger) HTTPRequest(method, path, remoteAddr string, statusCode int, duration time.Duration, userAgent string) {
	l.Info("HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
		"status_code": statusCode,
		"duration":    duration.String(),
		"user_agent":  userAgent,
	})
}

// DatabaseQuery logs database query information
func (l *Logger) DatabaseQuery(sql string, duration time.Duration, rows int64) {
	l.Debug("Database Query", map[string]interface{}{
		"sql":      sql,
		"duration": duration.String(),
		"rows":     rows,
	})
}

// CacheHit logs cache hit information
func (l *Logger) CacheHit(key string) {
	l.Debug("Cache Hit", map[string]interface{}{
		"key": key,
	})
}

// CacheMiss logs cache miss information
func (l *Logger) CacheMiss(key string) {
	l.Debug("Cache Miss", map[string]interface{}{
		"key": key,
	})
}

// Authentication logs authentication events
func (l *Logger) Authentication(userID uint, action string, success bool, ip string) {
	level := l.logger.Info()
	if !success {
		level = l.logger.Warn()
	}
	
	level.Interface("user_id", userID).
		Str("action", action).
		Bool("success", success).
		Str("ip", ip).
		Msg("Authentication Event")
}

// Authorization logs authorization events
func (l *Logger) Authorization(userID uint, resource string, action string, allowed bool, ip string) {
	level := l.logger.Info()
	if !allowed {
		level = l.logger.Warn()
	}
	
	level.Interface("user_id", userID).
		Str("resource", resource).
		Str("action", action).
		Bool("allowed", allowed).
		Str("ip", ip).
		Msg("Authorization Event")
} 
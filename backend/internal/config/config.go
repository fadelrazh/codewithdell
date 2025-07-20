package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	Storage  StorageConfig
}

// AppConfig holds application configuration
type AppConfig struct {
	Environment string
	Port        string
	CORSOrigin  string
	LogLevel    string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxOpen  int
	MaxIdle  int
	Timeout  time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Provider string
	APIKey   string
	From     string
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	Provider string
	Bucket   string
	Region   string
	AccessKey string
	SecretKey string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnv("PORT", "8080"),
			CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "codewithdell"),
			User:     getEnv("DB_USER", "codewithdell"),
			Password: getEnv("DB_PASSWORD", "codewithdell123"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxOpen:  getEnvAsInt("DB_MAX_OPEN", 25),
			MaxIdle:  getEnvAsInt("DB_MAX_IDLE", 5),
			Timeout:  getEnvAsDuration("DB_TIMEOUT", 5*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			Expiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		Email: EmailConfig{
			Provider: getEnv("EMAIL_PROVIDER", "sendgrid"),
			APIKey:   getEnv("EMAIL_API_KEY", ""),
			From:     getEnv("EMAIL_FROM", "noreply@codewithdell.com"),
		},
		Storage: StorageConfig{
			Provider:  getEnv("STORAGE_PROVIDER", "local"),
			Bucket:    getEnv("STORAGE_BUCKET", "codewithdell"),
			Region:    getEnv("STORAGE_REGION", "us-east-1"),
			AccessKey: getEnv("STORAGE_ACCESS_KEY", ""),
			SecretKey: getEnv("STORAGE_SECRET_KEY", ""),
		},
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	log.Info().Msg("Configuration loaded successfully")
	return config, nil
}

// validate validates the configuration
func (c *Config) validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("port is required")
	}

	if c.Database.Host == "" || c.Database.Name == "" || c.Database.User == "" {
		return fmt.Errorf("database configuration is incomplete")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
} 
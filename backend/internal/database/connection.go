package database

import (
	"fmt"
	"time"

	"codewithdell/backend/internal/config"
	"codewithdell/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.Timeout)

	DB = db
	return db, nil
}

// CloseConnection closes the database connection
func CloseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// RunMigrations runs database migrations
func RunMigrations(cfg config.DatabaseConfig) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Auto migrate all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Project{},
		&models.Category{},
		&models.Tag{},
		&models.Technology{},
		&models.Comment{},
		&models.Like{},
		&models.Bookmark{},
		&models.Screenshot{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
} 
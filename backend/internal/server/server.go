package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codewithdell/backend/internal/config"
	"codewithdell/backend/internal/database"
	"codewithdell/backend/internal/logger"
	"codewithdell/backend/internal/middleware"
	"codewithdell/backend/internal/redis"
	"codewithdell/backend/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
	server *http.Server
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	// Set Gin mode
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	return &Server{
		router: router,
		config: cfg,
	}
}

// Initialize sets up the server with all middleware and routes
func (s *Server) Initialize() error {
	// Initialize logger
	logger := logger.NewLogger()

	// Initialize database
	_, err := database.NewConnection(s.config.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run migrations
	if err := database.RunMigrations(s.config.Database); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Add database to context using the global DB instance
	s.router.Use(func(c *gin.Context) {
		c.Set("db", database.GetDB())
		c.Next()
	})

	// Initialize Redis
	redisClient, err := redis.NewClient(s.config.Redis)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}
	defer redisClient.Close()

	// Add middleware
	s.router.Use(middleware.Logger(logger))
	s.router.Use(middleware.CORS(s.config.App.CORSOrigin))
	s.router.Use(middleware.Security())
	s.router.Use(middleware.Prometheus())

	// Setup routes
	routes.Setup(s.router, s.config)

	return nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create HTTP server
	s.server = &http.Server{
		Addr:         ":" + s.config.App.Port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Msgf("Starting server on port %s", s.config.App.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
	return nil
} 
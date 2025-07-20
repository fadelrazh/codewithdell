package main

import (
	"codewithdell/backend/internal/config"
	"codewithdell/backend/internal/server"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// @title CodeWithDell API
// @version 1.0
// @description Advanced blog and showcase platform API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Create and initialize server
	srv := server.New(cfg)
	if err := srv.Initialize(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize server")
	}

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
} 
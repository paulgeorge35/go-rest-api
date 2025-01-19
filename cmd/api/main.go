package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"rest-api/internal/config"
	"rest-api/internal/database"
	"rest-api/internal/server"
	"rest-api/pkg/logger"

	_ "rest-api/docs" // for swagger

	"github.com/joho/godotenv"
)

// @title REST API
// @version 1.0
// @description A secure REST API with authentication and OAuth support
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger := logger.NewLogger()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", err)
	}

	// Initialize and start server
	srv := server.New(cfg, logger, db)
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}
}

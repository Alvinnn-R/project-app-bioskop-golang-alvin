package main

import (
	"log"
	"project-app-bioskop/cmd"
	"project-app-bioskop/internal/data/repository"
	"project-app-bioskop/internal/wire"
	"project-app-bioskop/pkg/database"
	"project-app-bioskop/pkg/utils"
)

func main() {
	// Load configuration from .env
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}

	// Initialize logger
	logger, err := utils.InitLogger(config.PathLogging, config.Debug)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database connection
	db, err := database.InitDB(config.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := repository.NewRepository(db)

	// Wire all dependencies and routes
	route := wire.Wiring(repo, config, logger)

	// Start HTTP server
	cmd.APiserver(route)
}

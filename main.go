package main

import (
	"log"

	"github.com/HersheyPlus/go-rbac/config"
	"github.com/HersheyPlus/go-rbac/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.InitializeDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Successfully connected to " + db.Name())
}

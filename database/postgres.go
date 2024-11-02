package database

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/HersheyPlus/go-rbac/config"
)

var DB *gorm.DB

func InitializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	var err error
	
	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	dsn := cfg.Database.GetDatabaseURL()
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB to set connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Run migrations
	if err := MigrateDatabase(DB); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Println("Database connection and migrations completed successfully")
	return DB, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
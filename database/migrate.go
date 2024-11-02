package database

import (
	"fmt"
	"log"
	"github.com/HersheyPlus/go-rbac/models"
	"github.com/HersheyPlus/go-rbac/pkg/password"
	"gorm.io/gorm"
)

// MigrateDatabase handles database migrations and seeding
func MigrateDatabase(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// Auto-migrate all models
	if err := autoMigrateModels(db); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	// Seed initial data
	if err := seedInitialData(db); err != nil {
		return fmt.Errorf("failed to seed initial data: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// autoMigrateModels performs automatic migration for all models
func autoMigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
	)
}

// seedInitialData seeds the database with initial required data
func seedInitialData(db *gorm.DB) error {
	// Create default permissions
	permissions := []models.Permission{
		{Name: "user:create", Description: "Can create users"},
		{Name: "user:read", Description: "Can read users"},
		{Name: "user:update", Description: "Can update users"},
		{Name: "user:delete", Description: "Can delete users"},
		{Name: "role:create", Description: "Can create roles"},
		{Name: "role:read", Description: "Can read roles"},
		{Name: "role:update", Description: "Can update roles"},
		{Name: "role:delete", Description: "Can delete roles"},
		{Name: "permission:read", Description: "Can read permissions"},
	}

	// Create default roles
	adminRole := models.Role{
		Name:        "admin",
		Description: "System administrator with full access",
	}

	userRole := models.Role{
		Name:        "user",
		Description: "Regular user with limited access",
	}

	// Create default admin user with a password that meets all requirements
	// Password contains: uppercase, lowercase, numbers, and special character
	defaultAdminPassword := "Admin@123!"
	hashedPassword, err := password.HashPassword(defaultAdminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	adminUser := models.User{
		Email:     "admin@example.com",
		Password:  hashedPassword,
		FirstName: "Admin",
		LastName:  "User",
		Active:    true,
	}

	// Begin transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Create permissions if they don't exist
		for _, perm := range permissions {
			var existingPerm models.Permission
			if err := tx.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&perm).Error; err != nil {
						return fmt.Errorf("failed to create permission %s: %w", perm.Name, err)
					}
				} else {
					return fmt.Errorf("failed to check permission %s: %w", perm.Name, err)
				}
			}
		}

		// Create admin role if it doesn't exist
		var existingAdminRole models.Role
		if err := tx.Where("name = ?", adminRole.Name).First(&existingAdminRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&adminRole).Error; err != nil {
					return fmt.Errorf("failed to create admin role: %w", err)
				}
				existingAdminRole = adminRole
			} else {
				return fmt.Errorf("failed to check admin role: %w", err)
			}
		}

		// Create user role if it doesn't exist
		var existingUserRole models.Role
		if err := tx.Where("name = ?", userRole.Name).First(&existingUserRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&userRole).Error; err != nil {
					return fmt.Errorf("failed to create user role: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check user role: %w", err)
			}
		}

		// Assign all permissions to admin role
		var allPermissions []models.Permission
		if err := tx.Find(&allPermissions).Error; err != nil {
			return fmt.Errorf("failed to fetch permissions: %w", err)
		}

		if err := tx.Model(&existingAdminRole).Association("Permissions").Replace(allPermissions); err != nil {
			return fmt.Errorf("failed to assign permissions to admin role: %w", err)
		}

		// Create admin user if it doesn't exist
		var existingAdminUser models.User
		if err := tx.Where("email = ?", adminUser.Email).First(&existingAdminUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&adminUser).Error; err != nil {
					return fmt.Errorf("failed to create admin user: %w", err)
				}
				existingAdminUser = adminUser
			} else {
				return fmt.Errorf("failed to check admin user: %w", err)
			}
		}

		// Assign admin role to admin user
		if err := tx.Model(&existingAdminUser).Association("Roles").Replace(&existingAdminRole); err != nil {
			return fmt.Errorf("failed to assign admin role to admin user: %w", err)
		}

		log.Printf("Successfully created admin user with email: %s", adminUser.Email)
		log.Println("Please remember to change the default admin password after first login")

		return nil
	})
}
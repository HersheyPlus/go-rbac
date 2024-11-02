package models

import (
	"time"
	"github.com/google/uuid"
)

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    uuid.UUID `gorm:"type:uuid;primary_key" json:"user_id"`
	RoleID    uuid.UUID `gorm:"type:uuid;primary_key" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}
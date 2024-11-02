package models

import (
	"time"
	"github.com/google/uuid"
)

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primary_key" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid;primary_key" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}
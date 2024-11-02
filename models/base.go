package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


// BeforeCreate will set Base fields before creating
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate will set UpdatedAt before updating
func (base *Base) BeforeUpdate(tx *gorm.DB) error {
	base.UpdatedAt = time.Now()
	return nil
}
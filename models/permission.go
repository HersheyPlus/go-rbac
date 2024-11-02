package models


type Permission struct {
	Base
	Name        string `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description string `gorm:"type:varchar(200)" json:"description"`
	// Relationships
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles"`
}
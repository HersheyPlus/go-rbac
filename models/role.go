package models

type Role struct {
	Base
	Name        string `gorm:"type:varchar(50);unique;not null" json:"name"`
	Description string `gorm:"type:varchar(200)" json:"description"`
	// Relationships
	Users       []User       `gorm:"many2many:user_roles;" json:"users"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
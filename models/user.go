package models

type User struct {
	Base
	Email     string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string `gorm:"type:varchar(255);not null" json:"-"`
	FirstName string `gorm:"type:varchar(50)" json:"first_name"`
	LastName  string `gorm:"type:varchar(50)" json:"last_name"`
	Active    bool   `gorm:"default:true" json:"active"`
	// Relationships
	Roles []Role `gorm:"many2many:user_roles;" json:"roles"`
}

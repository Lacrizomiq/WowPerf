package models

import (
	"gorm.io/gorm"
)

// User is the struct for the user model
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username" validate:"required,min=3,max=50"`
	Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"-" validate:"required,strongpassword"`
}

// UserCreate is the struct for creating a new user
type UserCreate struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,containsany=!@#$%^&*()_+"`
}

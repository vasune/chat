package models

import "gorm.io/gorm"

type User struct {
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	gorm.Model
}

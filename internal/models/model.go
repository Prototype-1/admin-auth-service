package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uint       `gorm:"primaryKey;autoIncrement;column:user_id"`
	FirstName      string     `gorm:"size:100"`
	LastName       string     `gorm:"size:100"`
	Email          string     `gorm:"size:255;unique;not null"`
	Password       string     `gorm:"size:255;not null"`
	Phone          string     `gorm:"size:20"`
	Role           string     `gorm:"size:50;not null"` 
	BlockedStatus  bool       `gorm:"default:false"`
	InactiveStatus bool       `gorm:"default:false"`
	SuspendedAt    *time.Time `gorm:"default:null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Admin struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement;column:admin_id"`
	Email     string    `gorm:"size:255;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	Role      string    `gorm:"size:50;not null;default:'admin'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthToken struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement;column:token_id"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:500;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type RefreshToken struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement;column:refresh_token_id"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:500;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
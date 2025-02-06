package models

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `gorm:"size:255;not null"`
	Email         string    `gorm:"uniqueIndex;size:255;not null"`
	BlockedStatus bool      `gorm:"default:false"`
	InactiveStatus bool     `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

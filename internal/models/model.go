package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Email     string         `gorm:"unique;not null"`
	Password  string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
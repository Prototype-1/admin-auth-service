package models

import (
	"time"
)

type Route struct {
    RouteID    int       `gorm:"primaryKey;autoIncrement" json:"route_id"`
    RouteName  string    `gorm:"type:varchar(255);not null" json:"route_name"`
    StartStopID int      `gorm:"not null" json:"start_stop_id"`
    EndStopID   int      `gorm:"not null" json:"end_stop_id"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    CategoryID int      `gorm:"not null" json:"category_id"`  
}
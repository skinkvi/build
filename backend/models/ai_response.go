package models

import "time"

type AIResponse struct {
	ID           uint `gorm:"primaryKey"`
	ProjectID    uint `gorm:"not null"`
	ResponseData string
	ImageURL     string
	CreatedAt    time.Time
}

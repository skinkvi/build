package models

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	CreatedAt    time.Time
	Provider     string `gorm:"index"`
	ProviderID   string `gorm:"uniqueIndex"`
	Projects     []Project
}

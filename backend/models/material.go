package models

import "time"

type Material struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"not null"`
	Category         string
	ImageURL         string
	Description      string
	PurchaseLocation string
	CreatedAt        time.Time
	ProjectMaterials []ProjectMaterial
}

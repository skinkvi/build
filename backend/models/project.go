package models

import "time"

type Project struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"not null"`
	ProjectType      string `gorm:"not null;check:project_type IN ('design', 'fix', 'identify')"`
	RequestData      string
	ImageURL         string
	CreatedAt        time.Time
	AIResponses      []AIResponse
	ProjectMaterials []ProjectMaterial
}

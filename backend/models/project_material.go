package models

type ProjectMaterial struct {
	ProjectID  uint `gorm:"primaryKey"`
	MaterialID uint `gorm:"primaryKey"`
}

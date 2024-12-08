package models

import "time"

// Material представляет материал, необходимый для ремонтов.
type Material struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	Category         string    `gorm:"type:varchar(100)" json:"category"`
	ImageURL         string    `gorm:"type:text" json:"image_url"`
	Description      string    `gorm:"type:text" json:"description"`
	PurchaseLocation string    `gorm:"type:text" json:"purchase_location"`
	CreatedAt        time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at"`
}

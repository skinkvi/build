package storage

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

type Material struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	ImageURL         string    `json:"image_url"`
	Description      string    `json:"description"`
	PurchaseLocation string    `json:"purchase_location"`
	CreatedAt        time.Time `json:"created_at"`
}

type Project struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProjectType string    `json:"project_type"`
	RequestData string    `json:"request_data"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectMaterial struct {
	ProjectID  int `json:"project_id"`
	MaterialID int `json:"material_id"`
}

type AIResponse struct {
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	ResponseData string    `json:"response_data"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
}

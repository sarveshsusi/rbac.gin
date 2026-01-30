package models

import "github.com/google/uuid"

type Model struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name    string    `json:"name" gorm:"not null"`
	BrandID uuid.UUID `json:"brand_id" gorm:"type:uuid;index"`
}

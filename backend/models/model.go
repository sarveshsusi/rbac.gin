package models

import "github.com/google/uuid"

type Model struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name    string    `gorm:"not null"`
	BrandID uuid.UUID `gorm:"type:uuid;not null"`
}

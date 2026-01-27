package models

import "github.com/google/uuid"

type Brand struct {
	ID   uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string    `json:"name" gorm:"unique;not null"`
}

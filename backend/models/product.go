package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	BrandID    uuid.UUID `json:"brand_id" gorm:"type:uuid;not null"`
	ModelID    uuid.UUID `json:"model_id" gorm:"type:uuid;not null"`
	CreatedBy  uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
}
type CreateProductRequest struct {
	Name       string    `json:"name" binding:"required"`
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
	BrandID    uuid.UUID `json:"brand_id" binding:"required"`
	ModelID    uuid.UUID `json:"model_id" binding:"required"`
}

package models

import (
	"time"
	"github.com/google/uuid"
)

type CustomerProduct struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID uuid.UUID `gorm:"index"`
	ProductID  uuid.UUID `gorm:"index"`
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
}

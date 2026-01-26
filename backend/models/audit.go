package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Entity     string
	EntityID   uuid.UUID
	Action     string
	PerformedBy uuid.UUID
	IP         string
	UserAgent  string
	CreatedAt  time.Time
}

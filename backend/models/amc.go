package models

import (
	"time"

	"github.com/google/uuid"
)

type AMCContract struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID uuid.UUID
	StartDate time.Time
	EndDate   time.Time
	Value     float64
	Status    string // active, expired, renewed
	CreatedAt time.Time
}

type AMCSchedule struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	AMCID     uuid.UUID
	VisitDate time.Time
	Completed bool
	TicketID  *uuid.UUID
}

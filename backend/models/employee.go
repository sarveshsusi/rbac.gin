package models

import (
	"time"

	"github.com/google/uuid"
)

type SupportEngineer struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;uniqueIndex"` // RoleSupport
	Designation string    `gorm:"type:varchar(100)"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type ServiceVisit struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID
	EngineerID uuid.UUID
	StartTime  time.Time
	EndTime    *time.Time
	Notes      string
}

type GPSLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EngineerID uuid.UUID
	Latitude   float64
	Longitude  float64
	LoggedAt   time.Time
}

type DigitalSignature struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID  uuid.UUID
	SignedBy  string
	FileURL   string
	SignedAt  time.Time
}

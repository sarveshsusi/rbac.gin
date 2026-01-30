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

func (SupportEngineer) TableName() string {
	return "support_engineers"
}

type ServiceVisit struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid"`
	EngineerID uuid.UUID `gorm:"type:uuid"`
	StartTime  time.Time
	EndTime    *time.Time
	Notes      string
}

func (ServiceVisit) TableName() string {
	return "service_visits"
}

type GPSLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	EngineerID uuid.UUID `gorm:"type:uuid"`
	Latitude   float64
	Longitude  float64
	LoggedAt   time.Time
}

func (GPSLog) TableName() string {
	return "gps_logs"
}

type DigitalSignature struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID uuid.UUID `gorm:"type:uuid"`
	SignedBy string
	FileURL  string
	SignedAt time.Time
}

func (DigitalSignature) TableName() string {
	return "digital_signatures"
}

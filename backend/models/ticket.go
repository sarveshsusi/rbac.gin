package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string


const (
	TicketCustomerCreated  TicketStatus = "customer_created"
	TicketAdminReviewed    TicketStatus = "admin_reviewed"
	TicketAssignedSupport  TicketStatus = "assigned_to_support"
	TicketResolvedSupport  TicketStatus = "resolved_by_support"
	TicketClosedByAdmin    TicketStatus = "closed_by_admin"
	TicketFeedbackGiven    TicketStatus = "feedback_given"
)

type TicketPriority string

const (
	PriorityLow      TicketPriority = "low"
	PriorityMedium   TicketPriority = "medium"
	PriorityHigh     TicketPriority = "high"
	PriorityCritical TicketPriority = "critical"
)

// models/ticket.go

type Ticket struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID  uuid.UUID       `gorm:"type:uuid;index"`
	ProductID   uuid.UUID       `gorm:"type:uuid;index"`
	AMCId       uuid.UUID       `gorm:"type:uuid;index"`
	Title       string          `gorm:"type:varchar(255);not null"`
	Description string          `gorm:"type:text"`
	Status      TicketStatus    `gorm:"type:varchar(50);index"`
	Priority    TicketPriority  `gorm:"type:varchar(30);index"`
	SLAHours    int
	TargetAt    *time.Time
	ClosedAt    *time.Time
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


type TicketAssignment struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID        uuid.UUID `gorm:"type:uuid;index"`
	EngineerID      uuid.UUID `gorm:"type:uuid;index"`
	AssignedBy      uuid.UUID `gorm:"type:uuid"`
	AssignedAt      time.Time
}
type TicketStatusHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID  uuid.UUID `gorm:"type:uuid;index"`
	OldStatus string
	NewStatus string
	ChangedBy uuid.UUID
	ChangedAt time.Time
}
type TicketComment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid;index"`
	UserID     uuid.UUID `gorm:"type:uuid"`
	Comment    string    `gorm:"type:text"`
	IsInternal bool
	CreatedAt  time.Time
}

type TicketAttachment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid;index"`
	FileURL    string
	FileType   string
	UploadedBy uuid.UUID
	CreatedAt  time.Time
}


type TicketFeedback struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID  uuid.UUID
	EngineerID uuid.UUID
	Rating    int
	Comment   string
	CreatedAt time.Time
}


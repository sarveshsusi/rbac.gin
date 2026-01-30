package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	StatusOpen       TicketStatus = "Open"
	StatusAssigned   TicketStatus = "Assigned"
	StatusInProgress TicketStatus = "In Progress"
	StatusClosed     TicketStatus = "Closed"
)

type TicketPriority string

const (
	PriorityLow      TicketPriority = "Low"
	PriorityStandard TicketPriority = "Standard"
	PriorityCritical TicketPriority = "Critical"
)

type SupportMode string

const (
	SupportModeOnSite SupportMode = "On-site"
	SupportModeRemote SupportMode = "Remote"
	SupportModePhone  SupportMode = "Phone"
)

type ServiceCallType string

const (
	ServiceTypeWarranty ServiceCallType = "Warranty"
	ServiceTypeService  ServiceCallType = "Service"
	ServiceTypeAMC      ServiceCallType = "AMC"
)

// models/ticket.go

type Ticket struct {
	ID                uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID        uuid.UUID       `gorm:"type:uuid;index"`
	ProductID         uuid.UUID       `gorm:"type:uuid;index"`
	AMCId             uuid.UUID       `gorm:"type:uuid;index"`
	Title             string          `gorm:"type:varchar(255);not null"`
	Description       string          `gorm:"type:text"`
	Status            TicketStatus    `gorm:"type:varchar(50);index"`
	Priority          TicketPriority  `gorm:"type:varchar(30);index"`
	SupportMode       SupportMode     `gorm:"type:varchar(50)"`
	ServiceCallType   ServiceCallType `gorm:"type:varchar(50)"`
	ClosureProofImage string          `gorm:"type:text"`
	SLAHours          int
	TargetAt          *time.Time
	ClosedAt          *time.Time
	CreatedBy         uuid.UUID `gorm:"type:uuid"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TicketAssignment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid;index"`
	EngineerID uuid.UUID `gorm:"type:uuid;index"`
	AssignedBy uuid.UUID `gorm:"type:uuid"`
	AssignedAt time.Time
}

func (TicketAssignment) TableName() string {
	return "ticket_assignments"
}

type TicketStatusHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID  uuid.UUID `gorm:"type:uuid;index"`
	OldStatus string
	NewStatus string
	ChangedBy uuid.UUID `gorm:"type:uuid"`
	ChangedAt time.Time
}

func (TicketStatusHistory) TableName() string {
	return "ticket_status_histories"
}

type TicketComment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid;index"`
	UserID     uuid.UUID `gorm:"type:uuid"`
	Comment    string    `gorm:"type:text"`
	IsInternal bool
	CreatedAt  time.Time
}

func (TicketComment) TableName() string {
	return "ticket_comments"
}

type TicketAttachment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid;index"`
	FileURL    string
	FileType   string
	UploadedBy uuid.UUID `gorm:"type:uuid"`
	CreatedAt  time.Time
}

func (TicketAttachment) TableName() string {
	return "ticket_attachments"
}

type TicketFeedback struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TicketID   uuid.UUID `gorm:"type:uuid"`
	EngineerID uuid.UUID `gorm:"type:uuid"`
	Rating     int
	Comment    string
	CreatedAt  time.Time
}

func (TicketFeedback) TableName() string {
	return "ticket_feedbacks"
}

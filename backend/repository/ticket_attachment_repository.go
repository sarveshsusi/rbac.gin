package repository

import (
	"rbac/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketAttachmentRepository struct {
	db *gorm.DB
}

func NewTicketAttachmentRepository(db *gorm.DB) *TicketAttachmentRepository {
	return &TicketAttachmentRepository{db: db}
}

func (r *TicketAttachmentRepository) Create(
	ticketID uuid.UUID,
	url string,
	fileType string,
	uploadedBy uuid.UUID,
) error {

	return r.db.Create(&models.TicketAttachment{
		TicketID:   ticketID,
		FileURL:    url,
		FileType:   fileType,
		UploadedBy: uploadedBy,
	}).Error
}

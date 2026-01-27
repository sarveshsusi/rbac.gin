// repository/ticket_repository.go
package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"rbac/models"
)

type TicketRepository struct {
	db *gorm.DB
}

func (r *TicketRepository) DB() *gorm.DB {
	return r.db
} 
func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

/*
=====================
 Transaction Helper
=====================
*/
func (r *TicketRepository) WithTx(
	fn func(tx *gorm.DB) error,
) error {
	return r.db.Transaction(fn)
} 

func (r *TicketRepository) WithTransaction(
	fn func(txRepo *TicketRepository) error,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &TicketRepository{db: tx}
		return fn(txRepo)
	})
}

/*
=====================
 Queries
=====================
*/
func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) FindByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) FindByEngineer(
	engineerID uuid.UUID,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	err := r.db.
		Joins("JOIN ticket_assignments ta ON ta.ticket_id = tickets.id").
		Where("ta.engineer_id = ?", engineerID).
		Order("tickets.created_at DESC").
		Find(&tickets).Error

	return tickets, err
}

func (r *TicketRepository) FindByCustomer(
	customerID uuid.UUID,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	err := r.db.
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&tickets).Error

	return tickets, err
}

/*
=====================
 Assignment
=====================
*/
func (r *TicketRepository) CreateAssignment(
	ticketID uuid.UUID,
	engineerID uuid.UUID,
	adminID uuid.UUID,
) error {

	return r.db.Create(&models.TicketAssignment{
		TicketID:   ticketID,
		EngineerID: engineerID,
		AssignedBy: adminID,
	}).Error
}



func (r *TicketRepository) UpdateStatusNoTx(
	ticketID uuid.UUID,
	newStatus models.TicketStatus,
	changedBy uuid.UUID,
) error {
	return r.UpdateStatus(r.db, ticketID, newStatus, changedBy)
}

/*
=====================
 Status + History
=====================
*/
func (r *TicketRepository) UpdateStatus(
	tx *gorm.DB,
	ticketID uuid.UUID,
	newStatus models.TicketStatus,
	changedBy uuid.UUID,
) error {

	

	var ticket models.Ticket
	if err := tx.First(&ticket, "id = ?", ticketID).Error; err != nil {
		return err
	}

	if err := tx.Model(&models.Ticket{}).
		Where("id = ?", ticketID).
		Update("status", newStatus).Error; err != nil {
		return err
	}

	return tx.Create(&models.TicketStatusHistory{
		TicketID:  ticketID,
		OldStatus: string(ticket.Status),
		NewStatus: string(newStatus),
		ChangedBy: changedBy,
	}).Error
}

func (r *TicketRepository) CreateStatusHistory(
	ticketID uuid.UUID,
	oldStatus string,
	newStatus string,
	changedBy uuid.UUID,
) error {

	return r.db.Create(&models.TicketStatusHistory{
		TicketID:  ticketID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		ChangedBy: changedBy,
	}).Error
}

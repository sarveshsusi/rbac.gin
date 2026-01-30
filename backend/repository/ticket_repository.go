// repository/ticket_repository.go
package repository

import (
	"rbac/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *TicketRepository) GetByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) GetAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
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
) error {
	return r.UpdateStatus(ticketID, newStatus)
}

/*
=====================
 Status + History
=====================
*/
/*
=====================
 Updates
=====================
*/
func (r *TicketRepository) UpdateFields(
	ticketID uuid.UUID,
	updates map[string]interface{},
) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", ticketID).Updates(updates).Error
}

func (r *TicketRepository) AssignEmployee(
	ticket *models.Ticket,
	assignment *models.TicketAssignment,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Update Ticket (Priority, Status, etc)
		if err := tx.Save(ticket).Error; err != nil {
			return err
		}

		// 2. Create Assignment Record
		if err := tx.Create(assignment).Error; err != nil {
			return err
		}

		// 3. Create History Log
		if err := tx.Create(&models.TicketStatusHistory{
			TicketID:  ticket.ID,
			OldStatus: string(models.StatusOpen), // Assuming from Open
			NewStatus: string(ticket.Status),
			ChangedBy: assignment.AssignedBy,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *TicketRepository) UpdateStatus(
	ticketID uuid.UUID,
	newStatus models.TicketStatus,
) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", ticketID).Update("status", newStatus).Error
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

func (r *TicketRepository) FindOverdueTickets(
	days int,
) ([]models.Ticket, error) {

	var tickets []models.Ticket

	cutoff := time.Now().AddDate(0, 0, -days)

	err := r.db.
		Where(
			"status NOT IN ? AND created_at < ?",
			[]models.TicketStatus{
				models.StatusClosed,
			},
			cutoff,
		).
		Find(&tickets).Error

	return tickets, err
}

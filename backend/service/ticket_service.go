package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"rbac/models"
	"rbac/repository"
)

type TicketService struct {
	repo *repository.TicketRepository
}

func NewTicketService(repo *repository.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

/*
	=========================
	  CUSTOMER: CREATE TICKET
	=========================
*/
// CreateCustomerTicket - Customer only provides title and description
func (s *TicketService) CreateCustomerTicket(
	customerID uuid.UUID,
	title, description string,
) (*models.Ticket, error) {

	// Customer creates ticket with minimal info
	// Admin will assign product, AMC, priority, support mode, etc. later
	ticket := &models.Ticket{
		ID:          uuid.New(),
		CustomerID:  customerID,
		Title:       title,
		Description: description,
		Status:      models.StatusOpen,
		CreatedBy:   customerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// CreateTicket - Legacy method, kept for backward compatibility if needed
func (s *TicketService) CreateTicket(
	customerID, productID, amcID uuid.UUID,
	title, description string,
) (*models.Ticket, error) {

	// 1. Create ticket
	ticket := &models.Ticket{
		ID:          uuid.New(),
		CustomerID:  customerID,
		ProductID:   productID,
		AMCId:       amcID,
		Title:       title,
		Description: description,
		Status:      models.StatusOpen, // Default Status
	}

	if err := s.repo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

/*
	=========================
	  ADMIN: CREATE TICKET (ON BEHALF)

=========================
*/
func (s *TicketService) AdminCreateTicket(
	ticket *models.Ticket,
) (*models.Ticket, error) {
	// Admin sets everything upfront
	ticket.Status = models.StatusOpen

	if err := s.repo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

/*
	=========================
	  ADMIN: ASSIGN TICKET

=========================
*/
func (s *TicketService) AssignTicket(
	ticketID, engineerID, adminID uuid.UUID,
	priority models.TicketPriority,
	supportMode models.SupportMode,
	serviceType models.ServiceCallType,
) error {

	// 1. Get Ticket
	ticket, err := s.repo.GetByID(ticketID)
	if err != nil {
		return err
	}

	// 2. Update Fields
	ticket.Priority = priority
	ticket.SupportMode = supportMode
	ticket.ServiceCallType = serviceType
	ticket.Status = models.StatusAssigned // Update Status

	// 3. Create Assignment Record
	assignment := &models.TicketAssignment{
		TicketID:   ticketID,
		EngineerID: engineerID,
		AssignedBy: adminID,
		AssignedAt: time.Now(),
	}

	// 4. Save Updates (Transaction ideally)
	// For simplicity, calling repository method that handles update + assignment
	return s.repo.AssignEmployee(ticket, assignment)
}

/*
	=========================
	  SUPPORT: START TICKET

=========================
*/
func (s *TicketService) StartTicket(ticketID uuid.UUID) error {
	return s.repo.UpdateStatus(ticketID, models.StatusInProgress)
}

/*
	=========================
	  SUPPORT: CLOSE TICKET

=========================
*/
func (s *TicketService) CloseTicket(
	ticketID uuid.UUID,
	proofImageURL string,
) error {

	if proofImageURL == "" {
		return errors.New("proof image is mandatory to close ticket")
	}

	// Update Ticket: Status Closed, ClosedAt, ProofImage
	updates := map[string]interface{}{
		"status":              models.StatusClosed,
		"closed_at":           time.Now(),
		"closure_proof_image": proofImageURL,
	}

	return s.repo.UpdateFields(ticketID, updates)
}

/*
	=========================
	  GET TICKETS

=========================
*/
func (s *TicketService) GetAll() ([]models.Ticket, error) {
	return s.repo.GetAll()
}

package service

import (
	"time"

	"github.com/google/uuid"

	"rbac/models"
	"rbac/repository"
)

type TicketService struct {
	repo *repository.TicketRepository
}

func NewTicketService(r *repository.TicketRepository) *TicketService {
	return &TicketService{repo: r}
}



/* =====================
   CREATE TICKET (CUSTOMER)
===================== */
func (s *TicketService) CreateTicket(
	customerID uuid.UUID,
	title string,
	description string,
	priority models.TicketPriority,
	productID uuid.UUID,
	amc models.AMCContract,
) (*models.Ticket, error) {

	target := time.Now().Add(time.Duration(amc.SLAHours) * time.Hour)

	ticket := &models.Ticket{
		CustomerID:  customerID,
		ProductID:   productID,
		AMCId:       amc.ID,
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      models.TicketCustomerCreated,
		SLAHours:    amc.SLAHours,
		TargetAt:    &target,
		CreatedBy:   customerID,
	}

	if err := s.repo.Create(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

/* =====================
   ADMIN ASSIGN
===================== */
func (s *TicketService) AssignTicket(
	ticketID uuid.UUID,
	engineerID uuid.UUID,
	adminID uuid.UUID,
) error {

	return s.repo.WithTransaction(func(txRepo *repository.TicketRepository) error {
		if err := txRepo.CreateAssignment(
			ticketID,
			engineerID,
			adminID,
		); err != nil {
			return err
		}

		return txRepo.UpdateStatus(
			txRepo.DB(),
			ticketID,
			models.TicketAssignedSupport,
			adminID,
		)
	})
}

/* =====================
   SUPPORT RESOLVE
===================== */
func (s *TicketService) ResolveTicket(
	ticketID uuid.UUID,
	engineerID uuid.UUID,
) error {

	return s.repo.UpdateStatus(
		s.repo.DB(),
		ticketID,
		models.TicketResolvedSupport,
		engineerID,
	)
}

/* =====================
   ADMIN CLOSE
===================== */
func (s *TicketService) CloseTicket(
	ticketID uuid.UUID,
	adminID uuid.UUID,
) error {

	return s.repo.UpdateStatus(
		s.repo.DB(),
		ticketID,
		models.TicketClosedByAdmin,
		adminID,
	)
}

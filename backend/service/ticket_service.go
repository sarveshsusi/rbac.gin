package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"rbac/domain"
	"rbac/models"
	"rbac/repository"
)

type TicketService struct {
	repo           *repository.TicketRepository
	attachmentRepo *repository.TicketAttachmentRepository
	escalationRepo *repository.TicketEscalationRepository
}

func NewTicketService(
	repo *repository.TicketRepository,
	attachmentRepo *repository.TicketAttachmentRepository,
	escalationRepo *repository.TicketEscalationRepository,
) *TicketService {
	return &TicketService{
		repo:           repo,
		attachmentRepo: attachmentRepo,
		escalationRepo: escalationRepo,
	}
}

/* =====================
   CREATE TICKET (CUSTOMER)
===================== */
func (s *TicketService) CreateTicket(
	customerID uuid.UUID,
	title string,
	description string,
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
		Priority:    models.PriorityLow, // enforced
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
   INTERNAL: STATE CHANGE
===================== */
func (s *TicketService) changeStatus(
	ticketID uuid.UUID,
	newStatus models.TicketStatus,
	by uuid.UUID,
) error {

	ticket, err := s.repo.FindByID(ticketID)
	if err != nil {
		return err
	}

	if !domain.CanTransition(ticket.Status, newStatus) {
		return errors.New("invalid ticket status transition")
	}

	return s.repo.UpdateStatus(
		s.repo.DB(),
		ticketID,
		newStatus,
		by,
	)
}

/* =====================
   ADMIN REVIEW
===================== */
func (s *TicketService) AdminReviewTicket(
	ticketID uuid.UUID,
	adminID uuid.UUID,
) error {
	return s.changeStatus(
		ticketID,
		models.TicketAdminReviewed,
		adminID,
	)
}

/* =====================
   ADMIN ASSIGN
===================== */
func (s *TicketService) AssignTicket(
	ticketID uuid.UUID,
	engineerID uuid.UUID,
	adminID uuid.UUID,
	productID uuid.UUID,
	priority models.TicketPriority,
) error {

	return s.repo.WithTransaction(func(tx *repository.TicketRepository) error {

		if err := tx.CreateAssignment(ticketID, engineerID, adminID); err != nil {
			return err
		}

		if err := tx.DB().
			Model(&models.Ticket{}).
			Where("id = ?", ticketID).
			Updates(map[string]interface{}{
				"priority":   priority,
				"product_id": productID,
			}).Error; err != nil {
			return err
		}

		return s.changeStatus(
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
	return s.changeStatus(
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

	if err := s.changeStatus(
		ticketID,
		models.TicketClosedByAdmin,
		adminID,
	); err != nil {
		return err
	}

	// clear escalation record if exists
	return s.escalationRepo.ResolveByTicket(ticketID)
}

/* =====================
   ATTACHMENT (IMAGEKIT URL)
===================== */
func (s *TicketService) AddAttachment(
	ticketID uuid.UUID,
	url string,
	fileType string,
	userID uuid.UUID,
) error {

	return s.attachmentRepo.Create(
		ticketID,
		url,
		fileType,
		userID,
	)
}

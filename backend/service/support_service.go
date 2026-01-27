// service/support_service.go
package service

import (
	"github.com/google/uuid"
	"rbac/models"
	"rbac/repository"
)

type SupportService struct {
	ticketRepo *repository.TicketRepository
}

func NewSupportService(t *repository.TicketRepository) *SupportService {
	return &SupportService{ticketRepo: t}
}

func (s *SupportService) GetAssignedTickets(
	engineerID uuid.UUID,
) ([]models.Ticket, error) {
	return s.ticketRepo.FindByEngineer(engineerID)
}

// service/customer_service.go
package service

import (
	"github.com/google/uuid"
	"rbac/models"
	"rbac/repository"
)

type CustomerService struct {
	ticketRepo *repository.TicketRepository
}

func NewCustomerService(t *repository.TicketRepository) *CustomerService {
	return &CustomerService{ticketRepo: t}
}

func (s *CustomerService) GetCustomerTickets(
	customerID uuid.UUID,
) ([]models.Ticket, error) {
	return s.ticketRepo.FindByCustomer(customerID)
}

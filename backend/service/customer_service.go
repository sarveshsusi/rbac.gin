package service

import (
	"rbac/models"
	"rbac/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerService struct {
	db           *gorm.DB
	authRepo     *repository.AuthRepository
	customerRepo *repository.CustomerRepository
	ticketRepo   *repository.TicketRepository
}

func NewCustomerService(
	db *gorm.DB,
	authRepo *repository.AuthRepository,
	customerRepo *repository.CustomerRepository,
	ticketRepo *repository.TicketRepository,
) *CustomerService {
	return &CustomerService{
		db:           db,
		authRepo:     authRepo,
		customerRepo: customerRepo,
		ticketRepo:   ticketRepo,
	}
}
func (s *CustomerService) CreateCustomer(
	name string,
	email string,
	password string,
	company string,
	phone string,
	address string,
) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}

		user := &models.User{
			Name:     name,
			Email:    email,
			Password: string(hash),
			Role:     models.RoleCustomer,
			IsActive: true,
		}

		// 🔥 USE TX VERSION
		if err := s.authRepo.CreateUserTx(tx, user); err != nil {
			return err
		}

		customer := &models.Customer{
			UserID:   user.ID,
			Company:  company,
			Phone:    phone,
			Address:  address,
			IsActive: true,
		}

		return s.customerRepo.Create(tx, customer)
	})
}
func (s *CustomerService) GetAllCustomers(
	page int,
) ([]models.Customer, int64, error) {

	limit := 3
	if page <= 0 {
		page = 1
	}

	return s.customerRepo.GetAllPaginated(page, limit)
}
func (s *CustomerService) GetCustomerTickets(
	customerID uuid.UUID,
) ([]models.Ticket, error) {
	return s.ticketRepo.FindByCustomer(customerID)
}

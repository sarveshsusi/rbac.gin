package service

import (
	"errors"

	"github.com/google/uuid"

	"rbac/repository"
)

type CustomerProductService struct {
	repo *repository.CustomerProductRepository
}

func NewCustomerProductService(
	repo *repository.CustomerProductRepository,
) *CustomerProductService {
	return &CustomerProductService{repo: repo}
}

// func (s *CustomerProductService) AssignProductToCustomer(
// 	customerID, productID uuid.UUID,
// ) error {

// 	exists, err := s.repo.Exists(customerID, productID)
// 	if err != nil {
// 		return err
// 	}

// 	if exists {
// 		return errors.New("product already assigned to customer")
// 	}

// 	return s.repo.Assign(customerID, productID)
// }

func (s *CustomerProductService) AssignProductToCustomer(
	userID uuid.UUID,
	productID uuid.UUID,
) error {

	// 🔥 STEP 1: find customer by user_id
	customer, err := s.repo.GetCustomerByUserID(userID)
	if err != nil {
		return errors.New("customer profile not found")
	}

	// 🔥 STEP 2: assign product using customer.id
	return s.repo.Assign(customer.ID, productID)
}

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

	customer, err := s.repo.GetCustomerByUserID(userID)
	if err != nil {
		return errors.New("customer profile not found")
	}

	return s.repo.Assign(customer.ID, productID)
}

type CustomerProductResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
}

func (s *CustomerProductService) GetCustomerProductsByUserID(userID uuid.UUID) ([]CustomerProductResponse, error) {
	// Resolve Customer ID from User ID
	customer, err := s.repo.GetCustomerByUserID(userID)
	if err != nil {
		return nil, errors.New("customer profile not found")
	}

	products, err := s.repo.GetByCustomerID(customer.ID)
	if err != nil {
		return nil, err
	}

	var response []CustomerProductResponse
	for _, p := range products {
		response = append(response, CustomerProductResponse{
			ID:          p.ID,
			ProductID:   p.ProductID,
			ProductName: p.Product.Name,
		})
	}
	return response, nil
}

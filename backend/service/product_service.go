package service

import (
	"errors"

	"github.com/google/uuid"

	"rbac/models"
	"rbac/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(
	req *models.CreateProductRequest,
	adminID uuid.UUID,
) (*models.Product, error) {

	ok, err := s.repo.IsBrandInCategory(
		req.BrandID.String(),
		req.CategoryID.String(),
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("brand not allowed for category")
	}

	ok, err = s.repo.IsModelInBrand(
		req.ModelID.String(),
		req.BrandID.String(),
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("model not allowed for brand")
	}

	product := &models.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		BrandID:    req.BrandID,
		ModelID:    req.ModelID,
		CreatedBy:  adminID,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := s.repo.GetAll(&products); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) AssignProductToCustomer(
	customerID, productID uuid.UUID,
) error {
	return s.repo.AssignToCustomer(customerID, productID)
}

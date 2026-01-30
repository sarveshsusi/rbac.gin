package service

import (
	"rbac/models"
	"rbac/repository"

	"github.com/google/uuid"
)

type BrandService struct {
	repo *repository.BrandRepository
}

func NewBrandService(repo *repository.BrandRepository) *BrandService {
	return &BrandService{repo: repo}
}

func (s *BrandService) GetByCategory(categoryID uuid.UUID) (interface{}, error) {
	return s.repo.GetByCategory(categoryID)
}

func (s *BrandService) GetAll() ([]models.Brand, error) {
	return s.repo.GetAll()
}

func (s *BrandService) Create(
	name string,
	categoryID uuid.UUID,
) (*models.Brand, error) {

	// create brand
	brand, err := s.repo.Create(name)
	if err != nil {
		return nil, err
	}

	// 🔥 CRITICAL: map brand to category
	if err := s.repo.AssignToCategory(brand.ID, categoryID); err != nil {
		return nil, err
	}

	return brand, nil
}

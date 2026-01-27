package service

import (
	"github.com/google/uuid"
	"rbac/repository"
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

func (s *BrandService) Create(name string) (interface{}, error) {
	return s.repo.Create(name)
}

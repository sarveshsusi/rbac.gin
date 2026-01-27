package service

import (
	"github.com/google/uuid"
	"rbac/repository"
)

type ModelService struct {
	repo *repository.ModelRepository
}

func NewModelService(repo *repository.ModelRepository) *ModelService {
	return &ModelService{repo: repo}
}

func (s *ModelService) GetByBrand(brandID uuid.UUID) (interface{}, error) {
	return s.repo.GetByBrand(brandID)
}

func (s *ModelService) Create(name string, brandID uuid.UUID) (interface{}, error) {
	return s.repo.Create(name, brandID)
}

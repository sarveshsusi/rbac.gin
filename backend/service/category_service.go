package service

import "rbac/repository"

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() (interface{}, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(name string) (interface{}, error) {
	return s.repo.Create(name)
}

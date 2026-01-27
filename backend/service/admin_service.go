// service/admin_service.go
package service

import "rbac/repository"

type AdminService struct {
	repo *repository.DashboardRepository
}

func NewAdminService(r *repository.DashboardRepository) *AdminService {
	return &AdminService{repo: r}
}

func (s *AdminService) GetDashboardStats() (map[string]int64, error) {
	return s.repo.FetchAdminStats()
}

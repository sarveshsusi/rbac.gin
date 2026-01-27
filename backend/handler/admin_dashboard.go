// handler/admin_dashboard.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rbac/service"
)

type AdminDashboardHandler struct {
	service *service.AdminService
}

func NewAdminDashboardHandler(s *service.AdminService) *AdminDashboardHandler {
	return &AdminDashboardHandler{service: s}
}

func (h *AdminDashboardHandler) Dashboard(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

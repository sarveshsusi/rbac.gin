// handler/support_dashboard.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rbac/service"
)

type SupportDashboardHandler struct {
	service *service.SupportService
}

func NewSupportDashboardHandler(s *service.SupportService) *SupportDashboardHandler {
	return &SupportDashboardHandler{service: s}
}

func (h *SupportDashboardHandler) MyTickets(c *gin.Context) {
	engineerID := c.MustGet("user_id").(uuid.UUID)

	tickets, err := h.service.GetAssignedTickets(engineerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

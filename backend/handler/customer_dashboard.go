// handler/customer_dashboard.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rbac/service"
)

type CustomerDashboardHandler struct {
	service *service.CustomerService
}

func NewCustomerDashboardHandler(s *service.CustomerService) *CustomerDashboardHandler {
	return &CustomerDashboardHandler{service: s}
}

func (h *CustomerDashboardHandler) MyTickets(c *gin.Context) {
	customerID := c.MustGet("user_id").(uuid.UUID)

	tickets, err := h.service.GetCustomerTickets(customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

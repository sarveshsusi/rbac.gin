package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/models"
	"rbac/service"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(s *service.TicketService) *TicketHandler {
	return &TicketHandler{service: s}
}

/* =====================
   CREATE TICKET
===================== */

type CreateTicketRequest struct {
	Title       string                `json:"title" binding:"required"`
	Description string                `json:"description"`
	Priority    models.TicketPriority `json:"priority" binding:"required"`
	ProductID   uuid.UUID             `json:"product_id" binding:"required"`
	AMCId       uuid.UUID             `json:"amc_id" binding:"required"`
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	customerID := c.MustGet("user_id").(uuid.UUID)

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// AMC should be fetched from service/db — simplified here
	amc := models.AMCContract{
		ID:       req.AMCId,
		SLAHours: 24,
	}

	ticket, err := h.service.CreateTicket(
		customerID,
		req.Title,
		req.Description,
		req.Priority,
		req.ProductID,
		amc,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

/* =====================
   ADMIN ASSIGN
===================== */

type AssignTicketRequest struct {
	EngineerID uuid.UUID `json:"engineer_id" binding:"required"`
}

func (h *TicketHandler) AssignTicket(c *gin.Context) {
	ticketID, _ := uuid.Parse(c.Param("id"))
	adminID := c.MustGet("user_id").(uuid.UUID)

	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.AssignTicket(ticketID, req.EngineerID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket assigned"})
}

/* =====================
   SUPPORT RESOLVE
===================== */

func (h *TicketHandler) ResolveTicket(c *gin.Context) {
	ticketID, _ := uuid.Parse(c.Param("id"))
	engineerID := c.MustGet("user_id").(uuid.UUID)

	if err := h.service.ResolveTicket(ticketID, engineerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket resolved"})
}

/* =====================
   ADMIN CLOSE
===================== */

func (h *TicketHandler) CloseTicket(c *gin.Context) {
	ticketID, _ := uuid.Parse(c.Param("id"))
	adminID := c.MustGet("user_id").(uuid.UUID)

	if err := h.service.CloseTicket(ticketID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket closed"})
}

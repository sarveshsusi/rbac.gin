package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/models"
	"rbac/service"
	"rbac/utils"
)

type TicketHandler struct {
	service  *service.TicketService
	uploader utils.ImageUploader
}

func NewTicketHandler(
	s *service.TicketService,
	uploader utils.ImageUploader,
) *TicketHandler {
	return &TicketHandler{
		service:  s,
		uploader: uploader,
	}
}

/*
	=========================
	  CUSTOMER: CREATE TICKET

=========================
*/
type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	customerID := c.MustGet("user_id").(uuid.UUID)

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Customer only provides title and description
	// Admin will assign product, AMC, priority, etc. later
	ticket, err := h.service.CreateCustomerTicket(
		customerID,
		req.Title,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

/*
	=========================
	  ADMIN: CREATE TICKET

=========================
*/
func (h *TicketHandler) AdminCreateTicket(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTicket, err := h.service.AdminCreateTicket(&ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTicket)
}

/*
	=========================
	  ADMIN: ASSIGN TICKET

=========================
*/
type AssignTicketRequest struct {
	EngineerID      uuid.UUID              `json:"engineer_id" binding:"required"`
	Priority        models.TicketPriority  `json:"priority" binding:"required"`
	SupportMode     models.SupportMode     `json:"support_mode" binding:"required"`
	ServiceCallType models.ServiceCallType `json:"service_call_type" binding:"required"`
}

func (h *TicketHandler) AssignTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	adminID := c.MustGet("user_id").(uuid.UUID)

	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AssignTicket(
		ticketID,
		req.EngineerID,
		adminID,
		req.Priority,
		req.SupportMode,
		req.ServiceCallType,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket assigned successfully"})
}

/*
	=========================
	  SUPPORT: START TICKET

=========================
*/
func (h *TicketHandler) StartTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	if err := h.service.StartTicket(ticketID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ticket started"})
}

/*
	=========================
	  SUPPORT: CLOSE TICKET (WITH PROOF)

=========================
*/
func (h *TicketHandler) CloseTicket(c *gin.Context) {
	ticketID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	file, err := c.FormFile("proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "proof image required"})
		return
	}

	// Upload to ImageKit
	url, err := h.uploader.Upload(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image upload failed"})
		return
	}

	if err := h.service.CloseTicket(ticketID, url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "ticket closed successfully",
		"proof_url": url,
	})
}

func (h *TicketHandler) GetAdminTickets(c *gin.Context) {
	tickets, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tickets"})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

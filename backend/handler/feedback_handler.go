package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rbac/service"
)

type FeedbackHandler struct {
	service *service.FeedbackService
}

func NewFeedbackHandler(s *service.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: s}
}

func (h *FeedbackHandler) Submit(c *gin.Context) {
	ticketID := uuid.MustParse(c.Param("id"))

	var req struct {
		EngineerID uuid.UUID `json:"engineer_id"`
		Rating     int       `json:"rating"`
		Comment    string    `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Submit(
		ticketID,
		req.EngineerID,
		req.Rating,
		req.Comment,
	); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "feedback submitted"})
}

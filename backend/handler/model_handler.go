package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/service"
)

type ModelHandler struct {
	service *service.ModelService
}

func NewModelHandler(s *service.ModelService) *ModelHandler {
	return &ModelHandler{service: s}
}

func (h *ModelHandler) GetByBrand(c *gin.Context) {
	brandID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid brand id"})
		return
	}

	data, err := h.service.GetByBrand(brandID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *ModelHandler) Create(c *gin.Context) {
	var body struct {
		Name    string    `json:"name" binding:"required"`
		BrandID uuid.UUID `json:"brand_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.service.Create(body.Name, body.BrandID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, data)
}

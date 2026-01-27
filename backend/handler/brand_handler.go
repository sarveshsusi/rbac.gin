package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/service"
)

type BrandHandler struct {
	service *service.BrandService
}

func NewBrandHandler(s *service.BrandService) *BrandHandler {
	return &BrandHandler{service: s}
}

/* =========================
   GET BRANDS BY CATEGORY
========================= */
func (h *BrandHandler) GetByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category id",
		})
		return
	}

	data, err := h.service.GetByCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

/* =========================
   CREATE BRAND (FIXED)
========================= */
func (h *BrandHandler) Create(c *gin.Context) {
	var body struct {
		Name       string    `json:"name" binding:"required"`
		CategoryID uuid.UUID `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	brand, err := h.service.Create(body.Name, body.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, brand)
}

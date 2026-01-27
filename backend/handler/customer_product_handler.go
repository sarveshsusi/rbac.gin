package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/service"
)

type CustomerProductHandler struct {
	service *service.CustomerProductService
}

func NewCustomerProductHandler(
	service *service.CustomerProductService,
) *CustomerProductHandler {
	return &CustomerProductHandler{service: service}
}

func (h *CustomerProductHandler) AssignToCustomer(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid customer id",
		})
		return
	}

	var body struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.AssignProductToCustomer(
		customerID,
		body.ProductID,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product assigned to customer successfully",
	})
}

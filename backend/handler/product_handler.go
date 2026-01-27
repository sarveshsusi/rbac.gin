package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/middleware"
	"rbac/models"
	"rbac/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

/* =========================
   CREATE PRODUCT (ADMIN)
========================= */
func (h *ProductHandler) Create(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get(middleware.CtxUserID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	adminID := userIDValue.(uuid.UUID)

	product, err := h.service.CreateProduct(&req, adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}


/* =========================
   GET ALL PRODUCTS (ADMIN)
========================= */
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

/* =========================
   ASSIGN PRODUCT TO CUSTOMER
========================= */
func (h *ProductHandler) AssignToCustomer(c *gin.Context) {
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

	if err := h.service.AssignProductToCustomer(customerID, body.ProductID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product assigned to customer",
	})
}

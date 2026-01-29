package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"rbac/service"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}

/* =========================
   CREATE CUSTOMER (ADMIN)
========================= */
func (h *CustomerHandler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Company  string `json:"company" binding:"required"`
		Phone    string `json:"phone"`
		Address  string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.CreateCustomer(
		req.Name,
		req.Email,
		req.Password,
		req.Company,
		req.Phone,
		req.Address,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create customer",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "customer created successfully",
	})
}

/* =========================
   GET ALL CUSTOMERS (ADMIN)
========================= */
func (h *CustomerHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	customers, total, err := h.service.GetAllCustomers(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch customers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
		"meta": gin.H{
			"page":  page,
			"limit": 3,
			"total": total,
		},
	})
}

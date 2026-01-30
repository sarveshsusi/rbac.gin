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

func (h *CustomerProductHandler) GetCustomerProducts(c *gin.Context) {
	idStr := c.Param("id")
	// The route is /admin/customers/:id/products
	// But :id here refers to CUSTOMER UUID (from customers table), NOT USER UUID.
	// HOWEVER, the frontend might be sending USER UUID if selecting from user list.
	// AdminCreateTicket.jsx uses selectedCustomer which is UserID from /admin/users endpoints.
	// So we might need to find CustomerID from UserID first if the input is UserID.

	// Assuming logic: if we pass UserID, we resolve it.
	targetID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Try to find if this is a userID that has a customer profile
	// Ideally we should fix frontend to pass CustomerID, but simpler to resolve here
	// NOTE: CustomerProductRepository.GetByCustomerID expects CUSTOMER ID.

	// Since we don't have easy access to CustomerRepo here directly (it is in Service),
	// let's rely on Service to handle "Get by UserID" or similar logic, OR just try to use it as CustomerID.
	// Given the previous flow: AdminCreateTicket uses `c.id` from `api.get("/admin/users?role=customer")`.
	// The `GetAllUsers` returns User UUID.
	// So `targetID` is `User.ID`.
	// We need to fetch `Customer` by `User.ID` to get `Customer.ID`.

	// Let's rely on service to handle this translation or add a method.
	// BUT, strict looking at Service layer `AssignProductToCustomer` it takes `userID`.
	// Let's modify service `GetCustomerProducts` to take UserID for consistency with the Assign method.

	// WAIT: I updated `GetCustomerProducts` in Service to take `customerID`.
	// I should probably change it to take `userID` to be safe and consistent with `AssignProductToCustomer`.

	// Let's update Service to GetByUserID logic.

	products, err := h.service.GetCustomerProductsByUserID(targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

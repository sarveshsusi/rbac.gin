package repository

import (
	"rbac/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerProductRepository struct {
	db *gorm.DB
}

func NewCustomerProductRepository(db *gorm.DB) *CustomerProductRepository {
	return &CustomerProductRepository{db: db}
}

func (r *CustomerProductRepository) Exists(
	customerID, productID uuid.UUID,
) (bool, error) {

	var count int64
	err := r.db.
		Model(&models.CustomerProduct{}).
		Where(
			"customer_id = ? AND product_id = ? AND is_active = true",
			customerID,
			productID,
		).
		Count(&count).Error

	return count > 0, err
}

func (r *CustomerProductRepository) Assign(
	customerID, productID uuid.UUID,
) error {

	cp := models.CustomerProduct{
		CustomerID: customerID,
		ProductID:  productID,
		IsActive:   true,
	}

	return r.db.Create(&cp).Error
}

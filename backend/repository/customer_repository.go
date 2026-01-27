package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"rbac/models"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// Create customer profile for a user
func (r *CustomerRepository) Create(userID uuid.UUID) error {
	customer := &models.Customer{
		UserID:   userID,
		IsActive: true,
	}

	return r.db.Create(customer).Error
}

// Optional helper (used elsewhere)
func (r *CustomerRepository) FindByUserID(userID uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.
		Where("user_id = ?", userID).
		First(&customer).Error; err != nil {
		return nil, errors.New("customer profile not found")
	}
	return &customer, nil
}

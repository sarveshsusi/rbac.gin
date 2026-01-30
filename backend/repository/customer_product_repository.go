package repository

import (
	"rbac/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *CustomerProductRepository) GetCustomerByUserID(
	userID uuid.UUID,
) (*models.Customer, error) {

	var customer models.Customer
	err := r.db.
		Where("user_id = ?", userID).
		First(&customer).Error

	return &customer, err
}

func (r *CustomerProductRepository) Assign(
	customerID uuid.UUID,
	productID uuid.UUID,
) error {

	return r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "customer_id"},
				{Name: "product_id"},
			},
			DoNothing: true,
		}).
		Create(&models.CustomerProduct{
			CustomerID: customerID,
			ProductID:  productID,
		}).Error
}

func (r *CustomerProductRepository) GetByCustomerID(customerID uuid.UUID) ([]models.CustomerProduct, error) {
	var results []models.CustomerProduct
	err := r.db.Preload("Product").Where("customer_id = ? AND is_active = true", customerID).
		Find(&results).Error
	return results, err
}

package repository

import (
	"rbac/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) IsBrandInCategory(
	brandID, categoryID string,
) (bool, error) {

	var count int64
	err := r.db.
		Table("brand_categories").
		Where("brand_id = ? AND category_id = ?", brandID, categoryID).
		Count(&count).Error

	return count > 0, err
}

func (r *ProductRepository) IsModelInBrand(
	modelID, brandID string,
) (bool, error) {

	var count int64
	err := r.db.
		Table("models").
		Where("id = ? AND brand_id = ?", modelID, brandID).
		Count(&count).Error

	return count > 0, err
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetAll(products *[]models.Product) error {
	return r.db.Find(products).Error
}

func (r *ProductRepository) AssignToCustomer(
	customerID, productID uuid.UUID,
) error {

	return r.db.Table("customer_products").Create(map[string]interface{}{
		"customer_id": customerID,
		"product_id":  productID,
	}).Error
}

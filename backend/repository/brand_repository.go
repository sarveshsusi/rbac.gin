package repository

import (
	"rbac/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BrandRepository struct {
	db *gorm.DB
}

func NewBrandRepository(db *gorm.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

func (r *BrandRepository) GetByCategory(categoryID uuid.UUID) ([]models.Brand, error) {
	var brands []models.Brand
	err := r.db.
		Joins("JOIN brand_categories bc ON bc.brand_id = brands.id").
		Where("bc.category_id = ?", categoryID).
		Find(&brands).Error
	return brands, err
}

func (r *BrandRepository) Create(name string) (*models.Brand, error) {
	brand := &models.Brand{Name: name}
	return brand, r.db.Create(brand).Error
}

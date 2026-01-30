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

/*
	=========================
	  CREATE BRAND

=========================
*/
func (r *BrandRepository) Create(name string) (*models.Brand, error) {
	brand := &models.Brand{
		Name: name,
	}
	if err := r.db.Create(brand).Error; err != nil {
		return nil, err
	}
	return brand, nil
}

/*
	=========================
	  GET ALL BRANDS

=========================
*/
func (r *BrandRepository) GetAll() ([]models.Brand, error) {
	var brands []models.Brand
	err := r.db.Find(&brands).Error
	return brands, err
}

/*
	=========================
	  ASSIGN BRAND → CATEGORY

=========================
*/
func (r *BrandRepository) AssignToCategory(
	brandID uuid.UUID,
	categoryID uuid.UUID,
) error {

	return r.db.Table("brand_categories").Create(map[string]interface{}{
		"brand_id":    brandID,
		"category_id": categoryID,
	}).Error
}

/*
	=========================
	  GET BRANDS BY CATEGORY

=========================
*/
func (r *BrandRepository) GetByCategory(
	categoryID uuid.UUID,
) ([]models.Brand, error) {

	var brands []models.Brand

	err := r.db.
		Joins("JOIN brand_categories bc ON bc.brand_id = brands.id").
		Where("bc.category_id = ?", categoryID).
		Find(&brands).Error

	return brands, err
}

/*
	=========================
	  VALIDATION: BRAND ↔ CATEGORY

=========================
*/
func (r *BrandRepository) IsAllowedForCategory(
	brandID uuid.UUID,
	categoryID uuid.UUID,
) (bool, error) {

	var count int64

	err := r.db.
		Table("brand_categories").
		Where("brand_id = ? AND category_id = ?", brandID, categoryID).
		Count(&count).Error

	return count > 0, err
}

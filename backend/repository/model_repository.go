package repository

import (
	"rbac/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) GetByBrand(brandID uuid.UUID) ([]models.Model, error) {
	var modelsList []models.Model
	err := r.db.
		Where("brand_id = ?", brandID).
		Find(&modelsList).Error
	return modelsList, err
}

func (r *ModelRepository) Create(name string, brandID uuid.UUID) (*models.Model, error) {
	model := &models.Model{
		Name:    name,
		BrandID: brandID,
	}
	return model, r.db.Create(model).Error
}

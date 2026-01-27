package repository

import (
	"gorm.io/gorm"
	"rbac/models"
)

type FeedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) Create(f *models.TicketFeedback) error {
	return r.db.Create(f).Error
}

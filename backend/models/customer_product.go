package models

import (
	"time"

	"github.com/google/uuid"
)

type CustomerProduct struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CustomerID uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_customer_product"`
	ProductID  uuid.UUID `gorm:"type:uuid;index;uniqueIndex:idx_customer_product"`
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time

	Product Product `gorm:"foreignKey:ProductID"`
}

func (CustomerProduct) TableName() string {
	return "customer_products"
}

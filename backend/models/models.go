package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSupport  Role = "support"
	RoleCustomer Role = "customer"
)

type User struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email             string     `gorm:"uniqueIndex;not null"`
	Password          string     `gorm:"not null"`
	Role              Role       `gorm:"type:varchar(20);not null"`
	IsActive          bool       `gorm:"default:true"`
	MustResetPassword bool       `gorm:"default:false"`
	CreatedBy         *uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsRevoked bool      `gorm:"default:false"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

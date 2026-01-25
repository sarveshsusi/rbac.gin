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
	Name              string    `gorm:"type:varchar(100)"` // 👈 pointer
	Email             string     `gorm:"uniqueIndex;not null"`
	Password          string     `gorm:"not null"`
	Role              Role       `gorm:"type:varchar(20);not null"`
	IsActive          bool       `gorm:"default:true"`
	MustResetPassword bool       `gorm:"default:false"`
	CreatedBy         *uuid.UUID `gorm:"type:uuid;index"`
	TwoFAEnabled      bool       `gorm:"column:two_fa_enabled;default:false"`
	LastLoginAt       *time.Time
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

type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}

type TwoFAOTP struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"index"`
	Code      string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}
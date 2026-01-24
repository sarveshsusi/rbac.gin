package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"rbac/models"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

/* =====================
   User Queries
===================== */

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.
		Where("email = ? AND is_active = true", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) FindUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.
		Where("id = ?", userID).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) UpdateUserPassword(
	userID uuid.UUID,
	hashedPassword string,
	mustReset bool,
) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":             hashedPassword,
			"must_reset_password":  mustReset,
		}).Error
}

/* =====================
   Refresh Tokens
===================== */

func (r *AuthRepository) CreateRefreshToken(rt *models.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *AuthRepository) RevokeAllUserTokens(userID uuid.UUID) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func (r *AuthRepository) FindRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken

	err := r.db.
		Preload("User").
		Where(
			"token = ? AND is_revoked = false AND expires_at > ?",
			tokenHash,
			time.Now(),
		).
		First(&rt).
		Error

	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	return &rt, nil
}

func (r *AuthRepository) RevokeRefreshToken(tokenHash string) error {
	result := r.db.Model(&models.RefreshToken{}).
		Where("token = ? AND is_revoked = false", tokenHash).
		Update("is_revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("refresh token already revoked or not found")
	}

	return nil
}

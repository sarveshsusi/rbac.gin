package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"rbac/config"
	"rbac/models"
	"rbac/repository"
	"rbac/utils"
)

type AuthService struct {
	repo *repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo *repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

/*
=====================
 Response DTOs
=====================
*/

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID    uuid.UUID  `json:"id"`
	Email string     `json:"email"`
	Role  models.Role `json:"role"`
}

/*
=====================
 Login
=====================
*/

func (s *AuthService) Login(
	email string,
	password string,
	ip string,
	userAgent string,
) (*LoginResponse, error) {

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(
		user,
		s.cfg.JWT.AccessSecret,
		s.cfg.JWT.AccessExpiry,
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (RAW)
	refreshRaw, err := utils.GenerateRefreshToken(
		s.cfg.JWT.RefreshSecret,
		s.cfg.JWT.RefreshExpiry,
	)
	if err != nil {
		return nil, err
	}

	// Store HASHED refresh token
	if err := s.repo.CreateRefreshToken(&models.RefreshToken{
		UserID:    user.ID,
		Token:     utils.HashToken(refreshRaw),
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
		User: &UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

/*
=====================
 Refresh Token (ROTATION)
=====================
*/

func (s *AuthService) RefreshAccessToken(
	oldRaw string,
	ip string,
	userAgent string,
) (*LoginResponse, error) {

	oldHash := utils.HashToken(oldRaw)

	rt, err := s.repo.FindRefreshToken(oldHash)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 🔒 Revoke old refresh token (rotation)
	if err := s.repo.RevokeRefreshToken(oldHash); err != nil {
		return nil, err
	}

	// New access token
	newAccess, err := utils.GenerateAccessToken(
		&rt.User,
		s.cfg.JWT.AccessSecret,
		s.cfg.JWT.AccessExpiry,
	)
	if err != nil {
		return nil, err
	}

	// New refresh token
	newRefresh, err := utils.GenerateRefreshToken(
		s.cfg.JWT.RefreshSecret,
		s.cfg.JWT.RefreshExpiry,
	)
	if err != nil {
		return nil, err
	}

	// Store new refresh token hash
	if err := s.repo.CreateRefreshToken(&models.RefreshToken{
		UserID:    rt.UserID,
		Token:     utils.HashToken(newRefresh),
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		User: &UserInfo{
			ID:    rt.User.ID,
			Email: rt.User.Email,
			Role:  rt.User.Role,
		},
	}, nil
}

/*
=====================
 Logout
=====================
*/

func (s *AuthService) Logout(
	refreshRaw string,
	userID uuid.UUID,
	role models.Role,
	ip string,
	userAgent string,
) error {
	return s.repo.RevokeRefreshToken(utils.HashToken(refreshRaw))
}

/*
=====================
 Admin: Create User
=====================
*/

func (s *AuthService) CreateUser(
	email string,
	role models.Role,
	createdBy uuid.UUID,
	ip string,
	userAgent string,
) (*models.User, string, error) {

	existing, _ := s.repo.FindUserByEmail(email)
	if existing != nil {
		return nil, "", errors.New("user already exists")
	}

	// Generate temporary password
	tempPassword, err := utils.GenerateRandomToken(12)
	if err != nil {
		return nil, "", err
	}

	// Hash password
	hashed, err := utils.HashPassword(tempPassword)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Email:             email,
		Password:          hashed,
		Role:              role,
		IsActive:          true,
		MustResetPassword: true,
		CreatedBy:         &createdBy,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, "", err
	}

	return user, tempPassword, nil
}

/*
=====================
 Change Password
=====================
*/

func (s *AuthService) ChangePassword(
	userID uuid.UUID,
	oldPassword string,
	newPassword string,
	role models.Role,
	ip string,
	userAgent string,
) error {

	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := utils.CheckPassword(oldPassword, user.Password); err != nil {
		return errors.New("invalid old password")
	}

	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateUserPassword(
		userID,
		hashed,
		false,
	); err != nil {
		return err
	}

	// 🔒 Revoke all refresh tokens
	if err := s.repo.RevokeAllUserTokens(userID); err != nil {
		return err
	}

	return nil
}

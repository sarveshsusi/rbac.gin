package service

import (
	
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"rbac/config"
	"rbac/models"
	"rbac/repository"
	"rbac/utils"
	"math/rand"
)

type AuthService struct {
	repo   *repository.AuthRepository
	mailer *utils.Mailer
	cfg    *config.Config
}

func NewAuthService(
	repo *repository.AuthRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		repo:   repo,
		mailer: utils.NewMailer(cfg.Mail),
		cfg:    cfg,
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
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
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

	// 🚫 Block login until password reset
	if user.MustResetPassword {
		return nil, errors.New("PASSWORD_RESET_REQUIRED")
	}

	// 🔐 Require 2FA only AFTER first successful login
	if user.TwoFAEnabled && user.LastLoginAt != nil {

	if err := s.sendOTP(user); err != nil {
		return nil, err
	}

	twoFAToken, err := utils.Generate2FAToken(
		user.ID,
		s.cfg.JWT.AccessSecret,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  "",
		RefreshToken: "",
		User: &UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, errors.New("TWO_FA_REQUIRED:" + twoFAToken)
}


	// ✅ Normal login (first time OR 2FA disabled)
	accessToken, err := utils.GenerateAccessToken(
		user,
		s.cfg.JWT.AccessSecret,
		s.cfg.JWT.AccessExpiry,
	)
	if err != nil {
		return nil, err
	}

	refreshRaw, err := utils.GenerateRefreshToken(
		s.cfg.JWT.RefreshSecret,
		s.cfg.JWT.RefreshExpiry,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateRefreshToken(&models.RefreshToken{
		UserID:    user.ID,
		Token:     utils.HashToken(refreshRaw),
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshExpiry),
	}); err != nil {
		return nil, err
	}

	// 🕒 Update last login ONLY after success
	now := time.Now()
	_ = s.repo.UpdateLastLogin(user.ID, &now)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
		User: &UserInfo{
			ID:    user.ID,
			Name:  user.Name,
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

func (s *AuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindUserByID(id)
}

func (s *AuthService) CreatePasswordReset(
	user *models.User,
) (string, error) {

	rawToken, err := utils.GenerateRandomToken(48)
	if err != nil {
		return "", err
	}

	reset := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     utils.HashToken(rawToken),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.CreatePasswordReset(reset); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *AuthService) ResetPassword(
	rawToken string,
	newPassword string,
) error {

	hashed := utils.HashToken(rawToken)

	reset, err := s.repo.FindValidPasswordReset(hashed)
	if err != nil {
		return errors.New("invalid or expired reset link")
	}

	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	hashedPwd, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateUserPassword(
		reset.UserID,
		hashedPwd,
		false, // 🔥 CLEAR must_reset_password
	); err != nil {
		return err
	}

	// 🔒 Revoke all sessions
	_ = s.repo.RevokeAllUserTokens(reset.UserID)

	return s.repo.MarkPasswordResetUsed(reset.ID)
}


func (s *AuthService) CreateUser(
	email string,
	role models.Role,
	createdBy uuid.UUID,
	ip string,
	userAgent string,
) (*models.User, error) {

	existing, _ := s.repo.FindUserByEmail(email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// 🔒 Generate temporary random password (never shown)
	tempPassword, err := utils.GenerateRandomToken(16)
	if err != nil {
		return nil, err
	}

	hashed, err := utils.HashPassword(tempPassword)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// 🔑 Create reset token
	resetToken, err := s.CreatePasswordReset(user)
	if err != nil {
		return nil, err
	}

	resetURL := s.cfg.FrontendURL + "/reset-password?token=" + resetToken

	// 📧 Send email
	body := `
		<h2>Welcome to RBAC App</h2>
		<p>Your account has been created.</p>
		<p>
			<a href="` + resetURL + `">Click here to set your password</a>
		</p>
		<p>This link expires in 24 hours.</p>
   	`

	if s.mailer == nil {
	return nil, errors.New("email service not configured")
}
	if err := s.mailer.Send(
	user.Email,
	"Set your password",
	body,
); err != nil {
	return nil, err
}

	return user, nil
}

func (s *AuthService) SendPasswordResetEmail(email string) error {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		// 🔒 Silently ignore (prevent user enumeration)
		return nil
	}

	// Create reset token
	token, err := s.CreatePasswordReset(user)
	if err != nil {
		return err
	}

	resetURL := s.cfg.FrontendURL + "/reset-password?token=" + token

	body := `
		<h2>Password Reset</h2>
		<p>Click the link below to reset your password:</p>
		<p>
			<a href="` + resetURL + `">Reset Password</a>
		</p>
		<p>This link expires in 24 hours.</p>
	`

	if s.mailer == nil {
		return errors.New("email service not configured")
	}

	return s.mailer.Send(
		user.Email,
		"Reset your password",
		body,
	)
}
func (s *AuthService) sendOTP(user *models.User) error {

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	hashed := utils.HashToken(code)

	_ = s.repo.MarkAllOTPUsed(user.ID)

	otp := &models.TwoFAOTP{
		UserID:    user.ID,
		Code:      hashed,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.repo.Create2FAOTP(otp); err != nil {
		return err
	}

	body := fmt.Sprintf(`
		<h2>Security Code</h2>
		<h1>%s</h1>
		<p>Expires in 5 minutes</p>
	`, code)

	return s.mailer.Send(user.Email, "Your login code", body)
}
func (s *AuthService) Verify2FA(userID uuid.UUID, code string) (*LoginResponse, error) {

	// 🔐 HASH OTP (enterprise-grade)
	hashed := utils.HashToken(code)

	otp, err := s.repo.FindValid2FAOTP(userID, hashed)
	if err != nil {
		return nil, errors.New("invalid or expired otp")
	}

	// ✅ Mark OTP as used
	_ = s.repo.MarkOTPUsed(otp.ID)

	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	now := time.Now()
	_ = s.repo.UpdateLastLogin(user.ID, &now)

	return s.issueTokens(user)
}

func (s *AuthService) issueTokens(user *models.User) (*LoginResponse, error) {

	accessToken, err := utils.GenerateAccessToken(
		user,
		s.cfg.JWT.AccessSecret,
		s.cfg.JWT.AccessExpiry,
	)
	if err != nil {
		return nil, err
	}

	refreshRaw, err := utils.GenerateRefreshToken(
		s.cfg.JWT.RefreshSecret,
		s.cfg.JWT.RefreshExpiry,
	)
	if err != nil {
		return nil, err
	}

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
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
func (s *AuthService) Enable2FA(userID uuid.UUID) error {
	return s.repo.Enable2FA(userID)
}

func (s *AuthService) Disable2FA(userID uuid.UUID) error {
	return s.repo.Disable2FA(userID)
} 


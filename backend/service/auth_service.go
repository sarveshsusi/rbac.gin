package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"math/big"
	"rbac/config"
	"rbac/models"
	"rbac/repository"
	"rbac/utils"
)

type AuthService struct {
	db           *gorm.DB
	repo         *repository.AuthRepository
	deviceRepo   *repository.RememberedDeviceRepo
	customerRepo *repository.CustomerRepository
	mailer       *utils.Mailer
	cfg          *config.Config
}

func NewAuthService(
	db *gorm.DB,
	repo *repository.AuthRepository,
	deviceRepo *repository.RememberedDeviceRepo,
	customerRepo *repository.CustomerRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		db:           db,
		repo:         repo,
		deviceRepo:   deviceRepo,
		customerRepo: customerRepo,
		mailer:       utils.NewMailer(cfg.Mail),
		cfg:          cfg,
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
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  models.Role `json:"role"`
}

type GetuserInfo struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      models.Role `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	IsActive  bool        `json:"is_active"`
}

/*
=====================
 Login
=====================
*/

func (s *AuthService) Login(
	c *gin.Context,
	email string,
	password string,
	rememberDevice bool, // 👈 checkbox from login page
	ip string,
	userAgent string,
) (*LoginResponse, error) {

	// 1️⃣ Find user
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 2️⃣ Verify password
	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3️⃣ Force password reset
	if user.MustResetPassword {
		return nil, errors.New("PASSWORD_RESET_REQUIRED")
	}

	// 4️⃣ 2FA FLOW (only after first successful login)
	if user.TwoFAEnabled && user.LastLoginAt != nil {

		// 🔍 Check remembered device FIRST
		deviceToken, _ := c.Cookie("remember_device")
		if deviceToken != "" {
			hashed := utils.HashRememberDeviceToken(deviceToken)
			if s.deviceRepo.ExistsValid(user.ID, hashed) {
				return s.issueTokens(user)
			}
		}

		// ❌ Device not trusted → send OTP
		if err := s.sendOTP(user); err != nil {
			return nil, err
		}

		// 🔐 Issue short-lived 2FA token (carry remember intent)
		twoFAToken, err := utils.Generate2FAToken(
			user.ID,
			rememberDevice,
			s.cfg.JWT.AccessSecret,
		)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("TWO_FA_REQUIRED:" + twoFAToken)
	}

	// 5️⃣ Normal login (first login OR 2FA disabled)
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

	// 6️⃣ First login → enable 2FA & set last login
	now := time.Now()
	if user.LastLoginAt == nil {
		_ = s.repo.UpdateLastLogin(user.ID, &now)
		_ = s.repo.Enable2FA(user.ID)

		// reload user
		user, _ = s.repo.FindUserByID(user.ID)
	}

	log.Println("2FA enabled:", user.TwoFAEnabled)
	log.Println("Last login:", user.LastLoginAt)

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
	name string,
	email string,
	role models.Role,
	createdBy uuid.UUID,
	ip string,
	userAgent string,

	company string,
	phone string,
	address string,
) (*models.User, error) {

	// 1️⃣ Check existing user
	existing, _ := s.repo.FindUserByEmail(email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// 2️⃣ Generate temp password
	tempPassword, err := utils.GenerateRandomToken(16)
	if err != nil {
		return nil, err
	}

	hashed, err := utils.HashPassword(tempPassword)
	if err != nil {
		return nil, err
	}

	var createdUser *models.User

	// 3️⃣ TRANSACTION START
	err = s.db.Transaction(func(tx *gorm.DB) error {

		user := &models.User{
			Name:              name,
			Email:             email,
			Password:          hashed,
			Role:              role,
			IsActive:          true,
			MustResetPassword: true,
			CreatedBy:         &createdBy,
		}

		// ✅ Create user (TX SAFE)
		if err := s.repo.CreateUserTx(tx, user); err != nil {
			return err
		}

		// ✅ Create customer profile ONLY for customers
		if role == models.RoleCustomer {
			customer := &models.Customer{
				UserID:   user.ID,
				Company:  company,
				Phone:    phone,
				Address:  address,
				IsActive: true,
			}

			if err := s.customerRepo.Create(tx, customer); err != nil {
				return err
			}
		}

		createdUser = user
		return nil
	})

	if err != nil {
		return nil, err
	}
	// 🔚 TRANSACTION END

	// 4️⃣ Create password reset token
	resetToken, err := s.CreatePasswordReset(createdUser)
	if err != nil {
		return nil, err
	}

	resetURL := s.cfg.FrontendURL + "/reset-password?token=" + resetToken

	// 5️⃣ Send email
	if s.mailer == nil {
		return nil, errors.New("email service not configured")
	}

	body := `
		<h2>Welcome to RBAC App</h2>
		<p>Your account has been created.</p>
		<p>
			<a href="` + resetURL + `">Click here to set your password</a>
		</p>
		<p>This link expires in 24 hours.</p>
	`

	if err := s.mailer.Send(
		createdUser.Email,
		"Set your password",
		body,
	); err != nil {
		return nil, err
	}

	return createdUser, nil
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

	code, err := generateOTP()
	if err != nil {
		return err
	}

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

func (s *AuthService) Verify2FA(
	userID uuid.UUID,
	code string,
	remember bool, // 👈 from 2FA JWT
	ip string,
	userAgent string,
) (*LoginResponse, *string, error) {

	// 1️⃣ Verify OTP
	hashedOTP := utils.HashToken(code)

	otp, err := s.repo.FindValid2FAOTP(userID, hashedOTP)
	if err != nil {
		return nil, nil, errors.New("invalid or expired otp")
	}

	_ = s.repo.MarkOTPUsed(otp.ID)

	// 2️⃣ Load user
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, nil, errors.New("user not found")
	}

	// 3️⃣ Update last login
	now := time.Now()
	_ = s.repo.UpdateLastLogin(user.ID, &now)

	// 4️⃣ Remember device (ONLY after OTP success)
	var rememberDeviceToken *string

	if remember {
		rawToken, err := utils.GenerateRandomToken(32)
		if err == nil {
			hashed := utils.HashRememberDeviceToken(rawToken)

			_ = s.deviceRepo.Create(&models.RememberedDevice{
				UserID:    user.ID,
				Token:     hashed,
				UserAgent: userAgent,
				IPAddress: ip,
				ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
			})

			rememberDeviceToken = &rawToken
		}
	}

	// 5️⃣ Issue access + refresh tokens
	resp, err := s.issueTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return resp, rememberDeviceToken, nil
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

func generateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

/* =====================
   Admin: Get Users (Paginated)
===================== */

func (s *AuthService) GetUsersPaginated(page int) ([]*GetuserInfo, int64, error) {
	const pageSize = 3

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize

	users, total, err := s.repo.GetUsersPaginated(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*GetuserInfo, 0, len(users))
	for _, u := range users {
		result = append(result, &GetuserInfo{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			IsActive:  u.IsActive,
		})
	}

	return result, total, nil
}

func (s *AuthService) GetUsersByRole(role models.Role) ([]*GetuserInfo, error) {
	users, err := s.repo.GetUsersByRole(role)
	if err != nil {
		return nil, err
	}

	result := make([]*GetuserInfo, 0, len(users))
	for _, u := range users {
		result = append(result, &GetuserInfo{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			IsActive:  u.IsActive,
		})
	}
	return result, nil
}

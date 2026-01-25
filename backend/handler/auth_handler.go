package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/config"
	"rbac/models"
	"rbac/service"
	"rbac/utils"
)

type AuthHandler struct {
	service *service.AuthService
	cfg     *config.Config
}

func NewAuthHandler(service *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
	}
}

/* =====================
   DTOs
===================== */

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type Verify2FARequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Code   string    `json:"code" binding:"required,len=6"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type CreateUserRequest struct {
	Email string      `json:"email" binding:"required,email"`
	Role  models.Role `json:"role" binding:"required,oneof=support customer"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

/* =====================
   Cookie Helpers
===================== */

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	secure := h.cfg.Server.Env == "production"

	c.SetCookie(
		"refresh_token",
		token,
		int(h.cfg.JWT.RefreshExpiry.Seconds()),
		"/",
		"",
		secure,
		true, // HttpOnly
	)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	secure := h.cfg.Server.Env == "production"

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		secure,
		true,
	)
}

/* =====================
   Handlers
===================== */

// Login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.service.Login(
		req.Email,
		req.Password,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "TWO_FA_REQUIRED:") {
		token := strings.TrimPrefix(err.Error(), "TWO_FA_REQUIRED:")

		c.JSON(200, gin.H{
			"two_fa_required": true,
			"two_fa_token": token,
		})
		return
	}
	if err.Error() == "PASSWORD_RESET_REQUIRED" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "password_reset_required",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "invalid credentials",
	})
	return
}

	h.setRefreshCookie(c, resp.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token": resp.AccessToken,
		"user":         resp.User,
	})
}

// Refresh token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		var req RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
			return
		}
		refreshToken = req.RefreshToken
	}

	resp, err := h.service.RefreshAccessToken(
		refreshToken,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	h.setRefreshCookie(c, resp.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token": resp.AccessToken,
		"user":         resp.User,
	})
}

// Logout
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		userID := c.MustGet("user_id").(uuid.UUID)
		role := c.MustGet("user_role").(models.Role)

		_ = h.service.Logout(
			refreshToken,
			userID,
			role,
			c.ClientIP(),
			c.GetHeader("User-Agent"),
		)
	}

	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Admin: Create User
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	createdBy := c.MustGet("user_id").(uuid.UUID)

	user, err := h.service.CreateUser(
		req.Email,
		req.Role,
		createdBy,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"message": "User created. Password setup email sent.",
	})
}


// Change password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	role := c.MustGet("user_role").(models.Role)

	err := h.service.ChangePassword(
		userID,
		req.OldPassword,
		req.NewPassword,
		role,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "password changed, login again"})
}

// Get current user
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,   // ✅ FROM DB
		"email": user.Email,
		"role":  user.Role,
	})
}


func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password reset successful",
	})
}
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 🔒 Do NOT reveal if user exists
	_ = h.service.SendPasswordResetEmail(req.Email)

	c.JSON(http.StatusOK, gin.H{
		"message": "If an account exists, a reset link has been sent.",
	})
}

// handler/auth_handler.go

type VerifyOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// func (h *AuthHandler) Verify2FA(c *gin.Context) {
// 	userID := c.MustGet("temp_user_id").(uuid.UUID)

// 	var req VerifyOTPRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(400, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	resp, err := h.service.Verify2FA(userID, req.Code)
// 	if err != nil {
// 		c.JSON(401, gin.H{"error": err.Error()})
// 		return
// 	}

// 	h.setRefreshCookie(c, resp.RefreshToken)

// 	c.JSON(200, gin.H{
// 		"access_token": resp.AccessToken,
// 		"user":         resp.User,
// 	})
// }


// func (h *AuthHandler) Verify2FA(c *gin.Context) {
// 	var req Verify2FARequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	resp, err := h.service.Verify2FA(req.UserID, req.Code)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	h.setRefreshCookie(c, resp.RefreshToken)

// 	c.JSON(http.StatusOK, gin.H{
// 		"access_token": resp.AccessToken,
// 		"user":         resp.User,
// 	})
// }
func (h *AuthHandler) Enable2FA(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.service.Enable2FA(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA enabled successfully",
	})
}
func (h *AuthHandler) Disable2FA(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.service.Disable2FA(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA disabled successfully",
	})
}
func (h *AuthHandler) Verify2FA(c *gin.Context) {

	rawToken := c.GetHeader("X-2FA-Token")
	if rawToken == "" {
		c.JSON(401, gin.H{"error": "missing 2fa token"})
		return
	}

	claims, err := utils.Parse2FAToken(
		rawToken,
		h.cfg.JWT.AccessSecret,
	)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid 2fa session"})
		return
	}

	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.service.Verify2FA(
		claims.UserID,
		req.Code,
	)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	h.setRefreshCookie(c, resp.RefreshToken)

	c.JSON(200, gin.H{
		"access_token": resp.AccessToken,
		"user": resp.User,
	})
}

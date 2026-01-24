package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac/config"
	"rbac/models"
	"rbac/service"
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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
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

	user, tempPwd, err := h.service.CreateUser(
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
		"id":                 user.ID,
		"email":              user.Email,
		"role":               user.Role,
		"temporary_password": tempPwd,
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
	c.JSON(http.StatusOK, gin.H{
		"id":    c.MustGet("user_id").(uuid.UUID),
		"email": c.GetString("user_email"),
		"role":  c.MustGet("user_role"),
	})
}

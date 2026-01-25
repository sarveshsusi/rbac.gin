package routes

import (
	"github.com/gin-gonic/gin"

	"rbac/config"
	"rbac/handler"
	"rbac/middleware"
)

func SetupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	cfg *config.Config,
) {
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)

		auth.POST(
	"/verify-2fa",
	middleware.Temp2FAMiddleware(cfg),
	authHandler.Verify2FA,
)

	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", authHandler.GetMe)
		protected.POST("/change-password", authHandler.ChangePassword)

		protected.POST("/2fa/enable", authHandler.Enable2FA)
		protected.POST("/2fa/disable", authHandler.Disable2FA)

		admin := protected.Group("/admin")
		admin.Use(middleware.RequireAdmin())
		{
			admin.POST("/users", authHandler.CreateUser)
			admin.GET("/users", authHandler.GetAllUsers)
		}
	}
}

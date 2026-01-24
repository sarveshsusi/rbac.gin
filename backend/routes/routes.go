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

	// ======================
	// Public Auth Routes
	// ======================
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// ======================
	// Protected Routes
	// ======================
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", authHandler.GetMe)
		protected.POST("/change-password", authHandler.ChangePassword)

		// ======================
		// Admin Routes
		// ======================
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireAdmin())
		{
			admin.POST("/users", authHandler.CreateUser)
		}
	}
}

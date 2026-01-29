package routes

import (
	"github.com/gin-gonic/gin"

	"rbac/config"
	"rbac/handler"
	"rbac/middleware"
	"rbac/models"
)

func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,

	// Auth
	authHandler *handler.AuthHandler,

	// Dashboards
	adminDashboard *handler.AdminDashboardHandler,
	supportDashboard *handler.SupportDashboardHandler,
	customerDashboard *handler.CustomerDashboardHandler,

	// Core
	ticketHandler *handler.TicketHandler,
	amcHandler *handler.AMCHandler,
	productHandler *handler.ProductHandler,
	customerProductHandler *handler.CustomerProductHandler,
	feedbackHandler *handler.FeedbackHandler,

	// Lookups
	categoryHandler *handler.CategoryHandler,
	brandHandler *handler.BrandHandler,
	modelHandler *handler.ModelHandler,
) {

	api := r.Group("/api/v1")

	/* =========================
	   AUTH (PUBLIC)
	========================= */
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

	/* =========================
	   PROTECTED (JWT)
	========================= */
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		/* ---------- COMMON ---------- */
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", authHandler.GetMe)
		protected.POST("/change-password", authHandler.ChangePassword)

		protected.POST("/2fa/enable", authHandler.Enable2FA)
		protected.POST("/2fa/disable", authHandler.Disable2FA)

		/* =========================
		   ADMIN
		========================= */
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole(models.RoleAdmin))
		{
			admin.GET("/dashboard", adminDashboard.Dashboard)

			// USERS
			admin.POST("/users", authHandler.CreateUser)
			admin.GET("/users", authHandler.GetAllUsers)

			// PRODUCTS
			admin.POST("/products", productHandler.Create)
			admin.GET("/products", productHandler.GetAll)

			admin.POST(
				"/customers/:id/products",
				customerProductHandler.AssignToCustomer,
			)

			// LOOKUPS
			admin.GET("/categories", categoryHandler.GetAll)
			admin.POST("/categories", categoryHandler.Create)

			admin.GET("/categories/:id/brands", brandHandler.GetByCategory)
			admin.POST("/brands", brandHandler.Create)

			admin.GET("/brands/:id/models", modelHandler.GetByBrand)
			admin.POST("/models", modelHandler.Create)

			// AMC
			admin.POST("/amc", amcHandler.Create)
			admin.GET("/amc", amcHandler.GetAllAMCs)

			// TICKETS
			admin.POST("/tickets/:id/assign", ticketHandler.AssignTicket)
			admin.POST("/tickets/:id/close", ticketHandler.CloseTicket)
		}

		/* =========================
		   SUPPORT
		========================= */
		support := protected.Group("/support")
		support.Use(middleware.RequireRole(models.RoleSupport))
		{
			support.GET("/tickets", supportDashboard.MyTickets)
			support.POST("/tickets/:id/resolve", ticketHandler.ResolveTicket)
		}

		/* =========================
		   CUSTOMER
		========================= */
		customer := protected.Group("/customer")
		customer.Use(middleware.RequireRole(models.RoleCustomer))
		{
			customer.GET("/tickets", customerDashboard.MyTickets)
			customer.POST("/tickets", ticketHandler.CreateTicket)
			customer.POST("/tickets/:id/feedback", feedbackHandler.Submit)
			customer.GET("/amc", amcHandler.GetMyAMCs)
		}
	}
}

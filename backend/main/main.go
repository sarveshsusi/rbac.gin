package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"rbac/config"
	"rbac/database"
	"rbac/handler"
	"rbac/middleware"
	"rbac/repository"
	"rbac/routes"
	"rbac/service"
	"rbac/utils"
)

func main() {
	/* =========================
	   ENV & CONFIG
	========================= */
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system env vars")
	}

	cfg := config.LoadConfig()

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	/* =========================
	   DATABASE
	========================= */
	if err := database.Init(cfg); err != nil {
		log.Fatalf("❌ database init failed: %v", err)
	}

	// Run Migrations
	database.Migrate(database.DB)

	/* =========================
	   GIN ENGINE
	========================= */
	db := database.DB
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	r.Use(middleware.CORSMiddleware([]string{
		"http://localhost:5173",
		"http://localhost:3000",
	}))

	/* =========================
	   REPOSITORIES
	========================= */
	authRepo := repository.NewAuthRepository(database.DB)
	rememberedDeviceRepo := repository.NewRememberedDeviceRepo(database.DB)

	ticketRepo := repository.NewTicketRepository(database.DB)

	amcRepo := repository.NewAMCRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	customerProductRepo := repository.NewCustomerProductRepository(database.DB)
	feedbackRepo := repository.NewFeedbackRepository(database.DB)
	customerRepo := repository.NewCustomerRepository(database.DB)
	dashboardRepo := repository.NewDashboardRepository(database.DB)

	categoryRepo := repository.NewCategoryRepository(database.DB)
	brandRepo := repository.NewBrandRepository(database.DB)
	modelRepo := repository.NewModelRepository(database.DB)

	/* =========================
	   SERVICES
	========================= */
	authService := service.NewAuthService(
		db,
		authRepo,
		rememberedDeviceRepo,
		customerRepo,
		cfg,
	)

	ticketService := service.NewTicketService(ticketRepo)

	adminService := service.NewAdminService(dashboardRepo)
	supportService := service.NewSupportService(ticketRepo)
	customerService := service.NewCustomerService(
		db,
		authRepo,
		customerRepo,
		ticketRepo,
	)

	amcService := service.NewAMCService(amcRepo)
	productService := service.NewProductService(productRepo)
	customerProductService := service.NewCustomerProductService(customerProductRepo)
	feedbackService := service.NewFeedbackService(feedbackRepo)

	categoryService := service.NewCategoryService(categoryRepo)
	brandService := service.NewBrandService(brandRepo)
	modelService := service.NewModelService(modelRepo)

	/* =========================
	   UTILS
	========================= */
	imageUploader := utils.NewImageKitUploader(cfg)

	/* =========================
	   HANDLERS
	========================= */
	authHandler := handler.NewAuthHandler(authService, cfg)

	adminDashboard := handler.NewAdminDashboardHandler(adminService)
	supportDashboard := handler.NewSupportDashboardHandler(supportService)
	customerDashboard := handler.NewCustomerDashboardHandler(customerService)

	ticketHandler := handler.NewTicketHandler(ticketService, imageUploader)
	amcHandler := handler.NewAMCHandler(amcService)
	productHandler := handler.NewProductHandler(productService)
	customerProductHandler := handler.NewCustomerProductHandler(customerProductService)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)

	categoryHandler := handler.NewCategoryHandler(categoryService)
	brandHandler := handler.NewBrandHandler(brandService)
	modelHandler := handler.NewModelHandler(modelService)

	/* =========================
	   ROUTES
	========================= */
	routes.SetupRoutes(
		r,
		cfg,

		// Auth
		authHandler,

		// Dashboards
		adminDashboard,
		supportDashboard,
		customerDashboard,

		// Core
		ticketHandler,
		amcHandler,
		productHandler,
		customerProductHandler,
		feedbackHandler,

		// Lookups
		categoryHandler,
		brandHandler,
		modelHandler,
	)

	/* =========================
	   START SERVER
	========================= */
	log.Printf(
		"🚀 Server running on port %s [%s]",
		cfg.Server.Port,
		cfg.Server.Env,
	)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}

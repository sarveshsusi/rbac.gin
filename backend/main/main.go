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
)

func main() {
	// Load .env (Supabase / local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, relying on system env vars")
	}

	// Load config
	cfg := config.LoadConfig()

	// Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init DB
	if err := database.Init(cfg); err != nil {
		log.Fatalf("❌ database init failed: %v", err)
	}

	// Router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// ✅ Fix proxy warning
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// ✅ CORS (VERY IMPORTANT for frontend)
	r.Use(middleware.CORSMiddleware([]string{
		"http://localhost:5173", // Vite
		"http://localhost:3000",
		// add prod domain later
	}))

	// Dependency injection
	repo := repository.NewAuthRepository(database.DB)
	svc := service.NewAuthService(repo, cfg)
	authHandler := handler.NewAuthHandler(svc, cfg)

	// Routes
	routes.SetupRoutes(r, authHandler, cfg)

	log.Printf(
		"🚀 Server running on port %s [%s]",
		cfg.Server.Port,
		cfg.Server.Env,
	)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}

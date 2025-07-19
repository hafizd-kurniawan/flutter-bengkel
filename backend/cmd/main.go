package main

//go:generate swag init

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/database"
	"flutter-bengkel/internal/handlers"
	"flutter-bengkel/internal/middleware"
	"flutter-bengkel/internal/repositories"
	"flutter-bengkel/internal/services"

	_ "flutter-bengkel/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title Workshop Management System API
// @version 1.0
// @description Complete workshop management system with service, sparepart, and vehicle trading
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@bengkel.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db.GetDB(), "migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	repos := repositories.New(db.GetDB())

	// Initialize services
	svc := services.New(repos, cfg)

	// Initialize handlers
	h := handlers.New(svc, cfg)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Workshop Management System API is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Setup routes
	h.SetupRoutes(api)

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	log.Printf("🚀 Server started on port %s", cfg.Server.Port)
	log.Printf("📚 Swagger documentation: http://localhost:%s/swagger/", cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("✅ Server exited")
}
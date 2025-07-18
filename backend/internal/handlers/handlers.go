package handlers

import (
	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	services *services.Services
	config   *config.Config
}

// New creates a new handlers instance
func New(services *services.Services, config *config.Config) *Handlers {
	return &Handlers{
		services: services,
		config:   config,
	}
}

// SetupRoutes sets up all API routes
func (h *Handlers) SetupRoutes(api fiber.Router) {
	// Auth routes (no authentication required)
	auth := api.Group("/auth")
	h.setupAuthRoutes(auth)

	// Protected routes
	protected := api.Use(h.jwtMiddleware())

	// User management routes
	users := protected.Group("/users")
	h.setupUserRoutes(users)

	// Customer management routes
	customers := protected.Group("/customers")
	h.setupCustomerRoutes(customers)

	// Vehicle management routes
	vehicles := protected.Group("/vehicles")
	h.setupVehicleRoutes(vehicles)

	// Service management routes
	servicesGroup := protected.Group("/services")
	h.setupServiceRoutes(servicesGroup)

	// Product management routes
	products := protected.Group("/products")
	h.setupProductRoutes(products)

	// Service job management routes
	serviceJobs := protected.Group("/service-jobs")
	h.setupServiceJobRoutes(serviceJobs)

	// Transaction management routes
	transactions := protected.Group("/transactions")
	h.setupTransactionRoutes(transactions)

	// Payment management routes
	payments := protected.Group("/payments")
	h.setupPaymentRoutes(payments)

	// Master data routes
	masterData := protected.Group("/master-data")
	h.setupMasterDataRoutes(masterData)
}
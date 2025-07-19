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

	// POS operations (Point of Sale)
	pos := protected.Group("/pos")
	h.setupPOSRoutes(pos)

	// Vehicle trading operations
	vehicleTrading := protected.Group("/vehicle-trading")
	h.setupVehicleTradingRoutes(vehicleTrading)

	// Financial management routes
	financial := protected.Group("/financial")
	h.setupFinancialRoutes(financial)

	// Reports and analytics
	reports := protected.Group("/reports")
	h.setupReportRoutes(reports)

	analytics := protected.Group("/analytics")
	h.setupAnalyticsRoutes(analytics)
}

// setupPOSRoutes sets up POS-related routes
func (h *Handlers) setupPOSRoutes(pos fiber.Router) {
	posHandler := NewPOSHandler()

	// Product operations for POS
	pos.Get("/products", posHandler.GetProducts)
	pos.Get("/products/search", posHandler.SearchProducts)
	pos.Get("/products/:id", posHandler.GetProducts) // Individual product detail

	// Transaction operations for POS
	pos.Post("/transactions", posHandler.CreateTransaction)
	pos.Get("/transactions", posHandler.CreateTransaction) // List transactions
	pos.Get("/transactions/:id", posHandler.CreateTransaction) // Transaction detail
}

// setupVehicleTradingRoutes sets up vehicle trading routes
func (h *Handlers) setupVehicleTradingRoutes(vt fiber.Router) {
	vtHandler := NewVehicleTradingHandler()

	// Vehicle purchase operations
	purchases := vt.Group("/purchases")
	purchases.Get("/", vtHandler.GetVehiclePurchases)
	purchases.Post("/", vtHandler.CreateVehiclePurchase)
	purchases.Get("/:id", vtHandler.GetVehiclePurchases) // Individual purchase detail

	// Vehicle sales inventory
	inventory := vt.Group("/inventory")
	inventory.Get("/", vtHandler.GetVehiclePurchases) // Available vehicles
	inventory.Get("/available", vtHandler.GetVehiclePurchases) // Available only

	// Vehicle sales operations
	sales := vt.Group("/sales")
	sales.Post("/", vtHandler.CreateVehiclePurchase) // Create sale
	sales.Get("/", vtHandler.GetVehiclePurchases) // List sales
}

// setupFinancialRoutes sets up financial management routes
func (h *Handlers) setupFinancialRoutes(financial fiber.Router) {
	// Payment methods
	paymentMethods := financial.Group("/payment-methods")
	paymentMethods.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Payment methods endpoint"})
	})

	// Accounts receivable
	receivables := financial.Group("/receivables")
	receivables.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Receivables endpoint"})
	})
	receivables.Get("/pending", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Pending receivables endpoint"})
	})
	receivables.Get("/overdue", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Overdue receivables endpoint"})
	})
}

// setupReportRoutes sets up reporting routes
func (h *Handlers) setupReportRoutes(reports fiber.Router) {
	// Daily reports
	daily := reports.Group("/daily")
	daily.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Daily reports endpoint"})
	})

	// Period reports
	period := reports.Group("/period")
	period.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Period reports endpoint"})
	})
}

// setupAnalyticsRoutes sets up analytics routes
func (h *Handlers) setupAnalyticsRoutes(analytics fiber.Router) {
	analytics.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Dashboard analytics endpoint"})
	})

	analytics.Get("/sales-trends", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Sales trends endpoint"})
	})

	analytics.Get("/top-products", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Top products endpoint"})
	})
}
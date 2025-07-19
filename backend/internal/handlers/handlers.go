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

	// POS Operations (Kasir-centric) - NEW
	pos := protected.Group("/pos")
	h.setupPOSRoutes(pos)

	// Vehicle Trading Operations - NEW
	vehicleTrading := protected.Group("/vehicle-trading")
	h.setupVehicleTradingRoutes(vehicleTrading)

	// Master data routes
	masterData := protected.Group("/master-data")
	h.setupMasterDataRoutes(masterData)
}

// setupPOSRoutes sets up POS (kasir-centric) routes
func (h *Handlers) setupPOSRoutes(pos fiber.Router) {
	// POS Transactions
	pos.Post("/transactions", h.requireRole("Kasir", "Manager", "Admin"), h.CreatePOSTransaction)
	pos.Get("/products/search", h.requireRole("Kasir", "Manager", "Admin"), h.SearchProducts)
	pos.Put("/transactions/:id/payment", h.requireRole("Kasir", "Manager", "Admin"), h.AddTransactionPayment)
	pos.Post("/transactions/:id/print", h.requireRole("Kasir", "Manager", "Admin"), h.PrintReceipt)
	
	// Service Operations
	pos.Get("/queue", h.requireRole("Kasir", "Manager", "Admin"), h.GetQueueManagement)
	pos.Put("/service-jobs/:id/assign", h.requireRole("Kasir", "Manager", "Admin"), h.AutoAssignMechanic)
	
	// Receivables Management
	pos.Get("/receivables/pending", h.requireRole("Kasir", "Manager", "Admin"), h.GetPendingReceivables)
	pos.Post("/receivables/:id/payment", h.requireRole("Kasir", "Manager", "Admin"), h.RecordReceivablePayment)
	pos.Get("/receivables/paid", h.requireRole("Kasir", "Manager", "Admin"), h.GetPaidReceivables)
	
	// Dashboard
	pos.Get("/dashboard/stats", h.requireRole("Kasir", "Manager", "Admin"), h.GetDashboardStats)
}

// setupVehicleTradingRoutes sets up vehicle trading routes
func (h *Handlers) setupVehicleTradingRoutes(vt fiber.Router) {
	// Vehicle Purchase (Admin/Manager only)
	vt.Post("/purchase", h.requireRole("Manager", "Admin"), h.PurchaseVehicle)
	vt.Put("/:id/service", h.requireRole("Manager", "Admin"), h.LinkServiceRequirement)
	vt.Put("/:id/complete-service", h.requireRole("Manager", "Admin", "Technician"), h.CompleteVehicleService)
	
	// Sales Inventory (Sales team access)
	vt.Get("/inventory", h.requireRole("Sales", "Manager", "Admin"), h.GetSalesInventory)
	vt.Put("/:id/price", h.requireRole("Sales", "Manager", "Admin"), h.UpdateVehiclePrice)
	vt.Get("/:id/profit", h.requireRole("Sales", "Manager", "Admin"), h.GetVehicleProfitCalculation)
	
	// Vehicle Sales
	vt.Post("/sales", h.requireRole("Sales", "Manager", "Admin"), h.CreateVehicleSale)
	vt.Get("/sales", h.requireRole("Sales", "Manager", "Admin"), h.GetVehicleSales)
	
	// Statistics
	vt.Get("/stats", h.requireRole("Sales", "Manager", "Admin"), h.GetVehicleTradingStats)
}
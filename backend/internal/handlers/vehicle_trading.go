package handlers

import (
	"strconv"

	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
)

// Vehicle Trading Operations

// PurchaseVehicle creates a new vehicle purchase record
// POST /api/v1/vehicle-trading/purchase
func (h *Handlers) PurchaseVehicle(c *fiber.Ctx) error {
	var req models.CreateVehiclePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate the request
	if err := h.validateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Get user from JWT context
	userID := c.Locals("user_id").(int64)
	outletID := c.Locals("outlet_id").(int64)

	// Create the vehicle purchase
	purchase, err := h.services.VehicleTrading.PurchaseVehicle(&req, userID, outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to create vehicle purchase",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Vehicle purchased successfully",
		Data:    purchase,
	})
}

// LinkServiceRequirement links service requirement to purchased vehicle
// PUT /api/v1/vehicle-trading/:id/service
func (h *Handlers) LinkServiceRequirement(c *fiber.Ctx) error {
	vehiclePurchaseID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle purchase ID",
		})
	}

	var req struct {
		ServiceRequired bool   `json:"service_required"`
		ServiceJobID    *int64 `json:"service_job_id"`
		Notes           string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Link service requirement
	err = h.services.VehicleTrading.LinkServiceRequirement(vehiclePurchaseID, req.ServiceRequired, req.ServiceJobID, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to link service requirement",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service requirement linked successfully",
	})
}

// GetSalesInventory gets vehicle inventory for sales team
// GET /api/v1/vehicle-trading/inventory
func (h *Handlers) GetSalesInventory(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	brand := c.Query("brand", "")
	model := c.Query("model", "")
	yearStr := c.Query("year", "")
	minPriceStr := c.Query("min_price", "")
	maxPriceStr := c.Query("max_price", "")
	serviceStatus := c.Query("service_status", "") // completed, pending, not_required
	saleStatus := c.Query("sale_status", "")       // Available, Reserved, Sold

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse optional filters
	var year *int
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = &y
		}
	}

	var minPrice, maxPrice *float64
	if minPriceStr != "" {
		if p, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = &p
		}
	}
	if maxPriceStr != "" {
		if p, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = &p
		}
	}

	filters := map[string]interface{}{
		"brand":          brand,
		"model":          model,
		"year":           year,
		"min_price":      minPrice,
		"max_price":      maxPrice,
		"service_status": serviceStatus,
		"sale_status":    saleStatus,
	}

	vehicles, total, err := h.services.VehicleTrading.GetSalesInventory(page, limit, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get sales inventory",
			Error:   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Sales inventory retrieved successfully",
		Data:    vehicles,
		Meta: models.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// UpdateVehiclePrice updates selling price (sales team only)
// PUT /api/v1/vehicle-trading/:id/price
func (h *Handlers) UpdateVehiclePrice(c *fiber.Ctx) error {
	vehiclePurchaseID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle purchase ID",
		})
	}

	var req models.UpdateVehiclePriceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate the request
	if err := h.validateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Check if user has sales role
	userRole := c.Locals("role_name").(string)
	if userRole != "Sales" && userRole != "Manager" && userRole != "Admin" {
		return c.Status(fiber.StatusForbidden).JSON(models.Response{
			Success: false,
			Message: "Only sales team can update vehicle prices",
		})
	}

	// Update vehicle price
	updatedVehicle, err := h.services.VehicleTrading.UpdateVehiclePrice(vehiclePurchaseID, req.SellingPrice, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to update vehicle price",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle price updated successfully",
		Data:    updatedVehicle,
	})
}

// CreateVehicleSale creates a vehicle sale transaction
// POST /api/v1/vehicle-trading/sales
func (h *Handlers) CreateVehicleSale(c *fiber.Ctx) error {
	var req models.CreateVehicleSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate the request
	if err := h.validateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Get sales user from JWT context
	salesUserID := c.Locals("user_id").(int64)

	// Create the vehicle sale
	sale, err := h.services.VehicleTrading.CreateVehicleSale(&req, salesUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to create vehicle sale",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Vehicle sale created successfully",
		Data:    sale,
	})
}

// GetVehicleSales gets vehicle sales history
// GET /api/v1/vehicle-trading/sales
func (h *Handlers) GetVehicleSales(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	salesUserIDStr := c.Query("sales_user_id", "")
	fromDate := c.Query("from_date", "")
	toDate := c.Query("to_date", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var salesUserID *int64
	if salesUserIDStr != "" {
		if id, err := strconv.ParseInt(salesUserIDStr, 10, 64); err == nil {
			salesUserID = &id
		}
	}

	filters := map[string]interface{}{
		"sales_user_id": salesUserID,
		"from_date":     fromDate,
		"to_date":       toDate,
	}

	sales, total, err := h.services.VehicleTrading.GetVehicleSales(page, limit, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get vehicle sales",
			Error:   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Vehicle sales retrieved successfully",
		Data:    sales,
		Meta: models.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// GetVehicleTradingStats gets trading statistics for dashboard
// GET /api/v1/vehicle-trading/stats
func (h *Handlers) GetVehicleTradingStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	outletID := c.Locals("outlet_id").(int64)
	period := c.Query("period", "monthly") // daily, weekly, monthly, yearly

	stats, err := h.services.VehicleTrading.GetTradingStats(userID, outletID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get trading statistics",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Trading statistics retrieved successfully",
		Data:    stats,
	})
}

// GetVehicleProfitCalculation calculates profit for a specific vehicle
// GET /api/v1/vehicle-trading/:id/profit
func (h *Handlers) GetVehicleProfitCalculation(c *fiber.Ctx) error {
	vehiclePurchaseID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle purchase ID",
		})
	}

	sellingPriceStr := c.Query("selling_price", "")
	if sellingPriceStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Selling price is required for profit calculation",
		})
	}

	sellingPrice, err := strconv.ParseFloat(sellingPriceStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid selling price format",
		})
	}

	profitData, err := h.services.VehicleTrading.CalculateVehicleProfit(vehiclePurchaseID, sellingPrice)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to calculate profit",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Profit calculation completed successfully",
		Data:    profitData,
	})
}

// Complete service for vehicle (marks service as completed)
// PUT /api/v1/vehicle-trading/:id/complete-service
func (h *Handlers) CompleteVehicleService(c *fiber.Ctx) error {
	vehiclePurchaseID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle purchase ID",
		})
	}

	var req struct {
		ServiceJobID   *int64  `json:"service_job_id"`
		ServiceCost    float64 `json:"service_cost"`
		CompletionNotes string  `json:"completion_notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Complete vehicle service
	updatedVehicle, err := h.services.VehicleTrading.CompleteVehicleService(vehiclePurchaseID, req.ServiceJobID, req.ServiceCost, req.CompletionNotes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to complete vehicle service",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle service completed successfully",
		Data:    updatedVehicle,
	})
}
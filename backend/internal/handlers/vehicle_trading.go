package handlers

import (
	"strconv"
	"time"

	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// Vehicle Trading Handlers

// @Summary Get available vehicles for sale
// @Description Get all vehicles available for sale with pagination and filtering
// @Tags Vehicle Trading
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param brand query string false "Filter by brand"
// @Param model query string false "Filter by model"
// @Param type query string false "Filter by vehicle type"
// @Param year_from query int false "Filter by year from"
// @Param year_to query int false "Filter by year to"
// @Param price_from query number false "Filter by price from"
// @Param price_to query number false "Filter by price to"
// @Param condition_min query int false "Filter by minimum condition rating"
// @Param condition_max query int false "Filter by maximum condition rating"
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/inventory [get]
func (h *Handlers) GetAvailableVehicles(c *fiber.Ctx) error {
	// Parse search parameters
	searchReq := &models.VehicleSearchRequest{
		Brand:        c.Query("brand"),
		Model:        c.Query("model"),
		Type:         c.Query("type"),
		Status:       "Available", // Only available vehicles
		Page:         1,
		PerPage:      10,
	}

	if page := c.QueryInt("page"); page > 0 {
		searchReq.Page = page
	}
	if perPage := c.QueryInt("per_page"); perPage > 0 {
		searchReq.PerPage = perPage
	}
	if yearFrom := c.QueryInt("year_from"); yearFrom > 0 {
		searchReq.YearFrom = yearFrom
	}
	if yearTo := c.QueryInt("year_to"); yearTo > 0 {
		searchReq.YearTo = yearTo
	}
	if conditionMin := c.QueryInt("condition_min"); conditionMin > 0 {
		searchReq.ConditionMin = conditionMin
	}
	if conditionMax := c.QueryInt("condition_max"); conditionMax > 0 {
		searchReq.ConditionMax = conditionMax
	}

	// Parse price filters
	if priceFromStr := c.Query("price_from"); priceFromStr != "" {
		if priceFrom, err := decimal.NewFromString(priceFromStr); err == nil {
			searchReq.PriceFrom = &priceFrom
		}
	}
	if priceToStr := c.Query("price_to"); priceToStr != "" {
		if priceTo, err := decimal.NewFromString(priceToStr); err == nil {
			searchReq.PriceTo = &priceTo
		}
	}

	vehicles, total, err := h.services.VehicleTrading.SearchVehicles(searchReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get available vehicles",
			Error:   err.Error(),
		})
	}

	totalPages := int((total + int64(searchReq.PerPage) - 1) / int64(searchReq.PerPage))

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Available vehicles retrieved successfully",
		Data:    vehicles,
		Meta: models.PaginationMeta{
			CurrentPage: searchReq.Page,
			PerPage:     searchReq.PerPage,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// @Summary Get vehicle inventory by ID
// @Description Get detailed information about a specific vehicle in inventory
// @Tags Vehicle Trading
// @Security Bearer
// @Param id path int true "Inventory ID"
// @Success 200 {object} models.Response{data=models.VehicleInventory}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/inventory/{id} [get]
func (h *Handlers) GetVehicleInventoryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid inventory ID",
			Error:   err.Error(),
		})
	}

	vehicle, err := h.services.VehicleTrading.GetVehicleInventoryByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Vehicle not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle inventory retrieved successfully",
		Data:    vehicle,
	})
}

// @Summary Record vehicle purchase
// @Description Record the purchase of a vehicle from a customer
// @Tags Vehicle Trading
// @Security Bearer
// @Param request body models.CreateVehiclePurchaseRequest true "Vehicle purchase data"
// @Success 201 {object} models.Response{data=models.VehiclePurchase}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/purchases [post]
func (h *Handlers) CreateVehiclePurchase(c *fiber.Ctx) error {
	var req models.CreateVehiclePurchaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Get user from context
	userID := c.Locals("userID").(int64)
	outletID := c.Locals("outletID").(int64)

	purchase, err := h.services.VehicleTrading.CreateVehiclePurchase(&req, userID, outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to create vehicle purchase",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Vehicle purchase recorded successfully",
		Data:    purchase,
	})
}

// @Summary Record vehicle sale
// @Description Record the sale of a vehicle to a customer
// @Tags Vehicle Trading
// @Security Bearer
// @Param request body models.CreateVehicleSaleRequest true "Vehicle sale data"
// @Success 201 {object} models.Response{data=models.VehicleSale}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/sales [post]
func (h *Handlers) CreateVehicleSale(c *fiber.Ctx) error {
	var req models.CreateVehicleSaleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Get user from context
	userID := c.Locals("userID").(int64)
	outletID := c.Locals("outletID").(int64)

	sale, err := h.services.VehicleTrading.CreateVehicleSale(&req, userID, outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to create vehicle sale",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Vehicle sale recorded successfully",
		Data:    sale,
	})
}

// @Summary Update vehicle inventory
// @Description Update vehicle inventory details like price, condition, or status
// @Tags Vehicle Trading
// @Security Bearer
// @Param id path int true "Inventory ID"
// @Param request body models.UpdateVehicleInventoryRequest true "Update data"
// @Success 200 {object} models.Response{data=models.VehicleInventory}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/inventory/{id} [put]
func (h *Handlers) UpdateVehicleInventory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid inventory ID",
			Error:   err.Error(),
		})
	}

	var req models.UpdateVehicleInventoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	vehicle, err := h.services.VehicleTrading.UpdateVehicleInventory(id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to update vehicle inventory",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle inventory updated successfully",
		Data:    vehicle,
	})
}

// @Summary Update vehicle selling price
// @Description Update the estimated selling price of a vehicle
// @Tags Vehicle Trading
// @Security Bearer
// @Param id path int true "Inventory ID"
// @Param request body object{price=number} true "New price"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/inventory/{id}/price [put]
func (h *Handlers) UpdateVehicleSellingPrice(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid inventory ID",
			Error:   err.Error(),
		})
	}

	var req struct {
		Price decimal.Decimal `json:"price" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	err = h.services.VehicleTrading.UpdateSellingPrice(id, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to update selling price",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Selling price updated successfully",
	})
}

// @Summary Upload vehicle photos
// @Description Upload photos for a vehicle in inventory
// @Tags Vehicle Trading
// @Security Bearer
// @Param id path int true "Inventory ID"
// @Param photos formData file true "Vehicle photos"
// @Success 200 {object} models.Response{data=[]models.VehiclePhoto}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/inventory/{id}/photos [post]
func (h *Handlers) UploadVehiclePhotos(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid inventory ID",
			Error:   err.Error(),
		})
	}

	// Handle file upload
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to parse multipart form",
			Error:   err.Error(),
		})
	}

	// Get user from context
	userID := c.Locals("userID").(int64)

	photos, err := h.services.VehicleTrading.UploadVehiclePhotos(id, form.File["photos"], userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to upload vehicle photos",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle photos uploaded successfully",
		Data:    photos,
	})
}

// @Summary Get vehicle trading profit analysis
// @Description Get profit analysis for vehicle trading operations
// @Tags Vehicle Trading
// @Security Bearer
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param outlet_id query int false "Filter by outlet"
// @Success 200 {object} models.Response{data=models.ProfitAnalysis}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/analytics/profit [get]
func (h *Handlers) GetVehicleTradingProfitAnalysis(c *fiber.Ctx) error {
	// Parse date parameters
	var startDate, endDate time.Time
	var err error

	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.Response{
				Success: false,
				Message: "Invalid start date format. Use YYYY-MM-DD",
				Error:   err.Error(),
			})
		}
	} else {
		// Default to beginning of current month
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	endDateStr := c.Query("end_date")
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.Response{
				Success: false,
				Message: "Invalid end date format. Use YYYY-MM-DD",
				Error:   err.Error(),
			})
		}
	} else {
		// Default to end of current month
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	}

	// Parse outlet filter
	var outletID *int64
	if outletIDStr := c.Query("outlet_id"); outletIDStr != "" {
		if id, err := strconv.ParseInt(outletIDStr, 10, 64); err == nil {
			outletID = &id
		}
	}

	analysis, err := h.services.VehicleTrading.GetProfitAnalysis(startDate, endDate, outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get profit analysis",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Profit analysis retrieved successfully",
		Data:    analysis,
	})
}

// @Summary Get inventory aging report
// @Description Get report showing how long vehicles have been in inventory
// @Tags Vehicle Trading
// @Security Bearer
// @Param outlet_id query int false "Filter by outlet"
// @Success 200 {object} models.Response{data=[]models.InventoryAgingReport}
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/reports/inventory-aging [get]
func (h *Handlers) GetInventoryAgingReport(c *fiber.Ctx) error {
	// Parse outlet filter
	var outletID *int64
	if outletIDStr := c.Query("outlet_id"); outletIDStr != "" {
		if id, err := strconv.ParseInt(outletIDStr, 10, 64); err == nil {
			outletID = &id
		}
	}

	report, err := h.services.VehicleTrading.GetInventoryAgingReport(outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get inventory aging report",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Inventory aging report retrieved successfully",
		Data:    report,
	})
}

// @Summary Get sales performance report
// @Description Get sales performance report for sales team
// @Tags Vehicle Trading
// @Security Bearer
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param outlet_id query int false "Filter by outlet"
// @Success 200 {object} models.Response{data=[]models.SalesPerformanceReport}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/reports/sales-performance [get]
func (h *Handlers) GetSalesPerformanceReport(c *fiber.Ctx) error {
	// Parse date parameters (similar to profit analysis)
	var startDate, endDate time.Time
	var err error

	startDateStr := c.Query("start_date")
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.Response{
				Success: false,
				Message: "Invalid start date format. Use YYYY-MM-DD",
				Error:   err.Error(),
			})
		}
	} else {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	endDateStr := c.Query("end_date")
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.Response{
				Success: false,
				Message: "Invalid end date format. Use YYYY-MM-DD",
				Error:   err.Error(),
			})
		}
	} else {
		now := time.Now()
		endDate = time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
	}

	var outletID *int64
	if outletIDStr := c.Query("outlet_id"); outletIDStr != "" {
		if id, err := strconv.ParseInt(outletIDStr, 10, 64); err == nil {
			outletID = &id
		}
	}

	report, err := h.services.VehicleTrading.GetSalesPerformanceReport(startDate, endDate, outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get sales performance report",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Sales performance report retrieved successfully",
		Data:    report,
	})
}

// @Summary Get vehicle purchases
// @Description Get list of vehicle purchases with pagination
// @Tags Vehicle Trading
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/purchases [get]
func (h *Handlers) GetVehiclePurchases(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)
	offset := (page - 1) * perPage

	purchases, total, err := h.services.VehicleTrading.GetVehiclePurchases(offset, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get vehicle purchases",
			Error:   err.Error(),
		})
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Vehicle purchases retrieved successfully",
		Data:    purchases,
		Meta: models.PaginationMeta{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// @Summary Get vehicle sales
// @Description Get list of vehicle sales with pagination
// @Tags Vehicle Trading
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.Response
// @Router /vehicle-trading/sales [get]
func (h *Handlers) GetVehicleSales(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)
	offset := (page - 1) * perPage

	sales, total, err := h.services.VehicleTrading.GetVehicleSales(offset, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get vehicle sales",
			Error:   err.Error(),
		})
	}

	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Vehicle sales retrieved successfully",
		Data:    sales,
		Meta: models.PaginationMeta{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// setupVehicleTradingRoutes sets up vehicle trading routes
func (h *Handlers) setupVehicleTradingRoutes(vt fiber.Router) {
	// Vehicle Purchase Routes
	purchases := vt.Group("/purchases")
	purchases.Get("/", h.requirePermission("vehicle_trading.read"), h.GetVehiclePurchases)
	purchases.Post("/", h.requirePermission("vehicle_trading.create"), h.CreateVehiclePurchase)
	purchases.Get("/:id", h.requirePermission("vehicle_trading.read"), h.GetVehiclePurchaseByID)
	purchases.Put("/:id", h.requirePermission("vehicle_trading.update"), h.UpdateVehiclePurchase)
	purchases.Delete("/:id", h.requirePermission("vehicle_trading.delete"), h.SoftDeleteVehiclePurchase)

	// Vehicle Inventory Routes
	inventory := vt.Group("/inventory")
	inventory.Get("/", h.requirePermission("vehicle_trading.read"), h.GetAvailableVehicles)
	inventory.Get("/:id", h.requirePermission("vehicle_trading.read"), h.GetVehicleInventoryByID)
	inventory.Put("/:id", h.requirePermission("vehicle_trading.update"), h.UpdateVehicleInventory)
	inventory.Put("/:id/price", h.requirePermission("vehicle_trading.update"), h.UpdateVehicleSellingPrice)
	inventory.Post("/:id/photos", h.requirePermission("vehicle_trading.update"), h.UploadVehiclePhotos)
	inventory.Delete("/:id", h.requirePermission("vehicle_trading.delete"), h.SoftDeleteVehicleInventory)

	// Vehicle Sales Routes
	sales := vt.Group("/sales")
	sales.Get("/", h.requirePermission("vehicle_trading.read"), h.GetVehicleSales)
	sales.Post("/", h.requirePermission("vehicle_trading.create"), h.CreateVehicleSale)
	sales.Get("/:id", h.requirePermission("vehicle_trading.read"), h.GetVehicleSaleByID)
	sales.Put("/:id", h.requirePermission("vehicle_trading.update"), h.UpdateVehicleSale)
	sales.Delete("/:id", h.requirePermission("vehicle_trading.delete"), h.SoftDeleteVehicleSale)

	// Analytics Routes
	analytics := vt.Group("/analytics")
	analytics.Get("/profit", h.requirePermission("vehicle_trading.read"), h.GetVehicleTradingProfitAnalysis)

	// Reports Routes
	reports := vt.Group("/reports")
	reports.Get("/inventory-aging", h.requirePermission("vehicle_trading.read"), h.GetInventoryAgingReport)
	reports.Get("/sales-performance", h.requirePermission("vehicle_trading.read"), h.GetSalesPerformanceReport)
	reports.Get("/commission", h.requirePermission("vehicle_trading.read"), h.GetCommissionReport)
}

// Additional placeholder handlers that need to be implemented
func (h *Handlers) GetVehiclePurchaseByID(c *fiber.Ctx) error {
	// Implementation similar to GetAvailableVehicles
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) UpdateVehiclePurchase(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) SoftDeleteVehiclePurchase(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) SoftDeleteVehicleInventory(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) GetVehicleSaleByID(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) UpdateVehicleSale(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) SoftDeleteVehicleSale(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}

func (h *Handlers) GetCommissionReport(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Not implemented"})
}
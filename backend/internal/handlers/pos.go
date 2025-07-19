package handlers

import (
	"strconv"
	"time"

	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
)

// POS Operations - Kasir-centric functionality

// CreatePOSTransaction creates a new POS transaction with multi-payment support
// POST /api/v1/pos/transactions
func (h *Handlers) CreatePOSTransaction(c *fiber.Ctx) error {
	var req models.CreatePOSTransactionRequest
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

	// Create the transaction
	transaction, err := h.services.Transaction.CreatePOSTransaction(&req, userID, outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to create transaction",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "POS transaction created successfully",
		Data:    transaction,
	})
}

// SearchProducts searches for products by query or barcode for POS
// GET /api/v1/pos/products/search
func (h *Handlers) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("query", "")
	barcode := c.Query("barcode", "")
	categoryStr := c.Query("category", "")
	limitStr := c.Query("limit", "20")

	if query == "" && barcode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Query or barcode is required",
		})
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var categoryID *int64
	if categoryStr != "" {
		if id, err := strconv.ParseInt(categoryStr, 10, 64); err == nil {
			categoryID = &id
		}
	}

	req := &models.ProductSearchRequest{
		Query:    query,
		Barcode:  barcode,
		Category: categoryID,
		Limit:    limit,
	}

	products, err := h.services.Product.SearchProducts(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to search products",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// AddTransactionPayment adds additional payment method to existing transaction
// PUT /api/v1/pos/transactions/:id/payment
func (h *Handlers) AddTransactionPayment(c *fiber.Ctx) error {
	transactionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid transaction ID",
		})
	}

	var req models.CreateTransactionPaymentRequest
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

	// Add payment to transaction
	payment, err := h.services.Transaction.AddTransactionPayment(transactionID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to add payment",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Payment added successfully",
		Data:    payment,
	})
}

// PrintReceipt generates receipt for transaction
// POST /api/v1/pos/transactions/:id/print
func (h *Handlers) PrintReceipt(c *fiber.Ctx) error {
	transactionID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid transaction ID",
		})
	}

	// Get transaction with details
	transaction, err := h.services.Transaction.GetTransactionWithDetails(transactionID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	}

	// Generate receipt data (simplified for now)
	receiptData := map[string]interface{}{
		"transaction":    transaction,
		"print_time":     time.Now(),
		"print_format":   "thermal", // Could be configurable
		"receipt_number": transaction.TransactionNumber,
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Receipt data generated successfully",
		Data:    receiptData,
	})
}

// GetQueueManagement gets current service queue for outlet
// GET /api/v1/pos/queue
func (h *Handlers) GetQueueManagement(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(int64)
	date := c.Query("date", time.Now().Format("2006-01-02"))

	req := &models.QueueManagementRequest{
		OutletID: outletID,
		Date:     date,
	}

	queueData, err := h.services.ServiceJob.GetQueueManagement(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get queue data",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Queue data retrieved successfully",
		Data:    queueData,
	})
}

// AutoAssignMechanic automatically assigns mechanic to service job based on workload
// PUT /api/v1/pos/service-jobs/:id/assign
func (h *Handlers) AutoAssignMechanic(c *fiber.Ctx) error {
	serviceJobID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	var req models.AutoAssignMechanicRequest
	req.ServiceJobID = serviceJobID
	
	if err := c.BodyParser(&req); err != nil {
		// If body parsing fails, just use the serviceJobID
		req.ServiceType = ""
	}

	// Auto-assign mechanic
	assignment, err := h.services.ServiceJob.AutoAssignMechanic(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to assign mechanic",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Mechanic assigned successfully",
		Data:    assignment,
	})
}

// Receivables Management

// GetPendingReceivables gets outstanding receivables (Belum Lunas)
// GET /api/v1/pos/receivables/pending
func (h *Handlers) GetPendingReceivables(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	receivables, total, err := h.services.Receivable.GetPendingReceivables(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get pending receivables",
			Error:   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Pending receivables retrieved successfully",
		Data:    receivables,
		Meta: models.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// RecordReceivablePayment records payment for receivable with kasir approval
// POST /api/v1/pos/receivables/:id/payment
func (h *Handlers) RecordReceivablePayment(c *fiber.Ctx) error {
	receivableID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid receivable ID",
		})
	}

	var req models.ApproveReceivablePaymentRequest
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

	// Get kasir user ID from JWT context
	kasirID := c.Locals("user_id").(int64)

	// Record payment with kasir approval
	payment, err := h.services.Receivable.RecordPaymentWithApproval(receivableID, &req, kasirID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to record payment",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Payment recorded and approved successfully",
		Data:    payment,
	})
}

// GetPaidReceivables gets paid receivables (Lunas)
// GET /api/v1/pos/receivables/paid
func (h *Handlers) GetPaidReceivables(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	receivables, total, err := h.services.Receivable.GetPaidReceivables(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get paid receivables",
			Error:   err.Error(),
		})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Paid receivables retrieved successfully",
		Data:    receivables,
		Meta: models.PaginationMeta{
			CurrentPage: page,
			PerPage:     limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	})
}

// Dashboard helpers for kasir

// GetDashboardStats gets dashboard statistics for kasir
// GET /api/v1/pos/dashboard/stats
func (h *Handlers) GetDashboardStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	outletID := c.Locals("outlet_id").(int64)
	date := c.Query("date", time.Now().Format("2006-01-02"))

	stats, err := h.services.Dashboard.GetKasirDashboardStats(userID, outletID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get dashboard statistics",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Dashboard statistics retrieved successfully",
		Data:    stats,
	})
}
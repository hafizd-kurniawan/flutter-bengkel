package handlers

import (
	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
)

// setupServiceRoutes sets up service management routes
func (h *Handlers) setupServiceRoutes(services fiber.Router) {
	services.Get("/", h.requirePermission("services.read"), h.getServices)
	services.Get("/:id", h.requirePermission("services.read"), h.getServiceByID)
	services.Post("/", h.requirePermission("services.create"), h.createService)
	services.Put("/:id", h.requirePermission("services.update"), h.updateService)
	services.Delete("/:id", h.requirePermission("services.delete"), h.deleteService)
}

// setupProductRoutes sets up product management routes
func (h *Handlers) setupProductRoutes(products fiber.Router) {
	products.Get("/", h.requirePermission("products.read"), h.getProducts)
	products.Get("/:id", h.requirePermission("products.read"), h.getProductByID)
	products.Post("/", h.requirePermission("products.create"), h.createProduct)
	products.Put("/:id", h.requirePermission("products.update"), h.updateProduct)
	products.Delete("/:id", h.requirePermission("products.delete"), h.deleteProduct)
	products.Put("/:id/stock", h.requirePermission("products.update"), h.updateProductStock)
	products.Get("/low-stock", h.requirePermission("products.read"), h.getLowStockProducts)
}

// setupTransactionRoutes sets up transaction management routes
func (h *Handlers) setupTransactionRoutes(transactions fiber.Router) {
	transactions.Get("/", h.requirePermission("transactions.read"), h.getTransactions)
	transactions.Get("/:id", h.requirePermission("transactions.read"), h.getTransactionByID)
	transactions.Post("/", h.requirePermission("transactions.create"), h.createTransaction)
	transactions.Put("/:id", h.requirePermission("transactions.update"), h.updateTransaction)
	transactions.Delete("/:id", h.requirePermission("transactions.delete"), h.deleteTransaction)
}

// setupPaymentRoutes sets up payment management routes
func (h *Handlers) setupPaymentRoutes(payments fiber.Router) {
	payments.Get("/", h.requirePermission("transactions.read"), h.getPayments)
	payments.Get("/:id", h.requirePermission("transactions.read"), h.getPaymentByID)
	payments.Post("/", h.requirePermission("transactions.create"), h.createPayment)
	payments.Delete("/:id", h.requirePermission("transactions.delete"), h.deletePayment)
}

// setupMasterDataRoutes sets up master data routes
func (h *Handlers) setupMasterDataRoutes(masterData fiber.Router) {
	masterData.Get("/service-categories", h.getServiceCategories)
	masterData.Get("/product-categories", h.getProductCategories)
	masterData.Get("/suppliers", h.getSuppliers)
	masterData.Get("/unit-types", h.getUnitTypes)
	masterData.Get("/payment-methods", h.getPaymentMethods)
}

// Service handlers
func (h *Handlers) getServices(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	var categoryID *int64
	if categoryIDParam := c.QueryInt("category_id", 0); categoryIDParam > 0 {
		id := int64(categoryIDParam)
		categoryID = &id
	}

	services, meta, err := h.services.Service.List(page, limit, categoryID, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get services",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Services retrieved successfully",
		Data:    services,
		Meta:    *meta,
	})
}

func (h *Handlers) getServiceByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service ID",
		})
	}

	service, err := h.services.Service.GetByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Service not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service retrieved successfully",
		Data:    service,
	})
}

func (h *Handlers) createService(c *fiber.Ctx) error {
	var req models.Service
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	service, err := h.services.Service.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to create service",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Service created successfully",
		Data:    service,
	})
}

func (h *Handlers) updateService(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service ID",
		})
	}

	var req models.Service
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	service, err := h.services.Service.Update(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update service",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service updated successfully",
		Data:    service,
	})
}

func (h *Handlers) deleteService(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service ID",
		})
	}

	if err := h.services.Service.Delete(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete service",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service deleted successfully",
	})
}

// Product handlers
func (h *Handlers) getProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	var categoryID *int64
	if categoryIDParam := c.QueryInt("category_id", 0); categoryIDParam > 0 {
		id := int64(categoryIDParam)
		categoryID = &id
	}

	var supplierID *int64
	if supplierIDParam := c.QueryInt("supplier_id", 0); supplierIDParam > 0 {
		id := int64(supplierIDParam)
		supplierID = &id
	}

	products, meta, err := h.services.Product.List(page, limit, categoryID, supplierID, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get products",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
		Meta:    *meta,
	})
}

func (h *Handlers) getProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	product, err := h.services.Product.GetByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

func (h *Handlers) createProduct(c *fiber.Ctx) error {
	var req models.Product
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	product, err := h.services.Product.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to create product",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

func (h *Handlers) updateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var req models.Product
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	product, err := h.services.Product.Update(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update product",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

func (h *Handlers) deleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	if err := h.services.Product.Delete(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete product",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Product deleted successfully",
	})
}

func (h *Handlers) updateProductStock(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid product ID",
		})
	}

	var req struct {
		Quantity  int    `json:"quantity"`
		Operation string `json:"operation"` // "add" or "subtract"
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.services.Product.UpdateStock(int64(id), req.Quantity, req.Operation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update product stock",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Product stock updated successfully",
	})
}

func (h *Handlers) getLowStockProducts(c *fiber.Ctx) error {
	var outletID *int64
	if outletIDParam := c.QueryInt("outlet_id", 0); outletIDParam > 0 {
		id := int64(outletIDParam)
		outletID = &id
	}

	products, err := h.services.Product.GetLowStockProducts(outletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get low stock products",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Low stock products retrieved successfully",
		Data:    products,
	})
}

// Transaction handlers (simplified)
func (h *Handlers) getTransactions(c *fiber.Ctx) error {
	// Implementation similar to other list handlers
	return c.JSON(models.Response{Success: true, Message: "Transactions endpoint"})
}

func (h *Handlers) getTransactionByID(c *fiber.Ctx) error {
	// Implementation similar to other get by ID handlers
	return c.JSON(models.Response{Success: true, Message: "Transaction by ID endpoint"})
}

func (h *Handlers) createTransaction(c *fiber.Ctx) error {
	// Implementation for transaction creation
	return c.JSON(models.Response{Success: true, Message: "Create transaction endpoint"})
}

func (h *Handlers) updateTransaction(c *fiber.Ctx) error {
	// Implementation for transaction update
	return c.JSON(models.Response{Success: true, Message: "Update transaction endpoint"})
}

func (h *Handlers) deleteTransaction(c *fiber.Ctx) error {
	// Implementation for transaction deletion
	return c.JSON(models.Response{Success: true, Message: "Delete transaction endpoint"})
}

// Payment handlers (simplified)
func (h *Handlers) getPayments(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Payments endpoint"})
}

func (h *Handlers) getPaymentByID(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Payment by ID endpoint"})
}

func (h *Handlers) createPayment(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Create payment endpoint"})
}

func (h *Handlers) deletePayment(c *fiber.Ctx) error {
	return c.JSON(models.Response{Success: true, Message: "Delete payment endpoint"})
}

// Master data handlers
func (h *Handlers) getServiceCategories(c *fiber.Ctx) error {
	categories, err := h.services.Service.ListCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get service categories",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service categories retrieved successfully",
		Data:    categories,
	})
}

func (h *Handlers) getProductCategories(c *fiber.Ctx) error {
	categories, err := h.services.Product.ListCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get product categories",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Product categories retrieved successfully",
		Data:    categories,
	})
}

func (h *Handlers) getSuppliers(c *fiber.Ctx) error {
	suppliers, err := h.services.Product.ListSuppliers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get suppliers",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Suppliers retrieved successfully",
		Data:    suppliers,
	})
}

func (h *Handlers) getUnitTypes(c *fiber.Ctx) error {
	unitTypes, err := h.services.Product.ListUnitTypes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get unit types",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Unit types retrieved successfully",
		Data:    unitTypes,
	})
}

func (h *Handlers) getPaymentMethods(c *fiber.Ctx) error {
	paymentMethods, err := h.services.Payment.ListPaymentMethods()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get payment methods",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Payment methods retrieved successfully",
		Data:    paymentMethods,
	})
}
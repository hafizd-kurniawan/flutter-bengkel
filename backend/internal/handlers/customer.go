package handlers

import (
	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
)

// setupCustomerRoutes sets up customer management routes
func (h *Handlers) setupCustomerRoutes(customers fiber.Router) {
	customers.Get("/", h.requirePermission("customers.read"), h.getCustomers)
	customers.Get("/:id", h.requirePermission("customers.read"), h.getCustomerByID)
	customers.Post("/", h.requirePermission("customers.create"), h.createCustomer)
	customers.Put("/:id", h.requirePermission("customers.update"), h.updateCustomer)
	customers.Delete("/:id", h.requirePermission("customers.delete"), h.deleteCustomer)
}

// @Summary Get customers
// @Description Get paginated list of customers
// @Tags Customers
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse{data=[]models.Customer}
// @Router /customers [get]
func (h *Handlers) getCustomers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	customers, meta, err := h.services.Customer.List(page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get customers",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Customers retrieved successfully",
		Data:    customers,
		Meta:    *meta,
	})
}

// @Summary Get customer by ID
// @Description Get customer details by ID including vehicles
// @Tags Customers
// @Security Bearer
// @Param id path int true "Customer ID"
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /customers/{id} [get]
func (h *Handlers) getCustomerByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid customer ID",
		})
	}

	customer, err := h.services.Customer.GetByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Customer not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Customer retrieved successfully",
		Data:    customer,
	})
}

// @Summary Create customer
// @Description Create a new customer
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body models.Customer true "Customer data"
// @Success 201 {object} models.Response{data=models.Customer}
// @Router /customers [post]
func (h *Handlers) createCustomer(c *fiber.Ctx) error {
	var req models.Customer
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	customer, err := h.services.Customer.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to create customer",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Customer created successfully",
		Data:    customer,
	})
}

// @Summary Update customer
// @Description Update customer details
// @Tags Customers
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param request body models.Customer true "Customer data"
// @Success 200 {object} models.Response{data=models.Customer}
// @Router /customers/{id} [put]
func (h *Handlers) updateCustomer(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid customer ID",
		})
	}

	var req models.Customer
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	customer, err := h.services.Customer.Update(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update customer",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Customer updated successfully",
		Data:    customer,
	})
}

// @Summary Delete customer
// @Description Delete customer by ID
// @Tags Customers
// @Security Bearer
// @Param id path int true "Customer ID"
// @Success 200 {object} models.Response
// @Router /customers/{id} [delete]
func (h *Handlers) deleteCustomer(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid customer ID",
		})
	}

	if err := h.services.Customer.Delete(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete customer",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Customer deleted successfully",
	})
}

// setupVehicleRoutes sets up vehicle management routes
func (h *Handlers) setupVehicleRoutes(vehicles fiber.Router) {
	vehicles.Get("/", h.requirePermission("vehicles.read"), h.getVehicles)
	vehicles.Get("/:id", h.requirePermission("vehicles.read"), h.getVehicleByID)
	vehicles.Post("/", h.requirePermission("vehicles.create"), h.createVehicle)
	vehicles.Put("/:id", h.requirePermission("vehicles.update"), h.updateVehicle)
	vehicles.Delete("/:id", h.requirePermission("vehicles.delete"), h.deleteVehicle)
}

// @Summary Get vehicles
// @Description Get paginated list of vehicles
// @Tags Vehicles
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param customer_id query int false "Filter by customer ID"
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse{data=[]models.CustomerVehicle}
// @Router /vehicles [get]
func (h *Handlers) getVehicles(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	var customerID *int64
	if customerIDParam := c.QueryInt("customer_id", 0); customerIDParam > 0 {
		id := int64(customerIDParam)
		customerID = &id
	}

	vehicles, meta, err := h.services.Vehicle.List(page, limit, customerID, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get vehicles",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Vehicles retrieved successfully",
		Data:    vehicles,
		Meta:    *meta,
	})
}

// @Summary Get vehicle by ID
// @Description Get vehicle details by ID
// @Tags Vehicles
// @Security Bearer
// @Param id path int true "Vehicle ID"
// @Success 200 {object} models.Response{data=models.CustomerVehicle}
// @Router /vehicles/{id} [get]
func (h *Handlers) getVehicleByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle ID",
		})
	}

	vehicle, err := h.services.Vehicle.GetByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Vehicle not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle retrieved successfully",
		Data:    vehicle,
	})
}

// @Summary Create vehicle
// @Description Create a new vehicle
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body models.CustomerVehicle true "Vehicle data"
// @Success 201 {object} models.Response{data=models.CustomerVehicle}
// @Router /vehicles [post]
func (h *Handlers) createVehicle(c *fiber.Ctx) error {
	var req models.CustomerVehicle
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	vehicle, err := h.services.Vehicle.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to create vehicle",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Vehicle created successfully",
		Data:    vehicle,
	})
}

// @Summary Update vehicle
// @Description Update vehicle details
// @Tags Vehicles
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Param request body models.CustomerVehicle true "Vehicle data"
// @Success 200 {object} models.Response{data=models.CustomerVehicle}
// @Router /vehicles/{id} [put]
func (h *Handlers) updateVehicle(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle ID",
		})
	}

	var req models.CustomerVehicle
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	vehicle, err := h.services.Vehicle.Update(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update vehicle",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle updated successfully",
		Data:    vehicle,
	})
}

// @Summary Delete vehicle
// @Description Delete vehicle by ID
// @Tags Vehicles
// @Security Bearer
// @Param id path int true "Vehicle ID"
// @Success 200 {object} models.Response
// @Router /vehicles/{id} [delete]
func (h *Handlers) deleteVehicle(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid vehicle ID",
		})
	}

	if err := h.services.Vehicle.Delete(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete vehicle",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Vehicle deleted successfully",
	})
}
package handlers

import (
	"flutter-bengkel/internal/middleware"
	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
)

// setupServiceJobRoutes sets up service job management routes
func (h *Handlers) setupServiceJobRoutes(serviceJobs fiber.Router) {
	serviceJobs.Get("/", h.requirePermission("service_jobs.read"), h.getServiceJobs)
	serviceJobs.Get("/:id", h.requirePermission("service_jobs.read"), h.getServiceJobByID)
	serviceJobs.Post("/", h.requirePermission("service_jobs.create"), h.createServiceJob)
	serviceJobs.Put("/:id", h.requirePermission("service_jobs.update"), h.updateServiceJob)
	serviceJobs.Put("/:id/status", h.requirePermission("service_jobs.update"), h.updateServiceJobStatus)
	serviceJobs.Delete("/:id", h.requirePermission("service_jobs.delete"), h.deleteServiceJob)
	
	// Service job details
	serviceJobs.Get("/:id/details", h.requirePermission("service_jobs.read"), h.getServiceJobDetails)
	serviceJobs.Post("/:id/details", h.requirePermission("service_jobs.update"), h.addServiceJobDetail)
	serviceJobs.Put("/details/:detail_id", h.requirePermission("service_jobs.update"), h.updateServiceJobDetail)
	serviceJobs.Delete("/details/:detail_id", h.requirePermission("service_jobs.update"), h.deleteServiceJobDetail)
}

// @Summary Get service jobs
// @Description Get paginated list of service jobs
// @Tags Service Jobs
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param outlet_id query int false "Filter by outlet ID"
// @Param status query string false "Filter by status"
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse{data=[]models.ServiceJob}
// @Router /service-jobs [get]
func (h *Handlers) getServiceJobs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status", "")
	search := c.Query("search", "")

	// Get user's outlet ID from context
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Message: "Unauthorized",
		})
	}

	var outletID *int64
	if claims.OutletID != nil {
		outletID = claims.OutletID
	}

	// Allow overriding outlet ID for super admin
	if claims.RoleID == 1 { // Super Admin
		if outletIDParam := c.QueryInt("outlet_id", 0); outletIDParam > 0 {
			id := int64(outletIDParam)
			outletID = &id
		}
	}

	serviceJobs, meta, err := h.services.ServiceJob.List(page, limit, outletID, status, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get service jobs",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Service jobs retrieved successfully",
		Data:    serviceJobs,
		Meta:    *meta,
	})
}

// @Summary Get service job by ID
// @Description Get service job details by ID including details
// @Tags Service Jobs
// @Security Bearer
// @Param id path int true "Service Job ID"
// @Success 200 {object} models.Response{data=models.ServiceJob}
// @Router /service-jobs/{id} [get]
func (h *Handlers) getServiceJobByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	serviceJob, err := h.services.ServiceJob.GetByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Service job not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job retrieved successfully",
		Data:    serviceJob,
	})
}

// @Summary Create service job
// @Description Create a new service job
// @Tags Service Jobs
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body models.CreateServiceJobRequest true "Service job data"
// @Success 201 {object} models.Response{data=models.ServiceJob}
// @Router /service-jobs [post]
func (h *Handlers) createServiceJob(c *fiber.Ctx) error {
	var req models.CreateServiceJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Get user and outlet from context
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Message: "Unauthorized",
		})
	}

	if claims.OutletID == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "User must be assigned to an outlet",
		})
	}

	serviceJob, err := h.services.ServiceJob.Create(&req, *claims.OutletID, claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to create service job",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Service job created successfully",
		Data:    serviceJob,
	})
}

// @Summary Update service job
// @Description Update service job details
// @Tags Service Jobs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Job ID"
// @Param request body models.UpdateServiceJobRequest true "Service job data"
// @Success 200 {object} models.Response{data=models.ServiceJob}
// @Router /service-jobs/{id} [put]
func (h *Handlers) updateServiceJob(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	var req models.UpdateServiceJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	serviceJob, err := h.services.ServiceJob.Update(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update service job",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job updated successfully",
		Data:    serviceJob,
	})
}

// @Summary Update service job status
// @Description Update service job status with history tracking
// @Tags Service Jobs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Job ID"
// @Param request body object{status=string,notes=string} true "Status update data"
// @Success 200 {object} models.Response
// @Router /service-jobs/{id}/status [put]
func (h *Handlers) updateServiceJobStatus(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Get user from context
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Message: "Unauthorized",
		})
	}

	if err := h.services.ServiceJob.UpdateStatus(int64(id), req.Status, claims.UserID, req.Notes); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update service job status",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job status updated successfully",
	})
}

// @Summary Delete service job
// @Description Delete service job by ID (sets status to cancelled)
// @Tags Service Jobs
// @Security Bearer
// @Param id path int true "Service Job ID"
// @Success 200 {object} models.Response
// @Router /service-jobs/{id} [delete]
func (h *Handlers) deleteServiceJob(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	if err := h.services.ServiceJob.Delete(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete service job",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job deleted successfully",
	})
}

// Service job details handlers

// @Summary Get service job details
// @Description Get all details for a service job
// @Tags Service Jobs
// @Security Bearer
// @Param id path int true "Service Job ID"
// @Success 200 {object} models.Response{data=[]models.ServiceDetail}
// @Router /service-jobs/{id}/details [get]
func (h *Handlers) getServiceJobDetails(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	details, err := h.services.ServiceJob.GetDetails(int64(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get service job details",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job details retrieved successfully",
		Data:    details,
	})
}

// @Summary Add service job detail
// @Description Add a new detail to a service job
// @Tags Service Jobs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Job ID"
// @Param request body models.ServiceDetail true "Service detail data"
// @Success 201 {object} models.Response{data=models.ServiceDetail}
// @Router /service-jobs/{id}/details [post]
func (h *Handlers) addServiceJobDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid service job ID",
		})
	}

	var req models.ServiceDetail
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	detail, err := h.services.ServiceJob.AddDetail(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to add service job detail",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "Service job detail added successfully",
		Data:    detail,
	})
}

// @Summary Update service job detail
// @Description Update a service job detail
// @Tags Service Jobs
// @Security Bearer
// @Accept json
// @Produce json
// @Param detail_id path int true "Detail ID"
// @Param request body models.ServiceDetail true "Service detail data"
// @Success 200 {object} models.Response
// @Router /service-jobs/details/{detail_id} [put]
func (h *Handlers) updateServiceJobDetail(c *fiber.Ctx) error {
	detailID, err := c.ParamsInt("detail_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid detail ID",
		})
	}

	var req models.ServiceDetail
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.services.ServiceJob.UpdateDetail(int64(detailID), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update service job detail",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job detail updated successfully",
	})
}

// @Summary Delete service job detail
// @Description Delete a service job detail
// @Tags Service Jobs
// @Security Bearer
// @Param detail_id path int true "Detail ID"
// @Success 200 {object} models.Response
// @Router /service-jobs/details/{detail_id} [delete]
func (h *Handlers) deleteServiceJobDetail(c *fiber.Ctx) error {
	detailID, err := c.ParamsInt("detail_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid detail ID",
		})
	}

	if err := h.services.ServiceJob.DeleteDetail(int64(detailID)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete service job detail",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Service job detail deleted successfully",
	})
}
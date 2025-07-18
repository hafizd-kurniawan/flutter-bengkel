package handlers

import (
	"flutter-bengkel/internal/middleware"
	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
)

// jwtMiddleware returns the JWT middleware
func (h *Handlers) jwtMiddleware() fiber.Handler {
	return middleware.JWTMiddleware(h.config.JWT.Secret)
}

// requirePermission returns middleware that checks for specific permission
func (h *Handlers) requirePermission(permission string) fiber.Handler {
	return middleware.RequirePermission(permission)
}

// setupAuthRoutes sets up authentication routes
func (h *Handlers) setupAuthRoutes(auth fiber.Router) {
	auth.Post("/login", h.login)
	auth.Post("/refresh", h.refreshToken)
	auth.Post("/logout", h.logout)
}

// @Summary User login
// @Description Authenticate user with username and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.Response{data=models.LoginResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /auth/login [post]
func (h *Handlers) login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	response, err := h.services.Auth.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.Response{data=models.LoginResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /auth/refresh [post]
func (h *Handlers) refreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	response, err := h.services.Auth.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Success: false,
			Message: "Token refresh failed",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// @Summary User logout
// @Description Logout user (client should discard tokens)
// @Tags Authentication
// @Security Bearer
// @Success 200 {object} models.Response
// @Router /auth/logout [post]
func (h *Handlers) logout(c *fiber.Ctx) error {
	// For now, logout is handled on the client side by discarding tokens
	// In the future, we could implement a token blacklist
	return c.JSON(models.Response{
		Success: true,
		Message: "Logged out successfully",
	})
}

// setupUserRoutes sets up user management routes
func (h *Handlers) setupUserRoutes(users fiber.Router) {
	users.Get("/", h.requirePermission("users.read"), h.getUsers)
	users.Get("/:id", h.requirePermission("users.read"), h.getUserByID)
	users.Post("/", h.requirePermission("users.create"), h.createUser)
	users.Put("/:id", h.requirePermission("users.update"), h.updateUser)
	users.Delete("/:id", h.requirePermission("users.delete"), h.deleteUser)
	users.Post("/:id/change-password", h.requirePermission("users.update"), h.changePassword)
}

// @Summary Get users
// @Description Get paginated list of users
// @Tags Users
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse{data=[]models.User}
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Router /users [get]
func (h *Handlers) getUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	users, meta, err := h.services.User.List(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Success: false,
			Message: "Failed to get users",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
		Meta:    *meta,
	})
}

// @Summary Get user by ID
// @Description Get user details by ID
// @Tags Users
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} models.Response{data=models.User}
// @Failure 404 {object} models.Response
// @Router /users/{id} [get]
func (h *Handlers) getUserByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	user, err := h.services.User.GetByID(int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// @Summary Create user
// @Description Create a new user
// @Tags Users
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Create user request"
// @Success 201 {object} models.Response{data=models.User}
// @Failure 400 {object} models.Response
// @Router /users [post]
func (h *Handlers) createUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	user, err := h.services.User.Create(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// @Summary Update user
// @Description Update user details
// @Tags Users
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "Update user request"
// @Success 200 {object} models.Response{data=models.User}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /users/{id} [put]
func (h *Handlers) updateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	user, err := h.services.User.Update(int64(id), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /users/{id} [delete]
func (h *Handlers) deleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	if err := h.services.User.Delete(int64(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.Response{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// @Summary Change password
// @Description Change user password
// @Tags Users
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /users/{id}/change-password [post]
func (h *Handlers) changePassword(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.services.User.ChangePassword(int64(id), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Failed to change password",
			Error:   err.Error(),
		})
	}

	return c.JSON(models.Response{
		Success: true,
		Message: "Password changed successfully",
	})
}
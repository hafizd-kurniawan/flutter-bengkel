package middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrorHandler handles all errors in the application
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Check for validation errors
	if strings.Contains(err.Error(), "validation") {
		code = fiber.StatusBadRequest
		message = "Validation failed"
	}

	// Check for JWT errors
	if strings.Contains(err.Error(), "token") {
		code = fiber.StatusUnauthorized
		message = "Invalid or expired token"
	}

	// Log the error
	log.Printf("Error: %v", err)

	return c.Status(code).JSON(models.Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}

// JWTMiddleware validates JWT tokens
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Authorization header is required",
			})
		}

		// Check Bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Invalid authorization header format",
			})
		}

		tokenString := tokenParts[1]

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Invalid token",
				Error:   err.Error(),
			})
		}

		// Validate token
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Invalid token claims",
			})
		}

		// Store claims in context
		userIDStr := (*claims)["user_id"].(string)
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Invalid user ID format",
			})
		}
		c.Locals("user_id", userID)
		c.Locals("username", (*claims)["username"].(string))
		
		roleIDStr := (*claims)["role_id"].(string)
		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				Success: false,
				Message: "Invalid role ID format",
			})
		}
		c.Locals("role_id", roleID)
		
		if outletIDRaw, exists := (*claims)["outlet_id"]; exists && outletIDRaw != nil {
			outletIDStr := outletIDRaw.(string)
			outletID, err := uuid.Parse(outletIDStr)
			if err == nil {
				c.Locals("outlet_id", outletID)
			}
		}

		if permissions, exists := (*claims)["permissions"]; exists {
			c.Locals("permissions", permissions)
		}

		return c.Next()
	}
}

// RequirePermission checks if user has required permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.Response{
				Success: false,
				Message: "No permissions found",
			})
		}

		// Check if user has the required permission
		hasPermission := false
		for _, p := range permissions {
			if p.(string) == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(models.Response{
				Success: false,
				Message: "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireRole checks if user has required role
func RequireRole(roleID int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleID, ok := c.Locals("role_id").(int64)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(models.Response{
				Success: false,
				Message: "No role found",
			})
		}

		if userRoleID != roleID {
			return c.Status(fiber.StatusForbidden).JSON(models.Response{
				Success: false,
				Message: "Insufficient role permissions",
			})
		}

		return c.Next()
	}
}

// GetUserFromContext extracts user information from context
func GetUserFromContext(c *fiber.Ctx) (*models.Claims, error) {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return nil, errors.New("username not found in context")
	}

	roleID, ok := c.Locals("role_id").(uuid.UUID)
	if !ok {
		return nil, errors.New("role ID not found in context")
	}

	claims := &models.Claims{
		UserID:   userID,
		Username: username,
		RoleID:   roleID,
	}

	// Outlet ID is optional
	if outletID, ok := c.Locals("outlet_id").(uuid.UUID); ok {
		claims.OutletID = &outletID
	}

	// Permissions are optional
	if permissions, ok := c.Locals("permissions").([]interface{}); ok {
		permissionStrings := make([]string, len(permissions))
		for i, p := range permissions {
			permissionStrings[i] = p.(string)
		}
		claims.Permissions = permissionStrings
	}

	return claims, nil
}

// ValidateRequestBody validates request body using struct tags
func ValidateRequestBody(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// You can add more validation logic here using go-playground/validator
	// For now, we'll just parse the body

	return nil
}
package main

import (
	"log"
	"time"

	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/database"
	"flutter-bengkel/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Simple outlet repository for demonstration
type OutletRepo struct {
	db *sqlx.DB
}

func NewOutletRepo(db *sqlx.DB) *OutletRepo {
	return &OutletRepo{db: db}
}

type Outlet struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Address   string     `json:"address" db:"address"`
	Phone     string     `json:"phone" db:"phone"`
	Email     string     `json:"email" db:"email"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty" db:"deleted_by"`
}

func (r *OutletRepo) GetAll(includeDeleted bool) ([]Outlet, error) {
	query := "SELECT id, name, address, phone, email, is_active, created_at, updated_at, deleted_at, deleted_by FROM outlets"
	if !includeDeleted {
		query += " WHERE deleted_at IS NULL"
	}
	query += " ORDER BY created_at DESC"

	var outlets []Outlet
	err := r.db.Select(&outlets, query)
	return outlets, err
}

func (r *OutletRepo) GetByID(id uuid.UUID, includeDeleted bool) (*Outlet, error) {
	query := "SELECT id, name, address, phone, email, is_active, created_at, updated_at, deleted_at, deleted_by FROM outlets WHERE id = $1"
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}

	var outlet Outlet
	err := r.db.Get(&outlet, query, id)
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func (r *OutletRepo) Create(outlet *Outlet) error {
	outlet.ID = uuid.New()
	query := `
		INSERT INTO outlets (id, name, address, phone, email, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	_, err := r.db.Exec(query, outlet.ID, outlet.Name, outlet.Address, outlet.Phone, outlet.Email, outlet.IsActive)
	return err
}

func (r *OutletRepo) SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error {
	query := "UPDATE outlets SET deleted_at = NOW(), deleted_by = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, deletedBy, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Outlet not found or already deleted")
	}

	return nil
}

func (r *OutletRepo) Restore(id uuid.UUID) error {
	query := "UPDATE outlets SET deleted_at = NULL, deleted_by = NULL, updated_at = NOW() WHERE id = $1 AND deleted_at IS NOT NULL"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Outlet not found or not deleted")
	}

	return nil
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	outletRepo := NewOutletRepo(db.GetDB())

	// Create Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	api := app.Group("/api/v1")

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "PostgreSQL Workshop Management System API is running",
			"database": "PostgreSQL with UUID and Soft Delete",
		})
	})

	// Outlets routes with soft delete demonstration
	outlets := api.Group("/outlets")

	// Get all outlets (active only by default)
	outlets.Get("/", func(c *fiber.Ctx) error {
		includeDeleted := c.Query("include_deleted") == "true"
		data, err := outletRepo.GetAll(includeDeleted)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Message: "Failed to get outlets",
				Error:   err.Error(),
			})
		}

		return c.JSON(models.Response{
			Success: true,
			Message: "Outlets retrieved successfully",
			Data:    data,
		})
	})

	// Get outlet by ID
	outlets.Get("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(models.Response{
				Success: false,
				Message: "Invalid UUID format",
			})
		}

		includeDeleted := c.Query("include_deleted") == "true"
		data, err := outletRepo.GetByID(id, includeDeleted)
		if err != nil {
			return c.Status(404).JSON(models.Response{
				Success: false,
				Message: "Outlet not found",
			})
		}

		return c.JSON(models.Response{
			Success: true,
			Message: "Outlet retrieved successfully",
			Data:    data,
		})
	})

	// Create outlet
	outlets.Post("/", func(c *fiber.Ctx) error {
		var outlet Outlet
		if err := c.BodyParser(&outlet); err != nil {
			return c.Status(400).JSON(models.Response{
				Success: false,
				Message: "Invalid request body",
			})
		}

		outlet.IsActive = true // Default to active
		err := outletRepo.Create(&outlet)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				Success: false,
				Message: "Failed to create outlet",
				Error:   err.Error(),
			})
		}

		return c.Status(201).JSON(models.Response{
			Success: true,
			Message: "Outlet created successfully",
			Data:    outlet,
		})
	})

	// Soft delete outlet
	outlets.Delete("/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(models.Response{
				Success: false,
				Message: "Invalid UUID format",
			})
		}

		// For demo purposes, use a dummy user ID
		deletedBy := uuid.New()
		err = outletRepo.SoftDelete(id, deletedBy)
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(models.Response{
					Success: false,
					Message: fiberErr.Message,
				})
			}
			return c.Status(500).JSON(models.Response{
				Success: false,
				Message: "Failed to delete outlet",
				Error:   err.Error(),
			})
		}

		return c.JSON(models.Response{
			Success: true,
			Message: "Outlet soft deleted successfully",
		})
	})

	// Restore outlet
	outlets.Post("/:id/restore", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(models.Response{
				Success: false,
				Message: "Invalid UUID format",
			})
		}

		err = outletRepo.Restore(id)
		if err != nil {
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(models.Response{
					Success: false,
					Message: fiberErr.Message,
				})
			}
			return c.Status(500).JSON(models.Response{
				Success: false,
				Message: "Failed to restore outlet",
				Error:   err.Error(),
			})
		}

		return c.JSON(models.Response{
			Success: true,
			Message: "Outlet restored successfully",
		})
	})

	// Start server
	log.Printf("🚀 PostgreSQL API Server started on port %s", cfg.Server.Port)
	log.Printf("📚 Available endpoints:")
	log.Printf("   GET /health - Health check")
	log.Printf("   GET /api/v1/outlets - Get all outlets")
	log.Printf("   GET /api/v1/outlets?include_deleted=true - Get all outlets including deleted")
	log.Printf("   GET /api/v1/outlets/:id - Get outlet by ID")
	log.Printf("   POST /api/v1/outlets - Create new outlet")
	log.Printf("   DELETE /api/v1/outlets/:id - Soft delete outlet")
	log.Printf("   POST /api/v1/outlets/:id/restore - Restore soft deleted outlet")

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
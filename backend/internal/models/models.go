package models

import (
	"time"
	"github.com/google/uuid"
)

// Common response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Pagination structure
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// Base model with common fields including soft delete
type BaseModel struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty" db:"deleted_by"`
}

// User model
type User struct {
	BaseModel
	Username     string     `json:"username" db:"username" validate:"required,min=3,max=100"`
	Email        string     `json:"email" db:"email" validate:"required,email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	FullName     string     `json:"full_name" db:"full_name" validate:"required"`
	Phone        string     `json:"phone" db:"phone"`
	RoleID       uuid.UUID  `json:"role_id" db:"role_id" validate:"required"`
	OutletID     *uuid.UUID `json:"outlet_id" db:"outlet_id"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	
	// Relations
	Role   *Role   `json:"role,omitempty"`
	Outlet *Outlet `json:"outlet,omitempty"`
}

// CreateUserRequest for user creation
type CreateUserRequest struct {
	Username string     `json:"username" validate:"required,min=3,max=100"`
	Email    string     `json:"email" validate:"required,email"`
	Password string     `json:"password" validate:"required,min=6"`
	FullName string     `json:"full_name" validate:"required"`
	Phone    string     `json:"phone"`
	RoleID   uuid.UUID  `json:"role_id" validate:"required"`
	OutletID *uuid.UUID `json:"outlet_id"`
}

// UpdateUserRequest for user updates
type UpdateUserRequest struct {
	Username string     `json:"username" validate:"required,min=3,max=100"`
	Email    string     `json:"email" validate:"required,email"`
	FullName string     `json:"full_name" validate:"required"`
	Phone    string     `json:"phone"`
	RoleID   uuid.UUID  `json:"role_id" validate:"required"`
	OutletID *uuid.UUID `json:"outlet_id"`
	IsActive bool       `json:"is_active"`
}

// ChangePasswordRequest for password changes
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// Role model
type Role struct {
	BaseModel
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description"`
	
	// Relations
	Permissions []Permission `json:"permissions,omitempty"`
}

// Permission model
type Permission struct {
	BaseModel
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description"`
	Resource    string `json:"resource" db:"resource" validate:"required"`
	Action      string `json:"action" db:"action" validate:"required"`
}

// Outlet model
type Outlet struct {
	BaseModel
	Name     string `json:"name" db:"name" validate:"required"`
	Address  string `json:"address" db:"address"`
	Phone    string `json:"phone" db:"phone"`
	Email    string `json:"email" db:"email"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

// Customer model
type Customer struct {
	BaseModel
	CustomerCode   string  `json:"customer_code" db:"customer_code"`
	Name           string  `json:"name" db:"name" validate:"required"`
	Email          string  `json:"email" db:"email"`
	Phone          string  `json:"phone" db:"phone" validate:"required"`
	Address        string  `json:"address" db:"address"`
	City           string  `json:"city" db:"city"`
	Province       string  `json:"province" db:"province"`
	PostalCode     string  `json:"postal_code" db:"postal_code"`
	DateOfBirth    *string `json:"date_of_birth" db:"date_of_birth"`
	Gender         string  `json:"gender" db:"gender"`
	CustomerType   string  `json:"customer_type" db:"customer_type"`
	LoyaltyPoints  int     `json:"loyalty_points" db:"loyalty_points"`
	Notes          string  `json:"notes" db:"notes"`
	IsActive       bool    `json:"is_active" db:"is_active"`
	
	// Relations
	Vehicles []CustomerVehicle `json:"vehicles,omitempty"`
}

// CustomerVehicle model
type CustomerVehicle struct {
	BaseModel
	CustomerID           uuid.UUID `json:"customer_id" db:"customer_id" validate:"required"`
	VehicleNumber        string  `json:"vehicle_number" db:"vehicle_number" validate:"required"`
	Brand                string  `json:"brand" db:"brand" validate:"required"`
	Model                string  `json:"model" db:"model" validate:"required"`
	Year                 int     `json:"year" db:"year" validate:"required"`
	Color                string  `json:"color" db:"color"`
	EngineNumber         string  `json:"engine_number" db:"engine_number"`
	ChassisNumber        string  `json:"chassis_number" db:"chassis_number"`
	FuelType             string  `json:"fuel_type" db:"fuel_type"`
	Transmission         string  `json:"transmission" db:"transmission"`
	Mileage              int64   `json:"mileage" db:"mileage"`
	LastServiceDate      *string `json:"last_service_date" db:"last_service_date"`
	NextServiceDate      *string `json:"next_service_date" db:"next_service_date"`
	InsuranceExpiry      *string `json:"insurance_expiry" db:"insurance_expiry"`
	RegistrationExpiry   *string `json:"registration_expiry" db:"registration_expiry"`
	Notes                string  `json:"notes" db:"notes"`
	IsActive             bool    `json:"is_active" db:"is_active"`
	
	// Relations
	Customer *Customer `json:"customer,omitempty"`
}

// Auth models
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// JWT Claims
type Claims struct {
	UserID   uuid.UUID   `json:"user_id"`
	Username string      `json:"username"`
	RoleID   uuid.UUID   `json:"role_id"`
	OutletID *uuid.UUID  `json:"outlet_id"`
	Permissions []string `json:"permissions"`
}
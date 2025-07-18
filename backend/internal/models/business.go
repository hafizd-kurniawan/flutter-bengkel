package models

import (
	"time"
)

// Service Category model
type ServiceCategory struct {
	BaseModel
	Name        string `json:"name" db:"name" validate:"required"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// Service model
type Service struct {
	BaseModel
	ServiceCode       string  `json:"service_code" db:"service_code" validate:"required"`
	Name              string  `json:"name" db:"name" validate:"required"`
	Description       string  `json:"description" db:"description"`
	CategoryID        int64   `json:"category_id" db:"category_id" validate:"required"`
	StandardPrice     float64 `json:"standard_price" db:"standard_price"`
	EstimatedDuration int     `json:"estimated_duration" db:"estimated_duration"` // in minutes
	IsActive          bool    `json:"is_active" db:"is_active"`
	
	// Relations
	Category *ServiceCategory `json:"category,omitempty"`
}

// UnitType model
type UnitType struct {
	BaseModel
	Name         string `json:"name" db:"name" validate:"required"`
	Abbreviation string `json:"abbreviation" db:"abbreviation" validate:"required"`
	Description  string `json:"description" db:"description"`
}

// Category model for products
type Category struct {
	BaseModel
	Name        string `json:"name" db:"name" validate:"required"`
	ParentID    *int64 `json:"parent_id" db:"parent_id"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	
	// Relations
	Parent   *Category  `json:"parent,omitempty"`
	Children []Category `json:"children,omitempty"`
}

// Supplier model
type Supplier struct {
	BaseModel
	SupplierCode   string `json:"supplier_code" db:"supplier_code" validate:"required"`
	Name           string `json:"name" db:"name" validate:"required"`
	Email          string `json:"email" db:"email"`
	Phone          string `json:"phone" db:"phone"`
	Address        string `json:"address" db:"address"`
	City           string `json:"city" db:"city"`
	Province       string `json:"province" db:"province"`
	PostalCode     string `json:"postal_code" db:"postal_code"`
	ContactPerson  string `json:"contact_person" db:"contact_person"`
	PaymentTerms   string `json:"payment_terms" db:"payment_terms"`
	IsActive       bool   `json:"is_active" db:"is_active"`
}

// Product model
type Product struct {
	BaseModel
	ProductCode      string  `json:"product_code" db:"product_code" validate:"required"`
	Name             string  `json:"name" db:"name" validate:"required"`
	Description      string  `json:"description" db:"description"`
	CategoryID       int64   `json:"category_id" db:"category_id" validate:"required"`
	UnitTypeID       int64   `json:"unit_type_id" db:"unit_type_id" validate:"required"`
	SupplierID       *int64  `json:"supplier_id" db:"supplier_id"`
	CostPrice        float64 `json:"cost_price" db:"cost_price"`
	SellingPrice     float64 `json:"selling_price" db:"selling_price"`
	StockQuantity    int     `json:"stock_quantity" db:"stock_quantity"`
	MinStockLevel    int     `json:"min_stock_level" db:"min_stock_level"`
	MaxStockLevel    int     `json:"max_stock_level" db:"max_stock_level"`
	HasSerialNumber  bool    `json:"has_serial_number" db:"has_serial_number"`
	IsService        bool    `json:"is_service" db:"is_service"`
	IsActive         bool    `json:"is_active" db:"is_active"`
	
	// Relations
	Category     *Category     `json:"category,omitempty"`
	UnitType     *UnitType     `json:"unit_type,omitempty"`
	Supplier     *Supplier     `json:"supplier,omitempty"`
	SerialNumbers []ProductSerialNumber `json:"serial_numbers,omitempty"`
}

// ProductSerialNumber model
type ProductSerialNumber struct {
	BaseModel
	ProductID    int64   `json:"product_id" db:"product_id" validate:"required"`
	SerialNumber string  `json:"serial_number" db:"serial_number" validate:"required"`
	Status       string  `json:"status" db:"status"` // available, sold, reserved, damaged
	PurchaseDate *string `json:"purchase_date" db:"purchase_date"`
	SaleDate     *string `json:"sale_date" db:"sale_date"`
	Notes        string  `json:"notes" db:"notes"`
	
	// Relations
	Product *Product `json:"product,omitempty"`
}

// ServiceJob model
type ServiceJob struct {
	BaseModel
	JobNumber            string    `json:"job_number" db:"job_number" validate:"required"`
	CustomerID           int64     `json:"customer_id" db:"customer_id" validate:"required"`
	VehicleID            int64     `json:"vehicle_id" db:"vehicle_id" validate:"required"`
	OutletID             int64     `json:"outlet_id" db:"outlet_id" validate:"required"`
	TechnicianID         *int64    `json:"technician_id" db:"technician_id"`
	QueueNumber          int       `json:"queue_number" db:"queue_number"`
	Priority             string    `json:"priority" db:"priority"` // low, normal, high, urgent
	Status               string    `json:"status" db:"status"`     // pending, in_progress, completed, cancelled, on_hold
	ProblemDescription   string    `json:"problem_description" db:"problem_description" validate:"required"`
	EstimatedCompletion  *time.Time `json:"estimated_completion" db:"estimated_completion"`
	ActualCompletion     *time.Time `json:"actual_completion" db:"actual_completion"`
	TotalAmount          float64   `json:"total_amount" db:"total_amount"`
	DiscountAmount       float64   `json:"discount_amount" db:"discount_amount"`
	TaxAmount            float64   `json:"tax_amount" db:"tax_amount"`
	FinalAmount          float64   `json:"final_amount" db:"final_amount"`
	WarrantyPeriodDays   int       `json:"warranty_period_days" db:"warranty_period_days"`
	Notes                string    `json:"notes" db:"notes"`
	
	// Relations
	Customer    *Customer        `json:"customer,omitempty"`
	Vehicle     *CustomerVehicle `json:"vehicle,omitempty"`
	Outlet      *Outlet          `json:"outlet,omitempty"`
	Technician  *User            `json:"technician,omitempty"`
	Details     []ServiceDetail  `json:"details,omitempty"`
	Histories   []ServiceJobHistory `json:"histories,omitempty"`
}

// ServiceDetail model
type ServiceDetail struct {
	BaseModel
	ServiceJobID int64   `json:"service_job_id" db:"service_job_id" validate:"required"`
	ProductID    *int64  `json:"product_id" db:"product_id"`
	ServiceID    *int64  `json:"service_id" db:"service_id"`
	Quantity     float64 `json:"quantity" db:"quantity" validate:"required"`
	UnitPrice    float64 `json:"unit_price" db:"unit_price" validate:"required"`
	TotalPrice   float64 `json:"total_price" db:"total_price" validate:"required"`
	Notes        string  `json:"notes" db:"notes"`
	
	// Relations
	ServiceJob *ServiceJob `json:"service_job,omitempty"`
	Product    *Product    `json:"product,omitempty"`
	Service    *Service    `json:"service,omitempty"`
}

// ServiceJobHistory model
type ServiceJobHistory struct {
	BaseModel
	ServiceJobID   int64  `json:"service_job_id" db:"service_job_id" validate:"required"`
	UserID         int64  `json:"user_id" db:"user_id" validate:"required"`
	PreviousStatus string `json:"previous_status" db:"previous_status"`
	NewStatus      string `json:"new_status" db:"new_status" validate:"required"`
	Notes          string `json:"notes" db:"notes"`
	
	// Relations
	ServiceJob *ServiceJob `json:"service_job,omitempty"`
	User       *User       `json:"user,omitempty"`
}

// Transaction model
type Transaction struct {
	BaseModel
	TransactionNumber string    `json:"transaction_number" db:"transaction_number" validate:"required"`
	TransactionType   string    `json:"transaction_type" db:"transaction_type" validate:"required"` // service, sparepart_sale, vehicle_purchase, vehicle_sale
	CustomerID        *int64    `json:"customer_id" db:"customer_id"`
	OutletID          int64     `json:"outlet_id" db:"outlet_id" validate:"required"`
	UserID            int64     `json:"user_id" db:"user_id" validate:"required"`
	ServiceJobID      *int64    `json:"service_job_id" db:"service_job_id"`
	SubtotalAmount    float64   `json:"subtotal_amount" db:"subtotal_amount"`
	DiscountAmount    float64   `json:"discount_amount" db:"discount_amount"`
	TaxAmount         float64   `json:"tax_amount" db:"tax_amount"`
	TotalAmount       float64   `json:"total_amount" db:"total_amount"`
	PaymentStatus     string    `json:"payment_status" db:"payment_status"` // pending, partial, paid, cancelled
	Notes             string    `json:"notes" db:"notes"`
	TransactionDate   time.Time `json:"transaction_date" db:"transaction_date"`
	
	// Relations
	Customer    *Customer           `json:"customer,omitempty"`
	Outlet      *Outlet             `json:"outlet,omitempty"`
	User        *User               `json:"user,omitempty"`
	ServiceJob  *ServiceJob         `json:"service_job,omitempty"`
	Details     []TransactionDetail `json:"details,omitempty"`
	Payments    []Payment           `json:"payments,omitempty"`
}

// TransactionDetail model
type TransactionDetail struct {
	BaseModel
	TransactionID int64   `json:"transaction_id" db:"transaction_id" validate:"required"`
	ProductID     *int64  `json:"product_id" db:"product_id"`
	ServiceID     *int64  `json:"service_id" db:"service_id"`
	Description   string  `json:"description" db:"description"`
	Quantity      float64 `json:"quantity" db:"quantity" validate:"required"`
	UnitPrice     float64 `json:"unit_price" db:"unit_price" validate:"required"`
	TotalPrice    float64 `json:"total_price" db:"total_price" validate:"required"`
	
	// Relations
	Transaction *Transaction `json:"transaction,omitempty"`
	Product     *Product     `json:"product,omitempty"`
	Service     *Service     `json:"service,omitempty"`
}

// PaymentMethod model
type PaymentMethod struct {
	BaseModel
	Name          string `json:"name" db:"name" validate:"required"`
	Type          string `json:"type" db:"type" validate:"required"` // cash, bank_transfer, credit_card, debit_card, e_wallet, check
	AccountNumber string `json:"account_number" db:"account_number"`
	BankName      string `json:"bank_name" db:"bank_name"`
	IsActive      bool   `json:"is_active" db:"is_active"`
}

// Payment model
type Payment struct {
	BaseModel
	PaymentNumber   string    `json:"payment_number" db:"payment_number" validate:"required"`
	TransactionID   int64     `json:"transaction_id" db:"transaction_id" validate:"required"`
	PaymentMethodID int64     `json:"payment_method_id" db:"payment_method_id" validate:"required"`
	Amount          float64   `json:"amount" db:"amount" validate:"required"`
	PaymentDate     time.Time `json:"payment_date" db:"payment_date"`
	ReferenceNumber string    `json:"reference_number" db:"reference_number"`
	Notes           string    `json:"notes" db:"notes"`
	
	// Relations
	Transaction   *Transaction   `json:"transaction,omitempty"`
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
}

// CreateServiceJobRequest
type CreateServiceJobRequest struct {
	CustomerID         int64  `json:"customer_id" validate:"required"`
	VehicleID          int64  `json:"vehicle_id" validate:"required"`
	Priority           string `json:"priority"`
	ProblemDescription string `json:"problem_description" validate:"required"`
	TechnicianID       *int64 `json:"technician_id"`
	WarrantyPeriodDays int    `json:"warranty_period_days"`
	Notes              string `json:"notes"`
}

// UpdateServiceJobRequest
type UpdateServiceJobRequest struct {
	TechnicianID        *int64     `json:"technician_id"`
	Priority            string     `json:"priority"`
	Status              string     `json:"status"`
	EstimatedCompletion *time.Time `json:"estimated_completion"`
	ActualCompletion    *time.Time `json:"actual_completion"`
	WarrantyPeriodDays  int        `json:"warranty_period_days"`
	Notes               string     `json:"notes"`
}

// CreateTransactionRequest
type CreateTransactionRequest struct {
	TransactionType string                      `json:"transaction_type" validate:"required"`
	CustomerID      *int64                      `json:"customer_id"`
	ServiceJobID    *int64                      `json:"service_job_id"`
	DiscountAmount  float64                     `json:"discount_amount"`
	TaxAmount       float64                     `json:"tax_amount"`
	Notes           string                      `json:"notes"`
	Details         []CreateTransactionDetailRequest `json:"details" validate:"required,dive"`
}

// CreateTransactionDetailRequest
type CreateTransactionDetailRequest struct {
	ProductID   *int64  `json:"product_id"`
	ServiceID   *int64  `json:"service_id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity" validate:"required"`
	UnitPrice   float64 `json:"unit_price" validate:"required"`
}

// CreatePaymentRequest
type CreatePaymentRequest struct {
	TransactionID   int64   `json:"transaction_id" validate:"required"`
	PaymentMethodID int64   `json:"payment_method_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
}
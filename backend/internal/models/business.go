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

// Multi-Payment Support Models (NEW)

// TransactionPayment model for multi-payment support
type TransactionPayment struct {
	ID              int64   `json:"id" db:"id"`
	PaymentID       int     `json:"payment_id" db:"payment_id"`
	TransactionID   int64   `json:"transaction_id" db:"transaction_id" validate:"required"`
	PaymentMethodID int64   `json:"payment_method_id" db:"payment_method_id" validate:"required"`
	Amount          float64 `json:"amount" db:"amount" validate:"required"`
	PaymentOrder    int     `json:"payment_order" db:"payment_order" validate:"required"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	
	// Relations
	Transaction   *Transaction   `json:"transaction,omitempty"`
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
}

// CreateMultiPaymentRequest for POS multi-payment transactions
type CreateMultiPaymentRequest struct {
	TransactionID int64                          `json:"transaction_id" validate:"required"`
	Payments      []CreateTransactionPaymentRequest `json:"payments" validate:"required,dive"`
}

// CreateTransactionPaymentRequest for individual payments in multi-payment
type CreateTransactionPaymentRequest struct {
	PaymentMethodID int64   `json:"payment_method_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	ReferenceNumber string  `json:"reference_number"`
}

// Enhanced Vehicle Trading Models

// VehicleSales model (from migration)
type VehicleSales struct {
	BaseModel
	SaleID            int     `json:"sale_id" db:"sale_id"`
	VehiclePurchaseID int64   `json:"vehicle_purchase_id" db:"vehicle_purchase_id" validate:"required"`
	CustomerID        int64   `json:"customer_id" db:"customer_id" validate:"required"`
	SalesUserID       int64   `json:"sales_user_id" db:"sales_user_id" validate:"required"`
	SaleDate          string  `json:"sale_date" db:"sale_date" validate:"required"`
	SellingPrice      float64 `json:"selling_price" db:"selling_price" validate:"required"`
	ProfitAmount      float64 `json:"profit_amount" db:"profit_amount"`
	CommissionAmount  float64 `json:"commission_amount" db:"commission_amount"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy         *int64  `json:"created_by" db:"created_by"`
	
	// Relations
	VehiclePurchase *VehiclePurchase `json:"vehicle_purchase,omitempty"`
	Customer        *Customer        `json:"customer,omitempty"`
	SalesUser       *User            `json:"sales_user,omitempty"`
}

// Enhanced VehiclePurchase model with new fields
type VehiclePurchase struct {
	BaseModel
	PurchaseNumber      string     `json:"purchase_number" db:"purchase_number" validate:"required"`
	VehicleNumber       string     `json:"vehicle_number" db:"vehicle_number" validate:"required"`
	Brand               string     `json:"brand" db:"brand" validate:"required"`
	Model               string     `json:"model" db:"model" validate:"required"`
	Year                int        `json:"year" db:"year" validate:"required"`
	Color               string     `json:"color" db:"color"`
	Mileage             *int64     `json:"mileage" db:"mileage"`
	ConditionRating     string     `json:"condition_rating" db:"condition_rating"`
	PurchasePrice       float64    `json:"purchase_price" db:"purchase_price" validate:"required"`
	EstimatedSellingPrice *float64 `json:"estimated_selling_price" db:"estimated_selling_price"`
	ActualSellingPrice  *float64   `json:"actual_selling_price" db:"actual_selling_price"`
	SellerName          string     `json:"seller_name" db:"seller_name"`
	SellerPhone         string     `json:"seller_phone" db:"seller_phone"`
	SellerAddress       string     `json:"seller_address" db:"seller_address"`
	Status              string     `json:"status" db:"status"`
	PurchaseDate        string     `json:"purchase_date" db:"purchase_date" validate:"required"`
	SaleDate            *string    `json:"sale_date" db:"sale_date"`
	OutletID            int64      `json:"outlet_id" db:"outlet_id" validate:"required"`
	UserID              int64      `json:"user_id" db:"user_id" validate:"required"`
	Notes               string     `json:"notes" db:"notes"`
	
	// New fields for enhanced vehicle trading
	ServiceRequired        bool       `json:"service_required" db:"service_required"`
	ServiceCompletionDate  *time.Time `json:"service_completion_date" db:"service_completion_date"`
	SellingPrice           *float64   `json:"selling_price" db:"selling_price"`
	SaleStatus             string     `json:"sale_status" db:"sale_status"`
	
	// Relations
	Outlet *Outlet `json:"outlet,omitempty"`
	User   *User   `json:"user,omitempty"`
}

// Commission Tracking Model
type CommissionTracking struct {
	BaseModel
	UserID           int64      `json:"user_id" db:"user_id" validate:"required"`
	TransactionID    *int64     `json:"transaction_id" db:"transaction_id"`
	VehicleSaleID    *int64     `json:"vehicle_sale_id" db:"vehicle_sale_id"`
	CommissionType   string     `json:"commission_type" db:"commission_type" validate:"required"`
	BaseAmount       float64    `json:"base_amount" db:"base_amount" validate:"required"`
	CommissionRate   float64    `json:"commission_rate" db:"commission_rate" validate:"required"`
	CommissionAmount float64    `json:"commission_amount" db:"commission_amount" validate:"required"`
	Status           string     `json:"status" db:"status"`
	ApprovedBy       *int64     `json:"approved_by" db:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at" db:"approved_at"`
	
	// Relations
	User          *User         `json:"user,omitempty"`
	Transaction   *Transaction  `json:"transaction,omitempty"`
	VehicleSale   *VehicleSales `json:"vehicle_sale,omitempty"`
	ApproverUser  *User         `json:"approver_user,omitempty"`
}

// Enhanced AccountsReceivable model with kasir approval
type AccountsReceivable struct {
	BaseModel
	ARNumber         string     `json:"ar_number" db:"ar_number" validate:"required"`
	CustomerID       int64      `json:"customer_id" db:"customer_id" validate:"required"`
	TransactionID    *int64     `json:"transaction_id" db:"transaction_id"`
	Amount           float64    `json:"amount" db:"amount" validate:"required"`
	PaidAmount       float64    `json:"paid_amount" db:"paid_amount"`
	RemainingAmount  float64    `json:"remaining_amount" db:"remaining_amount"`
	DueDate          string     `json:"due_date" db:"due_date" validate:"required"`
	Status           string     `json:"status" db:"status"`
	Notes            string     `json:"notes" db:"notes"`
	ApprovedBy       *int64     `json:"approved_by" db:"approved_by"`
	ApprovalDate     *time.Time `json:"approval_date" db:"approval_date"`
	
	// Relations
	Customer    *Customer    `json:"customer,omitempty"`
	Transaction *Transaction `json:"transaction,omitempty"`
	ApproverUser *User       `json:"approver_user,omitempty"`
}

// Enhanced ReceivablePayment model with kasir approval
type ReceivablePayment struct {
	BaseModel
	PaymentNumber        string     `json:"payment_number" db:"payment_number" validate:"required"`
	AccountsReceivableID int64      `json:"accounts_receivable_id" db:"accounts_receivable_id" validate:"required"`
	PaymentMethodID      int64      `json:"payment_method_id" db:"payment_method_id" validate:"required"`
	Amount               float64    `json:"amount" db:"amount" validate:"required"`
	PaymentDate          time.Time  `json:"payment_date" db:"payment_date"`
	ReferenceNumber      string     `json:"reference_number" db:"reference_number"`
	Notes                string     `json:"notes" db:"notes"`
	ApprovedBy           *int64     `json:"approved_by" db:"approved_by"`
	ApprovalDate         *time.Time `json:"approval_date" db:"approval_date"`
	
	// Relations
	AccountsReceivable *AccountsReceivable `json:"accounts_receivable,omitempty"`
	PaymentMethod      *PaymentMethod      `json:"payment_method,omitempty"`
	ApproverUser       *User               `json:"approver_user,omitempty"`
}

// Request models for new features

// CreateVehiclePurchaseRequest
type CreateVehiclePurchaseRequest struct {
	VehicleNumber         string  `json:"vehicle_number" validate:"required"`
	Brand                 string  `json:"brand" validate:"required"`
	Model                 string  `json:"model" validate:"required"`
	Year                  int     `json:"year" validate:"required"`
	Color                 string  `json:"color"`
	Mileage               *int64  `json:"mileage"`
	ConditionRating       string  `json:"condition_rating"`
	PurchasePrice         float64 `json:"purchase_price" validate:"required"`
	EstimatedSellingPrice *float64 `json:"estimated_selling_price"`
	SellerName            string  `json:"seller_name"`
	SellerPhone           string  `json:"seller_phone"`
	SellerAddress         string  `json:"seller_address"`
	Notes                 string  `json:"notes"`
	ServiceRequired       bool    `json:"service_required"`
}

// CreateVehicleSaleRequest
type CreateVehicleSaleRequest struct {
	VehiclePurchaseID int64   `json:"vehicle_purchase_id" validate:"required"`
	CustomerID        int64   `json:"customer_id" validate:"required"`
	SellingPrice      float64 `json:"selling_price" validate:"required"`
	CommissionAmount  float64 `json:"commission_amount"`
	Notes             string  `json:"notes"`
}

// UpdateVehiclePriceRequest (for sales team)
type UpdateVehiclePriceRequest struct {
	SellingPrice *float64 `json:"selling_price" validate:"required"`
	Notes        string   `json:"notes"`
}

// ApproveReceivablePaymentRequest (for kasir)
type ApproveReceivablePaymentRequest struct {
	PaymentMethodID int64   `json:"payment_method_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
}

// POS-specific request models

// CreatePOSTransactionRequest (enhanced for kasir)
type CreatePOSTransactionRequest struct {
	CustomerID      *int64                          `json:"customer_id"`
	DiscountAmount  float64                         `json:"discount_amount"`
	TaxAmount       float64                         `json:"tax_amount"`
	Notes           string                          `json:"notes"`
	Details         []CreateTransactionDetailRequest `json:"details" validate:"required,dive"`
	Payments        []CreateTransactionPaymentRequest `json:"payments" validate:"required,dive"`
}

// ProductSearchRequest for POS barcode/search
type ProductSearchRequest struct {
	Query    string `json:"query"`
	Barcode  string `json:"barcode"`
	Category *int64 `json:"category"`
	Limit    int    `json:"limit"`
}

// QueueManagementRequest for service jobs
type QueueManagementRequest struct {
	OutletID int64 `json:"outlet_id" validate:"required"`
	Date     string `json:"date"` // YYYY-MM-DD format
}

// Auto-assign mechanic request
type AutoAssignMechanicRequest struct {
	ServiceJobID int64  `json:"service_job_id" validate:"required"`
	ServiceType  string `json:"service_type"`
}
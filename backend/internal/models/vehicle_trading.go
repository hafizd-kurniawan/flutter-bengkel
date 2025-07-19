package models

import (
	"time"
	"github.com/shopspring/decimal"
)

// Enhanced BaseModel with soft delete support
type BaseModelWithSoftDelete struct {
	ID        int64      `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedBy *int64     `json:"created_by,omitempty" db:"created_by"`
}

// Vehicle Trading Models

// VehiclePurchase - Record when buying vehicles from customers
type VehiclePurchase struct {
	PurchaseID    int64           `json:"purchase_id" db:"purchase_id"`
	CustomerID    int64           `json:"customer_id" db:"customer_id" validate:"required"`
	OutletID      int64           `json:"outlet_id" db:"outlet_id" validate:"required"`
	PurchaseDate  time.Time       `json:"purchase_date" db:"purchase_date" validate:"required"`
	PurchasePrice decimal.Decimal `json:"purchase_price" db:"purchase_price" validate:"required"`
	PaymentMethod string          `json:"payment_method" db:"payment_method"`
	Notes         string          `json:"notes" db:"notes"`
	Status        string          `json:"status" db:"status"` // pending, completed, cancelled
	BaseModelWithSoftDelete
	
	// Relations
	Customer  *Customer     `json:"customer,omitempty"`
	Outlet    *Outlet       `json:"outlet,omitempty"`
	Inventory []VehicleInventory `json:"inventory,omitempty"`
}

// VehicleInventory - Vehicles available for sale
type VehicleInventory struct {
	InventoryID           int64             `json:"inventory_id" db:"inventory_id"`
	VehiclePurchaseID     int64             `json:"vehicle_purchase_id" db:"vehicle_purchase_id" validate:"required"`
	PlateNumber           string            `json:"plate_number" db:"plate_number" validate:"required"`
	Brand                 string            `json:"brand" db:"brand" validate:"required"`
	Model                 string            `json:"model" db:"model" validate:"required"`
	Type                  string            `json:"type" db:"type" validate:"required"`
	ProductionYear        int               `json:"production_year" db:"production_year" validate:"required"`
	ChassisNumber         string            `json:"chassis_number" db:"chassis_number" validate:"required"`
	EngineNumber          string            `json:"engine_number" db:"engine_number" validate:"required"`
	Color                 string            `json:"color" db:"color" validate:"required"`
	Mileage               int               `json:"mileage" db:"mileage"`
	ConditionRating       int               `json:"condition_rating" db:"condition_rating"` // 1-5 scale
	PurchasePrice         decimal.Decimal   `json:"purchase_price" db:"purchase_price" validate:"required"`
	EstimatedSellingPrice decimal.Decimal   `json:"estimated_selling_price" db:"estimated_selling_price" validate:"required"`
	ActualSellingPrice    *decimal.Decimal  `json:"actual_selling_price,omitempty" db:"actual_selling_price"`
	Status                string            `json:"status" db:"status"` // Available, Reserved, Sold, Under_Maintenance
	VehiclePhotos         []string          `json:"vehicle_photos,omitempty" db:"vehicle_photos"`
	ConditionNotes        string            `json:"condition_notes" db:"condition_notes"`
	SellingDate           *time.Time        `json:"selling_date,omitempty" db:"selling_date"`
	ProfitMargin          *decimal.Decimal  `json:"profit_margin,omitempty" db:"profit_margin"`
	BaseModelWithSoftDelete
	
	// Relations
	VehiclePurchase *VehiclePurchase           `json:"vehicle_purchase,omitempty"`
	Photos          []VehiclePhoto             `json:"photos,omitempty"`
	Assessments     []VehicleConditionAssessment `json:"assessments,omitempty"`
	Sales           []VehicleSale              `json:"sales,omitempty"`
}

// VehicleSale - Record when selling vehicles to customers
type VehicleSale struct {
	SaleID            int64           `json:"sale_id" db:"sale_id"`
	InventoryID       int64           `json:"inventory_id" db:"inventory_id" validate:"required"`
	CustomerID        int64           `json:"customer_id" db:"customer_id" validate:"required"`
	SalesPersonID     int64           `json:"sales_person_id" db:"sales_person_id" validate:"required"`
	OutletID          int64           `json:"outlet_id" db:"outlet_id" validate:"required"`
	SaleDate          time.Time       `json:"sale_date" db:"sale_date" validate:"required"`
	SellingPrice      decimal.Decimal `json:"selling_price" db:"selling_price" validate:"required"`
	CommissionRate    decimal.Decimal `json:"commission_rate" db:"commission_rate"`
	CommissionAmount  decimal.Decimal `json:"commission_amount" db:"commission_amount"`
	PaymentType       string          `json:"payment_type" db:"payment_type" validate:"required"` // cash, credit, trade_in, financing
	DownPayment       decimal.Decimal `json:"down_payment" db:"down_payment"`
	FinancingAmount   decimal.Decimal `json:"financing_amount" db:"financing_amount"`
	FinancingBank     string          `json:"financing_bank" db:"financing_bank"`
	FinancingTermMonths int           `json:"financing_term_months" db:"financing_term_months"`
	Status            string          `json:"status" db:"status"` // Pending, Completed, Cancelled
	Notes             string          `json:"notes" db:"notes"`
	BaseModelWithSoftDelete
	
	// Relations
	Inventory    *VehicleInventory `json:"inventory,omitempty"`
	Customer     *Customer         `json:"customer,omitempty"`
	SalesPerson  *User             `json:"sales_person,omitempty"`
	Outlet       *Outlet           `json:"outlet,omitempty"`
	Commissions  []SalesCommission `json:"commissions,omitempty"`
}

// VehicleConditionAssessment - Detailed vehicle condition assessment
type VehicleConditionAssessment struct {
	AssessmentID         int64           `json:"assessment_id" db:"assessment_id"`
	InventoryID          int64           `json:"inventory_id" db:"inventory_id" validate:"required"`
	AssessorID           int64           `json:"assessor_id" db:"assessor_id" validate:"required"`
	AssessmentDate       time.Time       `json:"assessment_date" db:"assessment_date" validate:"required"`
	ExteriorCondition    int             `json:"exterior_condition" db:"exterior_condition"` // 1-5
	InteriorCondition    int             `json:"interior_condition" db:"interior_condition"` // 1-5
	EngineCondition      int             `json:"engine_condition" db:"engine_condition"` // 1-5
	TransmissionCondition int            `json:"transmission_condition" db:"transmission_condition"` // 1-5
	TireCondition        int             `json:"tire_condition" db:"tire_condition"` // 1-5
	ElectricalCondition  int             `json:"electrical_condition" db:"electrical_condition"` // 1-5
	OverallRating        int             `json:"overall_rating" db:"overall_rating"` // 1-5
	AssessmentNotes      string          `json:"assessment_notes" db:"assessment_notes"`
	RecommendedRepairs   string          `json:"recommended_repairs" db:"recommended_repairs"`
	EstimatedRepairCost  decimal.Decimal `json:"estimated_repair_cost" db:"estimated_repair_cost"`
	BaseModelWithSoftDelete
	
	// Relations
	Inventory *VehicleInventory `json:"inventory,omitempty"`
	Assessor  *User             `json:"assessor,omitempty"`
}

// VehiclePhoto - Vehicle photo management
type VehiclePhoto struct {
	PhotoID     int64  `json:"photo_id" db:"photo_id"`
	InventoryID int64  `json:"inventory_id" db:"inventory_id" validate:"required"`
	PhotoURL    string `json:"photo_url" db:"photo_url" validate:"required"`
	PhotoType   string `json:"photo_type" db:"photo_type"` // exterior_front, exterior_back, etc.
	Description string `json:"description" db:"description"`
	IsPrimary   bool   `json:"is_primary" db:"is_primary"`
	SortOrder   int    `json:"sort_order" db:"sort_order"`
	BaseModelWithSoftDelete
	
	// Relations
	Inventory *VehicleInventory `json:"inventory,omitempty"`
}

// SalesCommission - Commission tracking for sales team
type SalesCommission struct {
	CommissionID    int64           `json:"commission_id" db:"commission_id"`
	SaleID          int64           `json:"sale_id" db:"sale_id" validate:"required"`
	SalesPersonID   int64           `json:"sales_person_id" db:"sales_person_id" validate:"required"`
	CommissionRate  decimal.Decimal `json:"commission_rate" db:"commission_rate" validate:"required"`
	CommissionAmount decimal.Decimal `json:"commission_amount" db:"commission_amount" validate:"required"`
	PaymentStatus   string          `json:"payment_status" db:"payment_status"` // pending, paid, cancelled
	PaymentDate     *time.Time      `json:"payment_date,omitempty" db:"payment_date"`
	Notes           string          `json:"notes" db:"notes"`
	BaseModelWithSoftDelete
	
	// Relations
	Sale        *VehicleSale `json:"sale,omitempty"`
	SalesPerson *User        `json:"sales_person,omitempty"`
}

// Request/Response Models for Vehicle Trading

// CreateVehiclePurchaseRequest
type CreateVehiclePurchaseRequest struct {
	CustomerID    int64           `json:"customer_id" validate:"required"`
	PurchasePrice decimal.Decimal `json:"purchase_price" validate:"required"`
	PaymentMethod string          `json:"payment_method"`
	Notes         string          `json:"notes"`
	VehicleDetails CreateVehicleInventoryRequest `json:"vehicle_details" validate:"required"`
}

// CreateVehicleInventoryRequest
type CreateVehicleInventoryRequest struct {
	PlateNumber           string          `json:"plate_number" validate:"required"`
	Brand                 string          `json:"brand" validate:"required"`
	Model                 string          `json:"model" validate:"required"`
	Type                  string          `json:"type" validate:"required"`
	ProductionYear        int             `json:"production_year" validate:"required"`
	ChassisNumber         string          `json:"chassis_number" validate:"required"`
	EngineNumber          string          `json:"engine_number" validate:"required"`
	Color                 string          `json:"color" validate:"required"`
	Mileage               int             `json:"mileage"`
	ConditionRating       int             `json:"condition_rating" validate:"min=1,max=5"`
	EstimatedSellingPrice decimal.Decimal `json:"estimated_selling_price" validate:"required"`
	ConditionNotes        string          `json:"condition_notes"`
}

// CreateVehicleSaleRequest
type CreateVehicleSaleRequest struct {
	InventoryID         int64           `json:"inventory_id" validate:"required"`
	CustomerID          int64           `json:"customer_id" validate:"required"`
	SellingPrice        decimal.Decimal `json:"selling_price" validate:"required"`
	PaymentType         string          `json:"payment_type" validate:"required"`
	DownPayment         decimal.Decimal `json:"down_payment"`
	FinancingAmount     decimal.Decimal `json:"financing_amount"`
	FinancingBank       string          `json:"financing_bank"`
	FinancingTermMonths int             `json:"financing_term_months"`
	Notes               string          `json:"notes"`
}

// UpdateVehicleInventoryRequest
type UpdateVehicleInventoryRequest struct {
	EstimatedSellingPrice *decimal.Decimal `json:"estimated_selling_price"`
	Status                *string          `json:"status"`
	ConditionRating       *int             `json:"condition_rating"`
	ConditionNotes        *string          `json:"condition_notes"`
	Mileage               *int             `json:"mileage"`
}

// VehicleSearchRequest
type VehicleSearchRequest struct {
	Brand        string           `json:"brand,omitempty"`
	Model        string           `json:"model,omitempty"`
	Type         string           `json:"type,omitempty"`
	YearFrom     int              `json:"year_from,omitempty"`
	YearTo       int              `json:"year_to,omitempty"`
	PriceFrom    *decimal.Decimal `json:"price_from,omitempty"`
	PriceTo      *decimal.Decimal `json:"price_to,omitempty"`
	Status       string           `json:"status,omitempty"`
	ConditionMin int              `json:"condition_min,omitempty"`
	ConditionMax int              `json:"condition_max,omitempty"`
	Page         int              `json:"page,omitempty"`
	PerPage      int              `json:"per_page,omitempty"`
}

// VehicleTradingAnalytics - Analytics and reporting models
type ProfitAnalysis struct {
	TotalPurchases        int             `json:"total_purchases"`
	TotalSales            int             `json:"total_sales"`
	TotalPurchaseAmount   decimal.Decimal `json:"total_purchase_amount"`
	TotalSalesAmount      decimal.Decimal `json:"total_sales_amount"`
	TotalProfitAmount     decimal.Decimal `json:"total_profit_amount"`
	AverageProfitMargin   decimal.Decimal `json:"average_profit_margin"`
	VehiclesInStock       int             `json:"vehicles_in_stock"`
	VehiclesStockValue    decimal.Decimal `json:"vehicles_stock_value"`
	TopSellingBrands      []BrandSalesData `json:"top_selling_brands"`
	MonthlyTrends         []MonthlySalesData `json:"monthly_trends"`
}

type BrandSalesData struct {
	Brand       string          `json:"brand"`
	TotalSales  int             `json:"total_sales"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	AvgPrice    decimal.Decimal `json:"avg_price"`
}

type MonthlySalesData struct {
	Month       string          `json:"month"`
	TotalSales  int             `json:"total_sales"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Profit      decimal.Decimal `json:"profit"`
}

// InventoryAgingReport
type InventoryAgingReport struct {
	InventoryID   int64           `json:"inventory_id"`
	PlateNumber   string          `json:"plate_number"`
	Brand         string          `json:"brand"`
	Model         string          `json:"model"`
	PurchaseDate  time.Time       `json:"purchase_date"`
	DaysInStock   int             `json:"days_in_stock"`
	PurchasePrice decimal.Decimal `json:"purchase_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	EstimatedLoss decimal.Decimal `json:"estimated_loss"`
}

// SalesPerformanceReport
type SalesPerformanceReport struct {
	SalesPersonID    int64           `json:"sales_person_id"`
	SalesPersonName  string          `json:"sales_person_name"`
	TotalSales       int             `json:"total_sales"`
	TotalSalesAmount decimal.Decimal `json:"total_sales_amount"`
	TotalCommission  decimal.Decimal `json:"total_commission"`
	AvgSaleValue     decimal.Decimal `json:"avg_sale_value"`
	ConversionRate   decimal.Decimal `json:"conversion_rate"`
}
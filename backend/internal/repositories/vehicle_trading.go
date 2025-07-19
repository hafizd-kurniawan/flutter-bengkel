package repositories

import (
	"fmt"
	"strings"
	"time"

	"flutter-bengkel/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// VehicleTradingRepository interface defines vehicle trading operations
type VehicleTradingRepository interface {
	// Vehicle Purchase
	CreateVehiclePurchase(purchase *models.VehiclePurchase) error
	GetVehiclePurchaseByID(id int64) (*models.VehiclePurchase, error)
	GetVehiclePurchases(offset, limit int) ([]*models.VehiclePurchase, int64, error)
	UpdateVehiclePurchase(id int64, purchase *models.VehiclePurchase) error
	SoftDeleteVehiclePurchase(id int64, deletedBy int64) error

	// Vehicle Inventory
	CreateVehicleInventory(inventory *models.VehicleInventory) error
	GetVehicleInventoryByID(id int64) (*models.VehicleInventory, error)
	GetAvailableVehicles() ([]*models.VehicleInventory, error)
	SearchVehicles(req *models.VehicleSearchRequest) ([]*models.VehicleInventory, int64, error)
	UpdateVehicleInventory(id int64, inventory *models.UpdateVehicleInventoryRequest) error
	UpdateSellingPrice(id int64, price decimal.Decimal) error
	MarkAsSold(id int64, saleData *models.VehicleSale) error
	SoftDeleteVehicleInventory(id int64, deletedBy int64) error

	// Vehicle Sales
	CreateVehicleSale(sale *models.VehicleSale) error
	GetVehicleSaleByID(id int64) (*models.VehicleSale, error)
	GetVehicleSales(offset, limit int) ([]*models.VehicleSale, int64, error)
	UpdateVehicleSale(id int64, sale *models.VehicleSale) error
	SoftDeleteVehicleSale(id int64, deletedBy int64) error

	// Vehicle Photos
	CreateVehiclePhoto(photo *models.VehiclePhoto) error
	GetVehiclePhotos(inventoryID int64) ([]*models.VehiclePhoto, error)
	UpdateVehiclePhoto(id int64, photo *models.VehiclePhoto) error
	SoftDeleteVehiclePhoto(id int64, deletedBy int64) error

	// Vehicle Assessments
	CreateVehicleAssessment(assessment *models.VehicleConditionAssessment) error
	GetVehicleAssessments(inventoryID int64) ([]*models.VehicleConditionAssessment, error)
	UpdateVehicleAssessment(id int64, assessment *models.VehicleConditionAssessment) error

	// Commission Management
	CreateSalesCommission(commission *models.SalesCommission) error
	GetSalesCommissions(salesPersonID int64, startDate, endDate time.Time) ([]*models.SalesCommission, error)
	UpdateCommissionPaymentStatus(id int64, status string, paymentDate *time.Time) error

	// Analytics & Reports
	GetProfitAnalysis(startDate, endDate time.Time, outletID *int64) (*models.ProfitAnalysis, error)
	GetInventoryAgingReport(outletID *int64) ([]*models.InventoryAgingReport, error)
	GetSalesPerformanceReport(startDate, endDate time.Time, outletID *int64) ([]*models.SalesPerformanceReport, error)
}

type vehicleTradingRepository struct {
	db *sqlx.DB
}

// NewVehicleTradingRepository creates a new vehicle trading repository
func NewVehicleTradingRepository(db *sqlx.DB) VehicleTradingRepository {
	return &vehicleTradingRepository{db: db}
}

// Vehicle Purchase Operations
func (r *vehicleTradingRepository) CreateVehiclePurchase(purchase *models.VehiclePurchase) error {
	query := `
		INSERT INTO vehicle_purchases (customer_id, outlet_id, purchase_date, purchase_price, 
			payment_method, notes, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING purchase_id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query, purchase.CustomerID, purchase.OutletID, 
		purchase.PurchaseDate, purchase.PurchasePrice, purchase.PaymentMethod,
		purchase.Notes, purchase.Status, purchase.CreatedBy).
		Scan(&purchase.PurchaseID, &purchase.CreatedAt, &purchase.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create vehicle purchase: %w", err)
	}
	
	return nil
}

func (r *vehicleTradingRepository) GetVehiclePurchaseByID(id int64) (*models.VehiclePurchase, error) {
	query := `
		SELECT vp.purchase_id, vp.customer_id, vp.outlet_id, vp.purchase_date, 
			   vp.purchase_price, vp.payment_method, vp.notes, vp.status,
			   vp.created_at, vp.updated_at, vp.deleted_at, vp.created_by,
			   c.name as "customer.name", c.phone as "customer.phone",
			   o.name as "outlet.name"
		FROM vehicle_purchases vp
		LEFT JOIN customers c ON vp.customer_id = c.customer_id
		LEFT JOIN outlets o ON vp.outlet_id = o.outlet_id
		WHERE vp.purchase_id = $1 AND vp.deleted_at IS NULL
	`
	
	var purchase models.VehiclePurchase
	err := r.db.Get(&purchase, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle purchase: %w", err)
	}
	
	return &purchase, nil
}

func (r *vehicleTradingRepository) GetVehiclePurchases(offset, limit int) ([]*models.VehiclePurchase, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM vehicle_purchases WHERE deleted_at IS NULL`
	var total int64
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get purchase count: %w", err)
	}
	
	// Get purchases
	query := `
		SELECT vp.purchase_id, vp.customer_id, vp.outlet_id, vp.purchase_date, 
			   vp.purchase_price, vp.payment_method, vp.notes, vp.status,
			   vp.created_at, vp.updated_at, vp.deleted_at, vp.created_by,
			   c.name as "customer.name", c.phone as "customer.phone",
			   o.name as "outlet.name"
		FROM vehicle_purchases vp
		LEFT JOIN customers c ON vp.customer_id = c.customer_id
		LEFT JOIN outlets o ON vp.outlet_id = o.outlet_id
		WHERE vp.deleted_at IS NULL
		ORDER BY vp.created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	var purchases []*models.VehiclePurchase
	err = r.db.Select(&purchases, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get vehicle purchases: %w", err)
	}
	
	return purchases, total, nil
}

func (r *vehicleTradingRepository) UpdateVehiclePurchase(id int64, purchase *models.VehiclePurchase) error {
	query := `
		UPDATE vehicle_purchases 
		SET purchase_date = $2, purchase_price = $3, payment_method = $4, 
			notes = $5, status = $6, updated_at = CURRENT_TIMESTAMP
		WHERE purchase_id = $1 AND deleted_at IS NULL
	`
	
	_, err := r.db.Exec(query, id, purchase.PurchaseDate, purchase.PurchasePrice,
		purchase.PaymentMethod, purchase.Notes, purchase.Status)
	if err != nil {
		return fmt.Errorf("failed to update vehicle purchase: %w", err)
	}
	
	return nil
}

func (r *vehicleTradingRepository) SoftDeleteVehiclePurchase(id int64, deletedBy int64) error {
	query := `
		UPDATE vehicle_purchases 
		SET deleted_at = CURRENT_TIMESTAMP, created_by = $2
		WHERE purchase_id = $1 AND deleted_at IS NULL
	`
	
	_, err := r.db.Exec(query, id, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle purchase: %w", err)
	}
	
	return nil
}

// Vehicle Inventory Operations
func (r *vehicleTradingRepository) CreateVehicleInventory(inventory *models.VehicleInventory) error {
	query := `
		INSERT INTO vehicle_inventory (vehicle_purchase_id, plate_number, brand, model, type,
			production_year, chassis_number, engine_number, color, mileage, condition_rating,
			purchase_price, estimated_selling_price, status, condition_notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING inventory_id, created_at, updated_at
	`
	
	err := r.db.QueryRow(query, inventory.VehiclePurchaseID, inventory.PlateNumber,
		inventory.Brand, inventory.Model, inventory.Type, inventory.ProductionYear,
		inventory.ChassisNumber, inventory.EngineNumber, inventory.Color,
		inventory.Mileage, inventory.ConditionRating, inventory.PurchasePrice,
		inventory.EstimatedSellingPrice, inventory.Status, inventory.ConditionNotes,
		inventory.CreatedBy).
		Scan(&inventory.InventoryID, &inventory.CreatedAt, &inventory.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create vehicle inventory: %w", err)
	}
	
	return nil
}

func (r *vehicleTradingRepository) GetVehicleInventoryByID(id int64) (*models.VehicleInventory, error) {
	query := `
		SELECT vi.inventory_id, vi.vehicle_purchase_id, vi.plate_number, vi.brand, 
			   vi.model, vi.type, vi.production_year, vi.chassis_number, vi.engine_number,
			   vi.color, vi.mileage, vi.condition_rating, vi.purchase_price, 
			   vi.estimated_selling_price, vi.actual_selling_price, vi.status, 
			   vi.condition_notes, vi.selling_date, vi.profit_margin,
			   vi.created_at, vi.updated_at, vi.deleted_at, vi.created_by,
			   vp.purchase_date, vp.purchase_price as original_purchase_price,
			   c.name as "customer.name", c.phone as "customer.phone"
		FROM vehicle_inventory vi
		LEFT JOIN vehicle_purchases vp ON vi.vehicle_purchase_id = vp.purchase_id
		LEFT JOIN customers c ON vp.customer_id = c.customer_id
		WHERE vi.inventory_id = $1 AND vi.deleted_at IS NULL
	`
	
	var inventory models.VehicleInventory
	err := r.db.Get(&inventory, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle inventory: %w", err)
	}
	
	return &inventory, nil
}

func (r *vehicleTradingRepository) GetAvailableVehicles() ([]*models.VehicleInventory, error) {
	query := `
		SELECT vi.inventory_id, vi.vehicle_purchase_id, vi.plate_number, vi.brand, 
			   vi.model, vi.type, vi.production_year, vi.chassis_number, vi.engine_number,
			   vi.color, vi.mileage, vi.condition_rating, vi.purchase_price, 
			   vi.estimated_selling_price, vi.actual_selling_price, vi.status, 
			   vi.condition_notes, vi.selling_date, vi.profit_margin,
			   vi.created_at, vi.updated_at, vi.deleted_at, vi.created_by,
			   vp.purchase_date, vp.purchase_price as original_purchase_price,
			   c.name as "seller.name", c.phone as "seller.phone"
		FROM vehicle_inventory vi
		LEFT JOIN vehicle_purchases vp ON vi.vehicle_purchase_id = vp.purchase_id
		LEFT JOIN customers c ON vp.customer_id = c.customer_id
		WHERE vi.status = 'Available' AND vi.deleted_at IS NULL
		ORDER BY vi.created_at DESC
	`
	
	var vehicles []*models.VehicleInventory
	err := r.db.Select(&vehicles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get available vehicles: %w", err)
	}
	
	return vehicles, nil
}

func (r *vehicleTradingRepository) SearchVehicles(req *models.VehicleSearchRequest) ([]*models.VehicleInventory, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	// Base condition
	conditions = append(conditions, "vi.deleted_at IS NULL")
	
	// Build dynamic WHERE clause
	if req.Brand != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(vi.brand) LIKE $%d", argIndex))
		args = append(args, "%"+strings.ToLower(req.Brand)+"%")
		argIndex++
	}
	
	if req.Model != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(vi.model) LIKE $%d", argIndex))
		args = append(args, "%"+strings.ToLower(req.Model)+"%")
		argIndex++
	}
	
	if req.Type != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(vi.type) LIKE $%d", argIndex))
		args = append(args, "%"+strings.ToLower(req.Type)+"%")
		argIndex++
	}
	
	if req.YearFrom > 0 {
		conditions = append(conditions, fmt.Sprintf("vi.production_year >= $%d", argIndex))
		args = append(args, req.YearFrom)
		argIndex++
	}
	
	if req.YearTo > 0 {
		conditions = append(conditions, fmt.Sprintf("vi.production_year <= $%d", argIndex))
		args = append(args, req.YearTo)
		argIndex++
	}
	
	if req.PriceFrom != nil {
		conditions = append(conditions, fmt.Sprintf("vi.estimated_selling_price >= $%d", argIndex))
		args = append(args, *req.PriceFrom)
		argIndex++
	}
	
	if req.PriceTo != nil {
		conditions = append(conditions, fmt.Sprintf("vi.estimated_selling_price <= $%d", argIndex))
		args = append(args, *req.PriceTo)
		argIndex++
	}
	
	if req.Status != "" {
		conditions = append(conditions, fmt.Sprintf("vi.status = $%d", argIndex))
		args = append(args, req.Status)
		argIndex++
	}
	
	if req.ConditionMin > 0 {
		conditions = append(conditions, fmt.Sprintf("vi.condition_rating >= $%d", argIndex))
		args = append(args, req.ConditionMin)
		argIndex++
	}
	
	if req.ConditionMax > 0 {
		conditions = append(conditions, fmt.Sprintf("vi.condition_rating <= $%d", argIndex))
		args = append(args, req.ConditionMax)
		argIndex++
	}
	
	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	
	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM vehicle_inventory vi
		LEFT JOIN vehicle_purchases vp ON vi.vehicle_purchase_id = vp.purchase_id
		LEFT JOIN customers c ON vp.customer_id = c.customer_id
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get vehicle count: %w", err)
	}
	
	// Get vehicles with pagination
	offset := (req.Page - 1) * req.PerPage
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	
	query := fmt.Sprintf(`
		SELECT vi.inventory_id, vi.vehicle_purchase_id, vi.plate_number, vi.brand, 
			   vi.model, vi.type, vi.production_year, vi.chassis_number, vi.engine_number,
			   vi.color, vi.mileage, vi.condition_rating, vi.purchase_price, 
			   vi.estimated_selling_price, vi.actual_selling_price, vi.status, 
			   vi.condition_notes, vi.selling_date, vi.profit_margin,
			   vi.created_at, vi.updated_at, vi.deleted_at, vi.created_by,
			   vp.purchase_date, vp.purchase_price as original_purchase_price,
			   c.name as "seller.name", c.phone as "seller.phone"
		FROM vehicle_inventory vi
		LEFT JOIN vehicle_purchases vp ON vi.vehicle_purchase_id = vp.purchase_id
		LEFT JOIN customers c ON vp.customer_id = c.customer_id
		%s
		ORDER BY vi.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	
	args = append(args, req.PerPage, offset)
	
	var vehicles []*models.VehicleInventory
	err = r.db.Select(&vehicles, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search vehicles: %w", err)
	}
	
	return vehicles, total, nil
}

func (r *vehicleTradingRepository) UpdateVehicleInventory(id int64, inventory *models.UpdateVehicleInventoryRequest) error {
	var setParts []string
	var args []interface{}
	argIndex := 1
	
	if inventory.EstimatedSellingPrice != nil {
		setParts = append(setParts, fmt.Sprintf("estimated_selling_price = $%d", argIndex))
		args = append(args, *inventory.EstimatedSellingPrice)
		argIndex++
	}
	
	if inventory.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *inventory.Status)
		argIndex++
	}
	
	if inventory.ConditionRating != nil {
		setParts = append(setParts, fmt.Sprintf("condition_rating = $%d", argIndex))
		args = append(args, *inventory.ConditionRating)
		argIndex++
	}
	
	if inventory.ConditionNotes != nil {
		setParts = append(setParts, fmt.Sprintf("condition_notes = $%d", argIndex))
		args = append(args, *inventory.ConditionNotes)
		argIndex++
	}
	
	if inventory.Mileage != nil {
		setParts = append(setParts, fmt.Sprintf("mileage = $%d", argIndex))
		args = append(args, *inventory.Mileage)
		argIndex++
	}
	
	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}
	
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)
	
	query := fmt.Sprintf(`
		UPDATE vehicle_inventory 
		SET %s
		WHERE inventory_id = $%d AND deleted_at IS NULL
	`, strings.Join(setParts, ", "), argIndex)
	
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update vehicle inventory: %w", err)
	}
	
	return nil
}

func (r *vehicleTradingRepository) UpdateSellingPrice(id int64, price decimal.Decimal) error {
	query := `
		UPDATE vehicle_inventory 
		SET estimated_selling_price = $2, updated_at = CURRENT_TIMESTAMP
		WHERE inventory_id = $1 AND deleted_at IS NULL
	`
	
	_, err := r.db.Exec(query, id, price)
	if err != nil {
		return fmt.Errorf("failed to update selling price: %w", err)
	}
	
	return nil
}

func (r *vehicleTradingRepository) MarkAsSold(id int64, saleData *models.VehicleSale) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Update inventory status
	_, err = tx.Exec(`
		UPDATE vehicle_inventory 
		SET status = 'Sold', actual_selling_price = $2, selling_date = $3,
			profit_margin = $2 - purchase_price, updated_at = CURRENT_TIMESTAMP
		WHERE inventory_id = $1 AND deleted_at IS NULL
	`, id, saleData.SellingPrice, saleData.SaleDate)
	if err != nil {
		return fmt.Errorf("failed to update inventory status: %w", err)
	}
	
	// Create sale record
	_, err = tx.Exec(`
		INSERT INTO vehicle_sales (inventory_id, customer_id, sales_person_id, outlet_id,
			sale_date, selling_price, commission_rate, commission_amount, payment_type,
			down_payment, financing_amount, financing_bank, financing_term_months,
			status, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, id, saleData.CustomerID, saleData.SalesPersonID, saleData.OutletID,
		saleData.SaleDate, saleData.SellingPrice, saleData.CommissionRate,
		saleData.CommissionAmount, saleData.PaymentType, saleData.DownPayment,
		saleData.FinancingAmount, saleData.FinancingBank, saleData.FinancingTermMonths,
		saleData.Status, saleData.Notes, saleData.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to create sale record: %w", err)
	}
	
	return tx.Commit()
}

func (r *vehicleTradingRepository) SoftDeleteVehicleInventory(id int64, deletedBy int64) error {
	query := `
		UPDATE vehicle_inventory 
		SET deleted_at = CURRENT_TIMESTAMP, created_by = $2
		WHERE inventory_id = $1 AND deleted_at IS NULL
	`
	
	_, err := r.db.Exec(query, id, deletedBy)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle inventory: %w", err)
	}
	
	return nil
}

// Implement remaining methods...
// (CreateVehicleSale, CreateVehiclePhoto, CreateVehicleAssessment, etc.)

// Placeholder implementations - these would be fully implemented similarly
func (r *vehicleTradingRepository) CreateVehicleSale(sale *models.VehicleSale) error {
	// Implementation similar to CreateVehiclePurchase
	return nil
}

func (r *vehicleTradingRepository) GetVehicleSaleByID(id int64) (*models.VehicleSale, error) {
	// Implementation similar to GetVehiclePurchaseByID
	return nil, nil
}

func (r *vehicleTradingRepository) GetVehicleSales(offset, limit int) ([]*models.VehicleSale, int64, error) {
	// Implementation similar to GetVehiclePurchases
	return nil, 0, nil
}

func (r *vehicleTradingRepository) UpdateVehicleSale(id int64, sale *models.VehicleSale) error {
	// Implementation similar to UpdateVehiclePurchase
	return nil
}

func (r *vehicleTradingRepository) SoftDeleteVehicleSale(id int64, deletedBy int64) error {
	// Implementation similar to SoftDeleteVehiclePurchase
	return nil
}

func (r *vehicleTradingRepository) CreateVehiclePhoto(photo *models.VehiclePhoto) error {
	// Implementation for photo management
	return nil
}

func (r *vehicleTradingRepository) GetVehiclePhotos(inventoryID int64) ([]*models.VehiclePhoto, error) {
	// Implementation for getting vehicle photos
	return nil, nil
}

func (r *vehicleTradingRepository) UpdateVehiclePhoto(id int64, photo *models.VehiclePhoto) error {
	// Implementation for updating photos
	return nil
}

func (r *vehicleTradingRepository) SoftDeleteVehiclePhoto(id int64, deletedBy int64) error {
	// Implementation for deleting photos
	return nil
}

func (r *vehicleTradingRepository) CreateVehicleAssessment(assessment *models.VehicleConditionAssessment) error {
	// Implementation for creating assessments
	return nil
}

func (r *vehicleTradingRepository) GetVehicleAssessments(inventoryID int64) ([]*models.VehicleConditionAssessment, error) {
	// Implementation for getting assessments
	return nil, nil
}

func (r *vehicleTradingRepository) UpdateVehicleAssessment(id int64, assessment *models.VehicleConditionAssessment) error {
	// Implementation for updating assessments
	return nil
}

func (r *vehicleTradingRepository) CreateSalesCommission(commission *models.SalesCommission) error {
	// Implementation for commission tracking
	return nil
}

func (r *vehicleTradingRepository) GetSalesCommissions(salesPersonID int64, startDate, endDate time.Time) ([]*models.SalesCommission, error) {
	// Implementation for getting commissions
	return nil, nil
}

func (r *vehicleTradingRepository) UpdateCommissionPaymentStatus(id int64, status string, paymentDate *time.Time) error {
	// Implementation for updating commission status
	return nil
}

func (r *vehicleTradingRepository) GetProfitAnalysis(startDate, endDate time.Time, outletID *int64) (*models.ProfitAnalysis, error) {
	// Implementation for profit analysis
	return nil, nil
}

func (r *vehicleTradingRepository) GetInventoryAgingReport(outletID *int64) ([]*models.InventoryAgingReport, error) {
	// Implementation for aging report
	return nil, nil
}

func (r *vehicleTradingRepository) GetSalesPerformanceReport(startDate, endDate time.Time, outletID *int64) ([]*models.SalesPerformanceReport, error) {
	// Implementation for sales performance
	return nil, nil
}
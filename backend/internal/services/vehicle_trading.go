package services

import (
	"fmt"
	"mime/multipart"
	"time"

	"flutter-bengkel/internal/models"
	"flutter-bengkel/internal/repositories"

	"github.com/shopspring/decimal"
)

// VehicleTradingService interface defines vehicle trading business logic
type VehicleTradingService interface {
	// Vehicle Purchase
	CreateVehiclePurchase(req *models.CreateVehiclePurchaseRequest, userID, outletID int64) (*models.VehiclePurchase, error)
	GetVehiclePurchaseByID(id int64) (*models.VehiclePurchase, error)
	GetVehiclePurchases(offset, limit int) ([]*models.VehiclePurchase, int64, error)
	UpdateVehiclePurchase(id int64, req *models.CreateVehiclePurchaseRequest) (*models.VehiclePurchase, error)
	DeleteVehiclePurchase(id int64, deletedBy int64) error

	// Vehicle Inventory
	GetVehicleInventoryByID(id int64) (*models.VehicleInventory, error)
	GetAvailableVehicles() ([]*models.VehicleInventory, error)
	SearchVehicles(req *models.VehicleSearchRequest) ([]*models.VehicleInventory, int64, error)
	UpdateVehicleInventory(id int64, req *models.UpdateVehicleInventoryRequest) (*models.VehicleInventory, error)
	UpdateSellingPrice(id int64, price decimal.Decimal) error
	DeleteVehicleInventory(id int64, deletedBy int64) error

	// Vehicle Sales
	CreateVehicleSale(req *models.CreateVehicleSaleRequest, salesPersonID, outletID int64) (*models.VehicleSale, error)
	GetVehicleSaleByID(id int64) (*models.VehicleSale, error)
	GetVehicleSales(offset, limit int) ([]*models.VehicleSale, int64, error)
	UpdateVehicleSale(id int64, req *models.CreateVehicleSaleRequest) (*models.VehicleSale, error)
	DeleteVehicleSale(id int64, deletedBy int64) error

	// Vehicle Photos
	UploadVehiclePhotos(inventoryID int64, files []*multipart.FileHeader, userID int64) ([]*models.VehiclePhoto, error)
	GetVehiclePhotos(inventoryID int64) ([]*models.VehiclePhoto, error)
	DeleteVehiclePhoto(id int64, deletedBy int64) error

	// Vehicle Assessments
	CreateVehicleAssessment(inventoryID int64, req *models.VehicleConditionAssessment, userID int64) (*models.VehicleConditionAssessment, error)
	GetVehicleAssessments(inventoryID int64) ([]*models.VehicleConditionAssessment, error)

	// Commission Management
	CalculateCommission(saleAmount decimal.Decimal, salesPersonID int64) (decimal.Decimal, decimal.Decimal, error)
	GetSalesCommissions(salesPersonID int64, startDate, endDate time.Time) ([]*models.SalesCommission, error)
	PayCommission(commissionID int64, paymentDate time.Time) error

	// Analytics & Reports
	GetProfitAnalysis(startDate, endDate time.Time, outletID *int64) (*models.ProfitAnalysis, error)
	GetInventoryAgingReport(outletID *int64) ([]*models.InventoryAgingReport, error)
	GetSalesPerformanceReport(startDate, endDate time.Time, outletID *int64) ([]*models.SalesPerformanceReport, error)
}

type vehicleTradingService struct {
	repos *repositories.Repositories
}

// NewVehicleTradingService creates a new vehicle trading service
func NewVehicleTradingService(repos *repositories.Repositories) VehicleTradingService {
	return &vehicleTradingService{repos: repos}
}

// Vehicle Purchase Operations
func (s *vehicleTradingService) CreateVehiclePurchase(req *models.CreateVehiclePurchaseRequest, userID, outletID int64) (*models.VehiclePurchase, error) {
	// Create purchase record
	purchase := &models.VehiclePurchase{
		CustomerID:    req.CustomerID,
		OutletID:      outletID,
		PurchaseDate:  time.Now(),
		PurchasePrice: req.PurchasePrice,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
		Status:        "completed",
		BaseModelWithSoftDelete: models.BaseModelWithSoftDelete{
			CreatedBy: &userID,
		},
	}

	err := s.repos.VehicleTrading.CreateVehiclePurchase(purchase)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle purchase: %w", err)
	}

	// Create inventory record
	inventory := &models.VehicleInventory{
		VehiclePurchaseID:     purchase.PurchaseID,
		PlateNumber:           req.VehicleDetails.PlateNumber,
		Brand:                 req.VehicleDetails.Brand,
		Model:                 req.VehicleDetails.Model,
		Type:                  req.VehicleDetails.Type,
		ProductionYear:        req.VehicleDetails.ProductionYear,
		ChassisNumber:         req.VehicleDetails.ChassisNumber,
		EngineNumber:          req.VehicleDetails.EngineNumber,
		Color:                 req.VehicleDetails.Color,
		Mileage:               req.VehicleDetails.Mileage,
		ConditionRating:       req.VehicleDetails.ConditionRating,
		PurchasePrice:         req.PurchasePrice,
		EstimatedSellingPrice: req.VehicleDetails.EstimatedSellingPrice,
		Status:                "Available",
		ConditionNotes:        req.VehicleDetails.ConditionNotes,
		BaseModelWithSoftDelete: models.BaseModelWithSoftDelete{
			CreatedBy: &userID,
		},
	}

	err = s.repos.VehicleTrading.CreateVehicleInventory(inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle inventory: %w", err)
	}

	return purchase, nil
}

func (s *vehicleTradingService) GetVehiclePurchaseByID(id int64) (*models.VehiclePurchase, error) {
	return s.repos.VehicleTrading.GetVehiclePurchaseByID(id)
}

func (s *vehicleTradingService) GetVehiclePurchases(offset, limit int) ([]*models.VehiclePurchase, int64, error) {
	return s.repos.VehicleTrading.GetVehiclePurchases(offset, limit)
}

func (s *vehicleTradingService) UpdateVehiclePurchase(id int64, req *models.CreateVehiclePurchaseRequest) (*models.VehiclePurchase, error) {
	// Get existing purchase
	existingPurchase, err := s.repos.VehicleTrading.GetVehiclePurchaseByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingPurchase.PurchasePrice = req.PurchasePrice
	existingPurchase.PaymentMethod = req.PaymentMethod
	existingPurchase.Notes = req.Notes

	err = s.repos.VehicleTrading.UpdateVehiclePurchase(id, existingPurchase)
	if err != nil {
		return nil, err
	}

	return s.repos.VehicleTrading.GetVehiclePurchaseByID(id)
}

func (s *vehicleTradingService) DeleteVehiclePurchase(id int64, deletedBy int64) error {
	return s.repos.VehicleTrading.SoftDeleteVehiclePurchase(id, deletedBy)
}

// Vehicle Inventory Operations
func (s *vehicleTradingService) GetVehicleInventoryByID(id int64) (*models.VehicleInventory, error) {
	return s.repos.VehicleTrading.GetVehicleInventoryByID(id)
}

func (s *vehicleTradingService) GetAvailableVehicles() ([]*models.VehicleInventory, error) {
	return s.repos.VehicleTrading.GetAvailableVehicles()
}

func (s *vehicleTradingService) SearchVehicles(req *models.VehicleSearchRequest) ([]*models.VehicleInventory, int64, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}

	return s.repos.VehicleTrading.SearchVehicles(req)
}

func (s *vehicleTradingService) UpdateVehicleInventory(id int64, req *models.UpdateVehicleInventoryRequest) (*models.VehicleInventory, error) {
	err := s.repos.VehicleTrading.UpdateVehicleInventory(id, req)
	if err != nil {
		return nil, err
	}

	return s.repos.VehicleTrading.GetVehicleInventoryByID(id)
}

func (s *vehicleTradingService) UpdateSellingPrice(id int64, price decimal.Decimal) error {
	return s.repos.VehicleTrading.UpdateSellingPrice(id, price)
}

func (s *vehicleTradingService) DeleteVehicleInventory(id int64, deletedBy int64) error {
	return s.repos.VehicleTrading.SoftDeleteVehicleInventory(id, deletedBy)
}

// Vehicle Sales Operations
func (s *vehicleTradingService) CreateVehicleSale(req *models.CreateVehicleSaleRequest, salesPersonID, outletID int64) (*models.VehicleSale, error) {
	// Get vehicle inventory to check availability
	inventory, err := s.repos.VehicleTrading.GetVehicleInventoryByID(req.InventoryID)
	if err != nil {
		return nil, fmt.Errorf("vehicle not found: %w", err)
	}

	if inventory.Status != "Available" {
		return nil, fmt.Errorf("vehicle is not available for sale")
	}

	// Calculate commission
	commissionRate, commissionAmount, err := s.CalculateCommission(req.SellingPrice, salesPersonID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate commission: %w", err)
	}

	// Create sale record
	sale := &models.VehicleSale{
		InventoryID:         req.InventoryID,
		CustomerID:          req.CustomerID,
		SalesPersonID:       salesPersonID,
		OutletID:            outletID,
		SaleDate:            time.Now(),
		SellingPrice:        req.SellingPrice,
		CommissionRate:      commissionRate,
		CommissionAmount:    commissionAmount,
		PaymentType:         req.PaymentType,
		DownPayment:         req.DownPayment,
		FinancingAmount:     req.FinancingAmount,
		FinancingBank:       req.FinancingBank,
		FinancingTermMonths: req.FinancingTermMonths,
		Status:              "Completed",
		Notes:               req.Notes,
		BaseModelWithSoftDelete: models.BaseModelWithSoftDelete{
			CreatedBy: &salesPersonID,
		},
	}

	// Mark vehicle as sold and create sale record in transaction
	err = s.repos.VehicleTrading.MarkAsSold(req.InventoryID, sale)
	if err != nil {
		return nil, fmt.Errorf("failed to record vehicle sale: %w", err)
	}

	// Create commission record
	commission := &models.SalesCommission{
		SaleID:           sale.SaleID,
		SalesPersonID:    salesPersonID,
		CommissionRate:   commissionRate,
		CommissionAmount: commissionAmount,
		PaymentStatus:    "pending",
		BaseModelWithSoftDelete: models.BaseModelWithSoftDelete{
			CreatedBy: &salesPersonID,
		},
	}

	err = s.repos.VehicleTrading.CreateSalesCommission(commission)
	if err != nil {
		// Log error but don't fail the sale
		// This should be handled better in production with proper logging
		fmt.Printf("Warning: failed to create commission record: %v\n", err)
	}

	return sale, nil
}

func (s *vehicleTradingService) GetVehicleSaleByID(id int64) (*models.VehicleSale, error) {
	return s.repos.VehicleTrading.GetVehicleSaleByID(id)
}

func (s *vehicleTradingService) GetVehicleSales(offset, limit int) ([]*models.VehicleSale, int64, error) {
	return s.repos.VehicleTrading.GetVehicleSales(offset, limit)
}

func (s *vehicleTradingService) UpdateVehicleSale(id int64, req *models.CreateVehicleSaleRequest) (*models.VehicleSale, error) {
	// Implementation similar to purchase update
	return nil, fmt.Errorf("not implemented")
}

func (s *vehicleTradingService) DeleteVehicleSale(id int64, deletedBy int64) error {
	return s.repos.VehicleTrading.SoftDeleteVehicleSale(id, deletedBy)
}

// Commission Calculation
func (s *vehicleTradingService) CalculateCommission(saleAmount decimal.Decimal, salesPersonID int64) (decimal.Decimal, decimal.Decimal, error) {
	// Default commission rate (can be configured per salesperson)
	defaultRate := decimal.NewFromFloat(0.025) // 2.5%

	// TODO: Get salesperson-specific commission rate from database
	// For now, use default rate

	commissionAmount := saleAmount.Mul(defaultRate)
	return defaultRate, commissionAmount, nil
}

func (s *vehicleTradingService) GetSalesCommissions(salesPersonID int64, startDate, endDate time.Time) ([]*models.SalesCommission, error) {
	return s.repos.VehicleTrading.GetSalesCommissions(salesPersonID, startDate, endDate)
}

func (s *vehicleTradingService) PayCommission(commissionID int64, paymentDate time.Time) error {
	return s.repos.VehicleTrading.UpdateCommissionPaymentStatus(commissionID, "paid", &paymentDate)
}

// Vehicle Photos
func (s *vehicleTradingService) UploadVehiclePhotos(inventoryID int64, files []*multipart.FileHeader, userID int64) ([]*models.VehiclePhoto, error) {
	var photos []*models.VehiclePhoto

	for i, file := range files {
		// TODO: Implement actual file upload to storage (S3, local storage, etc.)
		// For now, we'll create mock URLs
		photoURL := fmt.Sprintf("/uploads/vehicles/%d/%s", inventoryID, file.Filename)

		photo := &models.VehiclePhoto{
			InventoryID: inventoryID,
			PhotoURL:    photoURL,
			PhotoType:   s.determinePhotoType(file.Filename),
			Description: fmt.Sprintf("Vehicle photo %d", i+1),
			IsPrimary:   i == 0, // First photo is primary
			SortOrder:   i,
			BaseModelWithSoftDelete: models.BaseModelWithSoftDelete{
				CreatedBy: &userID,
			},
		}

		err := s.repos.VehicleTrading.CreateVehiclePhoto(photo)
		if err != nil {
			return nil, fmt.Errorf("failed to save photo record: %w", err)
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

func (s *vehicleTradingService) determinePhotoType(filename string) string {
	// Simple logic to determine photo type based on filename
	// In production, this could be more sophisticated
	return "other"
}

func (s *vehicleTradingService) GetVehiclePhotos(inventoryID int64) ([]*models.VehiclePhoto, error) {
	return s.repos.VehicleTrading.GetVehiclePhotos(inventoryID)
}

func (s *vehicleTradingService) DeleteVehiclePhoto(id int64, deletedBy int64) error {
	return s.repos.VehicleTrading.SoftDeleteVehiclePhoto(id, deletedBy)
}

// Vehicle Assessments
func (s *vehicleTradingService) CreateVehicleAssessment(inventoryID int64, req *models.VehicleConditionAssessment, userID int64) (*models.VehicleConditionAssessment, error) {
	req.InventoryID = inventoryID
	req.AssessorID = userID
	req.AssessmentDate = time.Now()
	req.CreatedBy = &userID

	err := s.repos.VehicleTrading.CreateVehicleAssessment(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle assessment: %w", err)
	}

	return req, nil
}

func (s *vehicleTradingService) GetVehicleAssessments(inventoryID int64) ([]*models.VehicleConditionAssessment, error) {
	return s.repos.VehicleTrading.GetVehicleAssessments(inventoryID)
}

// Analytics & Reports
func (s *vehicleTradingService) GetProfitAnalysis(startDate, endDate time.Time, outletID *int64) (*models.ProfitAnalysis, error) {
	return s.repos.VehicleTrading.GetProfitAnalysis(startDate, endDate, outletID)
}

func (s *vehicleTradingService) GetInventoryAgingReport(outletID *int64) ([]*models.InventoryAgingReport, error) {
	return s.repos.VehicleTrading.GetInventoryAgingReport(outletID)
}

func (s *vehicleTradingService) GetSalesPerformanceReport(startDate, endDate time.Time, outletID *int64) ([]*models.SalesPerformanceReport, error) {
	return s.repos.VehicleTrading.GetSalesPerformanceReport(startDate, endDate, outletID)
}
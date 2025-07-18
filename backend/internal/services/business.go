package services

import (
	"errors"

	"flutter-bengkel/internal/models"
	"flutter-bengkel/internal/repositories"
)

// Customer Service
type CustomerService interface {
	Create(req *models.Customer) (*models.Customer, error)
	GetByID(id int64) (*models.Customer, error)
	Update(id int64, req *models.Customer) (*models.Customer, error)
	Delete(id int64) error
	List(page, limit int, search string) ([]models.Customer, *models.PaginationMeta, error)
}

type customerService struct {
	repos *repositories.Repositories
}

func NewCustomerService(repos *repositories.Repositories) CustomerService {
	return &customerService{repos: repos}
}

func (s *customerService) Create(req *models.Customer) (*models.Customer, error) {
	// Generate customer code if not provided
	if req.CustomerCode == "" {
		code, err := s.repos.Customer.GenerateCustomerCode()
		if err != nil {
			return nil, err
		}
		req.CustomerCode = code
	}

	// Check if customer code already exists
	if _, err := s.repos.Customer.GetByCustomerCode(req.CustomerCode); err == nil {
		return nil, errors.New("customer code already exists")
	}

	// Set defaults
	if req.CustomerType == "" {
		req.CustomerType = "individual"
	}
	req.IsActive = true
	req.LoyaltyPoints = 0

	if err := s.repos.Customer.Create(req); err != nil {
		return nil, err
	}

	return s.repos.Customer.GetByID(req.ID)
}

func (s *customerService) GetByID(id int64) (*models.Customer, error) {
	customer, err := s.repos.Customer.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get customer vehicles
	vehicles, err := s.repos.Vehicle.GetByCustomerID(id)
	if err == nil {
		customer.Vehicles = vehicles
	}

	return customer, nil
}

func (s *customerService) Update(id int64, req *models.Customer) (*models.Customer, error) {
	// Get existing customer
	existingCustomer, err := s.repos.Customer.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if customer code is being changed and already exists
	if req.CustomerCode != existingCustomer.CustomerCode {
		if _, err := s.repos.Customer.GetByCustomerCode(req.CustomerCode); err == nil {
			return nil, errors.New("customer code already exists")
		}
	}

	if err := s.repos.Customer.Update(id, req); err != nil {
		return nil, err
	}

	return s.repos.Customer.GetByID(id)
}

func (s *customerService) Delete(id int64) error {
	return s.repos.Customer.Delete(id)
}

func (s *customerService) List(page, limit int, search string) ([]models.Customer, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	customers, total, err := s.repos.Customer.List(offset, limit, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return customers, meta, nil
}

// Vehicle Service
type VehicleService interface {
	Create(req *models.CustomerVehicle) (*models.CustomerVehicle, error)
	GetByID(id int64) (*models.CustomerVehicle, error)
	Update(id int64, req *models.CustomerVehicle) (*models.CustomerVehicle, error)
	Delete(id int64) error
	List(page, limit int, customerID *int64, search string) ([]models.CustomerVehicle, *models.PaginationMeta, error)
	GetByCustomerID(customerID int64) ([]models.CustomerVehicle, error)
}

type vehicleService struct {
	repos *repositories.Repositories
}

func NewVehicleService(repos *repositories.Repositories) VehicleService {
	return &vehicleService{repos: repos}
}

func (s *vehicleService) Create(req *models.CustomerVehicle) (*models.CustomerVehicle, error) {
	// Check if vehicle number already exists
	if _, err := s.repos.Vehicle.GetByVehicleNumber(req.VehicleNumber); err == nil {
		return nil, errors.New("vehicle number already exists")
	}

	// Validate customer exists
	if _, err := s.repos.Customer.GetByID(req.CustomerID); err != nil {
		return nil, errors.New("customer not found")
	}

	// Set defaults
	req.IsActive = true
	if req.FuelType == "" {
		req.FuelType = "gasoline"
	}
	if req.Transmission == "" {
		req.Transmission = "manual"
	}

	if err := s.repos.Vehicle.Create(req); err != nil {
		return nil, err
	}

	return s.repos.Vehicle.GetByID(req.ID)
}

func (s *vehicleService) GetByID(id int64) (*models.CustomerVehicle, error) {
	return s.repos.Vehicle.GetByID(id)
}

func (s *vehicleService) Update(id int64, req *models.CustomerVehicle) (*models.CustomerVehicle, error) {
	// Get existing vehicle
	existingVehicle, err := s.repos.Vehicle.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if vehicle number is being changed and already exists
	if req.VehicleNumber != existingVehicle.VehicleNumber {
		if _, err := s.repos.Vehicle.GetByVehicleNumber(req.VehicleNumber); err == nil {
			return nil, errors.New("vehicle number already exists")
		}
	}

	if err := s.repos.Vehicle.Update(id, req); err != nil {
		return nil, err
	}

	return s.repos.Vehicle.GetByID(id)
}

func (s *vehicleService) Delete(id int64) error {
	return s.repos.Vehicle.Delete(id)
}

func (s *vehicleService) List(page, limit int, customerID *int64, search string) ([]models.CustomerVehicle, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	vehicles, total, err := s.repos.Vehicle.List(offset, limit, customerID, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return vehicles, meta, nil
}

func (s *vehicleService) GetByCustomerID(customerID int64) ([]models.CustomerVehicle, error) {
	return s.repos.Vehicle.GetByCustomerID(customerID)
}

// Service Service
type ServiceService interface {
	Create(req *models.Service) (*models.Service, error)
	GetByID(id int64) (*models.Service, error)
	Update(id int64, req *models.Service) (*models.Service, error)
	Delete(id int64) error
	List(page, limit int, categoryID *int64, search string) ([]models.Service, *models.PaginationMeta, error)
	ListCategories() ([]models.ServiceCategory, error)
}

type serviceService struct {
	repos *repositories.Repositories
}

func NewServiceService(repos *repositories.Repositories) ServiceService {
	return &serviceService{repos: repos}
}

func (s *serviceService) Create(req *models.Service) (*models.Service, error) {
	// Generate service code if not provided
	if req.ServiceCode == "" {
		code, err := s.repos.Service.GenerateServiceCode()
		if err != nil {
			return nil, err
		}
		req.ServiceCode = code
	}

	// Check if service code already exists
	if _, err := s.repos.Service.GetByServiceCode(req.ServiceCode); err == nil {
		return nil, errors.New("service code already exists")
	}

	// Set defaults
	req.IsActive = true

	if err := s.repos.Service.Create(req); err != nil {
		return nil, err
	}

	return s.repos.Service.GetByID(req.ID)
}

func (s *serviceService) GetByID(id int64) (*models.Service, error) {
	return s.repos.Service.GetByID(id)
}

func (s *serviceService) Update(id int64, req *models.Service) (*models.Service, error) {
	// Get existing service
	existingService, err := s.repos.Service.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if service code is being changed and already exists
	if req.ServiceCode != existingService.ServiceCode {
		if _, err := s.repos.Service.GetByServiceCode(req.ServiceCode); err == nil {
			return nil, errors.New("service code already exists")
		}
	}

	if err := s.repos.Service.Update(id, req); err != nil {
		return nil, err
	}

	return s.repos.Service.GetByID(id)
}

func (s *serviceService) Delete(id int64) error {
	return s.repos.Service.Delete(id)
}

func (s *serviceService) List(page, limit int, categoryID *int64, search string) ([]models.Service, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	services, total, err := s.repos.Service.List(offset, limit, categoryID, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return services, meta, nil
}

func (s *serviceService) ListCategories() ([]models.ServiceCategory, error) {
	return s.repos.Service.ListCategories()
}

// Product Service
type ProductService interface {
	Create(req *models.Product) (*models.Product, error)
	GetByID(id int64) (*models.Product, error)
	Update(id int64, req *models.Product) (*models.Product, error)
	Delete(id int64) error
	List(page, limit int, categoryID *int64, supplierID *int64, search string) ([]models.Product, *models.PaginationMeta, error)
	UpdateStock(id int64, quantity int, operation string) error
	GetLowStockProducts(outletID *int64) ([]models.Product, error)
	ListCategories() ([]models.Category, error)
	ListSuppliers() ([]models.Supplier, error)
	ListUnitTypes() ([]models.UnitType, error)
}

type productService struct {
	repos *repositories.Repositories
}

func NewProductService(repos *repositories.Repositories) ProductService {
	return &productService{repos: repos}
}

func (s *productService) Create(req *models.Product) (*models.Product, error) {
	// Generate product code if not provided
	if req.ProductCode == "" {
		code, err := s.repos.Product.GenerateProductCode()
		if err != nil {
			return nil, err
		}
		req.ProductCode = code
	}

	// Check if product code already exists
	if _, err := s.repos.Product.GetByProductCode(req.ProductCode); err == nil {
		return nil, errors.New("product code already exists")
	}

	// Set defaults
	req.IsActive = true

	if err := s.repos.Product.Create(req); err != nil {
		return nil, err
	}

	return s.repos.Product.GetByID(req.ID)
}

func (s *productService) GetByID(id int64) (*models.Product, error) {
	return s.repos.Product.GetByID(id)
}

func (s *productService) Update(id int64, req *models.Product) (*models.Product, error) {
	// Get existing product
	existingProduct, err := s.repos.Product.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if product code is being changed and already exists
	if req.ProductCode != existingProduct.ProductCode {
		if _, err := s.repos.Product.GetByProductCode(req.ProductCode); err == nil {
			return nil, errors.New("product code already exists")
		}
	}

	if err := s.repos.Product.Update(id, req); err != nil {
		return nil, err
	}

	return s.repos.Product.GetByID(id)
}

func (s *productService) Delete(id int64) error {
	return s.repos.Product.Delete(id)
}

func (s *productService) List(page, limit int, categoryID *int64, supplierID *int64, search string) ([]models.Product, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	products, total, err := s.repos.Product.List(offset, limit, categoryID, supplierID, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return products, meta, nil
}

func (s *productService) UpdateStock(id int64, quantity int, operation string) error {
	if operation != "add" && operation != "subtract" {
		return errors.New("invalid operation: must be 'add' or 'subtract'")
	}

	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	return s.repos.Product.UpdateStock(id, quantity, operation)
}

func (s *productService) GetLowStockProducts(outletID *int64) ([]models.Product, error) {
	return s.repos.Product.GetLowStockProducts(outletID)
}

func (s *productService) ListCategories() ([]models.Category, error) {
	return s.repos.Product.ListCategories()
}

func (s *productService) ListSuppliers() ([]models.Supplier, error) {
	return s.repos.Product.ListSuppliers()
}

func (s *productService) ListUnitTypes() ([]models.UnitType, error) {
	return s.repos.Product.ListUnitTypes()
}
package repositories

import (
	"fmt"

	"flutter-bengkel/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

// Service Repository
type ServiceRepository interface {
	Create(service *models.Service) error
	GetByID(id uuid.UUID) (*models.Service, error)
	GetByServiceCode(code string) (*models.Service, error)
	Update(id uuid.UUID, service *models.Service) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Restore(id uuid.UUID) error
	PermanentDelete(id uuid.UUID) error
	List(offset, limit int, categoryID *uuid.UUID, search string, includeDeleted bool) ([]models.Service, int64, error)
	ListCategories() ([]models.ServiceCategory, error)
	GenerateServiceCode() (string, error)
}

type serviceRepository struct {
	db *sqlx.DB
}

func NewServiceRepository(db *sqlx.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(service *models.Service) error {
	query := `
		INSERT INTO services (service_code, name, description, category_id, 
							  standard_price, estimated_duration, is_active)
		VALUES (:service_code, :name, :description, :category_id, 
				:standard_price, :estimated_duration, :is_active)
	`
	
	result, err := r.db.NamedExec(query, service)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get service ID: %w", err)
	}
	
	service.ID = id
	return nil
}

func (r *serviceRepository) GetByID(id int64) (*models.Service, error) {
	query := `
		SELECT s.id, s.service_code, s.name, s.description, s.category_id, 
			   s.standard_price, s.estimated_duration, s.is_active, s.created_at, s.updated_at,
			   sc.id as "category.id", sc.name as "category.name", 
			   sc.description as "category.description"
		FROM services s
		LEFT JOIN service_categories sc ON s.category_id = sc.id
		WHERE s.id = ?
	`
	
	var service models.Service
	err := r.db.Get(&service, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	
	return &service, nil
}

func (r *serviceRepository) GetByServiceCode(code string) (*models.Service, error) {
	query := `
		SELECT s.id, s.service_code, s.name, s.description, s.category_id, 
			   s.standard_price, s.estimated_duration, s.is_active, s.created_at, s.updated_at,
			   sc.id as "category.id", sc.name as "category.name", 
			   sc.description as "category.description"
		FROM services s
		LEFT JOIN service_categories sc ON s.category_id = sc.id
		WHERE s.service_code = ?
	`
	
	var service models.Service
	err := r.db.Get(&service, query, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get service by code: %w", err)
	}
	
	return &service, nil
}

func (r *serviceRepository) Update(id int64, service *models.Service) error {
	query := `
		UPDATE services 
		SET name = :name, description = :description, category_id = :category_id, 
			standard_price = :standard_price, estimated_duration = :estimated_duration, 
			is_active = :is_active
		WHERE id = :id
	`
	
	service.ID = id
	_, err := r.db.NamedExec(query, service)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}
	
	return nil
}

func (r *serviceRepository) Delete(id int64) error {
	query := `UPDATE services SET is_active = false WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	
	return nil
}

func (r *serviceRepository) List(offset, limit int, categoryID *int64, search string) ([]models.Service, int64, error) {
	whereClause := "WHERE s.is_active = true"
	args := []interface{}{}
	
	if categoryID != nil {
		whereClause += " AND s.category_id = ?"
		args = append(args, *categoryID)
	}
	
	if search != "" {
		whereClause += " AND (s.name LIKE ? OR s.service_code LIKE ? OR s.description LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM services s 
		LEFT JOIN service_categories sc ON s.category_id = sc.id 
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count services: %w", err)
	}
	
	// Get services
	query := fmt.Sprintf(`
		SELECT s.id, s.service_code, s.name, s.description, s.category_id, 
			   s.standard_price, s.estimated_duration, s.is_active, s.created_at, s.updated_at,
			   sc.id as "category.id", sc.name as "category.name", 
			   sc.description as "category.description"
		FROM services s
		LEFT JOIN service_categories sc ON s.category_id = sc.id
		%s
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var services []models.Service
	err = r.db.Select(&services, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list services: %w", err)
	}
	
	return services, total, nil
}

func (r *serviceRepository) ListCategories() ([]models.ServiceCategory, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at 
		FROM service_categories 
		WHERE is_active = true 
		ORDER BY name
	`
	
	var categories []models.ServiceCategory
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list service categories: %w", err)
	}
	
	return categories, nil
}

func (r *serviceRepository) GenerateServiceCode() (string, error) {
	query := `SELECT COUNT(*) FROM services WHERE service_code LIKE 'SVC%'`
	
	var count int
	err := r.db.Get(&count, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate service code: %w", err)
	}
	
	return fmt.Sprintf("SVC%03d", count+1), nil
}

// Product Repository
type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uuid.UUID) (*models.Product, error)
	GetByProductCode(code string) (*models.Product, error)
	Update(id uuid.UUID, product *models.Product) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Restore(id uuid.UUID) error
	PermanentDelete(id uuid.UUID) error
	List(offset, limit int, categoryID *uuid.UUID, supplierID *uuid.UUID, search string, includeDeleted bool) ([]models.Product, int64, error)
	UpdateStock(id uuid.UUID, quantity int, operation string) error // operation: "add" or "subtract"
	GetLowStockProducts(outletID *uuid.UUID) ([]models.Product, error)
	ListCategories() ([]models.Category, error)
	ListSuppliers() ([]models.Supplier, error)
	ListUnitTypes() ([]models.UnitType, error)
	GenerateProductCode() (string, error)
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (product_code, name, description, category_id, unit_type_id, 
							  supplier_id, cost_price, selling_price, stock_quantity, 
							  min_stock_level, max_stock_level, has_serial_number, 
							  is_service, is_active)
		VALUES (:product_code, :name, :description, :category_id, :unit_type_id, 
				:supplier_id, :cost_price, :selling_price, :stock_quantity, 
				:min_stock_level, :max_stock_level, :has_serial_number, 
				:is_service, :is_active)
	`
	
	result, err := r.db.NamedExec(query, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get product ID: %w", err)
	}
	
	product.ID = id
	return nil
}

func (r *productRepository) GetByID(id int64) (*models.Product, error) {
	query := `
		SELECT p.id, p.product_code, p.name, p.description, p.category_id, p.unit_type_id,
			   p.supplier_id, p.cost_price, p.selling_price, p.stock_quantity, 
			   p.min_stock_level, p.max_stock_level, p.has_serial_number, 
			   p.is_service, p.is_active, p.created_at, p.updated_at,
			   c.id as "category.id", c.name as "category.name",
			   ut.id as "unit_type.id", ut.name as "unit_type.name", 
			   ut.abbreviation as "unit_type.abbreviation",
			   s.id as "supplier.id", s.supplier_code as "supplier.supplier_code",
			   s.name as "supplier.name"
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN unit_types ut ON p.unit_type_id = ut.id
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		WHERE p.id = ?
	`
	
	var product models.Product
	err := r.db.Get(&product, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	
	return &product, nil
}

func (r *productRepository) GetByProductCode(code string) (*models.Product, error) {
	query := `
		SELECT p.id, p.product_code, p.name, p.description, p.category_id, p.unit_type_id,
			   p.supplier_id, p.cost_price, p.selling_price, p.stock_quantity, 
			   p.min_stock_level, p.max_stock_level, p.has_serial_number, 
			   p.is_service, p.is_active, p.created_at, p.updated_at,
			   c.id as "category.id", c.name as "category.name",
			   ut.id as "unit_type.id", ut.name as "unit_type.name", 
			   ut.abbreviation as "unit_type.abbreviation",
			   s.id as "supplier.id", s.supplier_code as "supplier.supplier_code",
			   s.name as "supplier.name"
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN unit_types ut ON p.unit_type_id = ut.id
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		WHERE p.product_code = ?
	`
	
	var product models.Product
	err := r.db.Get(&product, query, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by code: %w", err)
	}
	
	return &product, nil
}

func (r *productRepository) Update(id int64, product *models.Product) error {
	query := `
		UPDATE products 
		SET name = :name, description = :description, category_id = :category_id, 
			unit_type_id = :unit_type_id, supplier_id = :supplier_id, 
			cost_price = :cost_price, selling_price = :selling_price, 
			min_stock_level = :min_stock_level, max_stock_level = :max_stock_level, 
			has_serial_number = :has_serial_number, is_service = :is_service, 
			is_active = :is_active
		WHERE id = :id
	`
	
	product.ID = id
	_, err := r.db.NamedExec(query, product)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	
	return nil
}

func (r *productRepository) Delete(id int64) error {
	query := `UPDATE products SET is_active = false WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	
	return nil
}

func (r *productRepository) List(offset, limit int, categoryID *int64, supplierID *int64, search string) ([]models.Product, int64, error) {
	whereClause := "WHERE p.is_active = true"
	args := []interface{}{}
	
	if categoryID != nil {
		whereClause += " AND p.category_id = ?"
		args = append(args, *categoryID)
	}
	
	if supplierID != nil {
		whereClause += " AND p.supplier_id = ?"
		args = append(args, *supplierID)
	}
	
	if search != "" {
		whereClause += " AND (p.name LIKE ? OR p.product_code LIKE ? OR p.description LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM products p 
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}
	
	// Get products
	query := fmt.Sprintf(`
		SELECT p.id, p.product_code, p.name, p.description, p.category_id, p.unit_type_id,
			   p.supplier_id, p.cost_price, p.selling_price, p.stock_quantity, 
			   p.min_stock_level, p.max_stock_level, p.has_serial_number, 
			   p.is_service, p.is_active, p.created_at, p.updated_at,
			   c.id as "category.id", c.name as "category.name",
			   ut.id as "unit_type.id", ut.name as "unit_type.name", 
			   ut.abbreviation as "unit_type.abbreviation",
			   s.id as "supplier.id", s.supplier_code as "supplier.supplier_code",
			   s.name as "supplier.name"
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN unit_types ut ON p.unit_type_id = ut.id
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		%s
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var products []models.Product
	err = r.db.Select(&products, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	
	return products, total, nil
}

func (r *productRepository) UpdateStock(id int64, quantity int, operation string) error {
	var query string
	
	switch operation {
	case "add":
		query = `UPDATE products SET stock_quantity = stock_quantity + ? WHERE id = ?`
	case "subtract":
		query = `UPDATE products SET stock_quantity = stock_quantity - ? WHERE id = ? AND stock_quantity >= ?`
	default:
		return fmt.Errorf("invalid operation: %s", operation)
	}
	
	if operation == "subtract" {
		_, err := r.db.Exec(query, quantity, id, quantity)
		if err != nil {
			return fmt.Errorf("failed to update stock (insufficient stock): %w", err)
		}
	} else {
		_, err := r.db.Exec(query, quantity, id)
		if err != nil {
			return fmt.Errorf("failed to update stock: %w", err)
		}
	}
	
	return nil
}

func (r *productRepository) GetLowStockProducts(outletID *int64) ([]models.Product, error) {
	query := `
		SELECT p.id, p.product_code, p.name, p.description, p.category_id, p.unit_type_id,
			   p.supplier_id, p.cost_price, p.selling_price, p.stock_quantity, 
			   p.min_stock_level, p.max_stock_level, p.has_serial_number, 
			   p.is_service, p.is_active, p.created_at, p.updated_at,
			   c.id as "category.id", c.name as "category.name",
			   ut.id as "unit_type.id", ut.name as "unit_type.name", 
			   ut.abbreviation as "unit_type.abbreviation",
			   s.id as "supplier.id", s.supplier_code as "supplier.supplier_code",
			   s.name as "supplier.name"
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		LEFT JOIN unit_types ut ON p.unit_type_id = ut.id
		LEFT JOIN suppliers s ON p.supplier_id = s.id
		WHERE p.is_active = true AND p.stock_quantity <= p.min_stock_level
		ORDER BY p.stock_quantity ASC
	`
	
	var products []models.Product
	err := r.db.Select(&products, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock products: %w", err)
	}
	
	return products, nil
}

func (r *productRepository) ListCategories() ([]models.Category, error) {
	query := `
		SELECT id, name, parent_id, description, is_active, created_at, updated_at 
		FROM categories 
		WHERE is_active = true 
		ORDER BY name
	`
	
	var categories []models.Category
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list product categories: %w", err)
	}
	
	return categories, nil
}

func (r *productRepository) ListSuppliers() ([]models.Supplier, error) {
	query := `
		SELECT id, supplier_code, name, email, phone, address, city, province, 
			   postal_code, contact_person, payment_terms, is_active, created_at, updated_at
		FROM suppliers 
		WHERE is_active = true 
		ORDER BY name
	`
	
	var suppliers []models.Supplier
	err := r.db.Select(&suppliers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}
	
	return suppliers, nil
}

func (r *productRepository) ListUnitTypes() ([]models.UnitType, error) {
	query := `
		SELECT id, name, abbreviation, description, created_at, updated_at 
		FROM unit_types 
		ORDER BY name
	`
	
	var unitTypes []models.UnitType
	err := r.db.Select(&unitTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list unit types: %w", err)
	}
	
	return unitTypes, nil
}

func (r *productRepository) GenerateProductCode() (string, error) {
	query := `SELECT COUNT(*) FROM products WHERE product_code LIKE 'PRD%'`
	
	var count int
	err := r.db.Get(&count, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate product code: %w", err)
	}
	
	return fmt.Sprintf("PRD%06d", count+1), nil
}
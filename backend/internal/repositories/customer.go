package repositories

import (
	"fmt"

	"flutter-bengkel/internal/models"

	"github.com/jmoiron/sqlx"
)

// Customer Repository
type CustomerRepository interface {
	Create(customer *models.Customer) error
	GetByID(id int64) (*models.Customer, error)
	GetByCustomerCode(code string) (*models.Customer, error)
	Update(id int64, customer *models.Customer) error
	Delete(id int64) error
	List(offset, limit int, search string) ([]models.Customer, int64, error)
	GenerateCustomerCode() (string, error)
}

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	query := `
		INSERT INTO customers (customer_code, name, email, phone, address, city, province, 
							  postal_code, date_of_birth, gender, customer_type, notes, is_active)
		VALUES (:customer_code, :name, :email, :phone, :address, :city, :province, 
				:postal_code, :date_of_birth, :gender, :customer_type, :notes, :is_active)
	`
	
	result, err := r.db.NamedExec(query, customer)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get customer ID: %w", err)
	}
	
	customer.ID = id
	return nil
}

func (r *customerRepository) GetByID(id int64) (*models.Customer, error) {
	query := `
		SELECT id, customer_code, name, email, phone, address, city, province, 
			   postal_code, date_of_birth, gender, customer_type, loyalty_points, notes, 
			   is_active, created_at, updated_at
		FROM customers 
		WHERE id = ?
	`
	
	var customer models.Customer
	err := r.db.Get(&customer, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	
	return &customer, nil
}

func (r *customerRepository) GetByCustomerCode(code string) (*models.Customer, error) {
	query := `
		SELECT id, customer_code, name, email, phone, address, city, province, 
			   postal_code, date_of_birth, gender, customer_type, loyalty_points, notes, 
			   is_active, created_at, updated_at
		FROM customers 
		WHERE customer_code = ?
	`
	
	var customer models.Customer
	err := r.db.Get(&customer, query, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer by code: %w", err)
	}
	
	return &customer, nil
}

func (r *customerRepository) Update(id int64, customer *models.Customer) error {
	query := `
		UPDATE customers 
		SET name = :name, email = :email, phone = :phone, address = :address, 
			city = :city, province = :province, postal_code = :postal_code, 
			date_of_birth = :date_of_birth, gender = :gender, customer_type = :customer_type, 
			notes = :notes, is_active = :is_active
		WHERE id = :id
	`
	
	customer.ID = id
	_, err := r.db.NamedExec(query, customer)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	
	return nil
}

func (r *customerRepository) Delete(id int64) error {
	query := `UPDATE customers SET is_active = false WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	
	return nil
}

func (r *customerRepository) List(offset, limit int, search string) ([]models.Customer, int64, error) {
	whereClause := "WHERE is_active = true"
	args := []interface{}{}
	
	if search != "" {
		whereClause += " AND (name LIKE ? OR customer_code LIKE ? OR phone LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customers %s", whereClause)
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}
	
	// Get customers
	query := fmt.Sprintf(`
		SELECT id, customer_code, name, email, phone, address, city, province, 
			   postal_code, date_of_birth, gender, customer_type, loyalty_points, notes, 
			   is_active, created_at, updated_at
		FROM customers %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var customers []models.Customer
	err = r.db.Select(&customers, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	
	return customers, total, nil
}

func (r *customerRepository) GenerateCustomerCode() (string, error) {
	query := `SELECT COUNT(*) FROM customers WHERE customer_code LIKE 'CUST%'`
	
	var count int
	err := r.db.Get(&count, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate customer code: %w", err)
	}
	
	return fmt.Sprintf("CUST%06d", count+1), nil
}

// Vehicle Repository
type VehicleRepository interface {
	Create(vehicle *models.CustomerVehicle) error
	GetByID(id int64) (*models.CustomerVehicle, error)
	GetByVehicleNumber(vehicleNumber string) (*models.CustomerVehicle, error)
	GetByCustomerID(customerID int64) ([]models.CustomerVehicle, error)
	Update(id int64, vehicle *models.CustomerVehicle) error
	Delete(id int64) error
	List(offset, limit int, customerID *int64, search string) ([]models.CustomerVehicle, int64, error)
}

type vehicleRepository struct {
	db *sqlx.DB
}

func NewVehicleRepository(db *sqlx.DB) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Create(vehicle *models.CustomerVehicle) error {
	query := `
		INSERT INTO customer_vehicles (customer_id, vehicle_number, brand, model, year, color, 
									   engine_number, chassis_number, fuel_type, transmission, 
									   mileage, last_service_date, next_service_date, 
									   insurance_expiry, registration_expiry, notes, is_active)
		VALUES (:customer_id, :vehicle_number, :brand, :model, :year, :color, 
				:engine_number, :chassis_number, :fuel_type, :transmission, 
				:mileage, :last_service_date, :next_service_date, 
				:insurance_expiry, :registration_expiry, :notes, :is_active)
	`
	
	result, err := r.db.NamedExec(query, vehicle)
	if err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get vehicle ID: %w", err)
	}
	
	vehicle.ID = id
	return nil
}

func (r *vehicleRepository) GetByID(id int64) (*models.CustomerVehicle, error) {
	query := `
		SELECT cv.id, cv.customer_id, cv.vehicle_number, cv.brand, cv.model, cv.year, cv.color,
			   cv.engine_number, cv.chassis_number, cv.fuel_type, cv.transmission, cv.mileage,
			   cv.last_service_date, cv.next_service_date, cv.insurance_expiry, cv.registration_expiry,
			   cv.notes, cv.is_active, cv.created_at, cv.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone"
		FROM customer_vehicles cv
		LEFT JOIN customers c ON cv.customer_id = c.id
		WHERE cv.id = ?
	`
	
	var vehicle models.CustomerVehicle
	err := r.db.Get(&vehicle, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}
	
	return &vehicle, nil
}

func (r *vehicleRepository) GetByVehicleNumber(vehicleNumber string) (*models.CustomerVehicle, error) {
	query := `
		SELECT cv.id, cv.customer_id, cv.vehicle_number, cv.brand, cv.model, cv.year, cv.color,
			   cv.engine_number, cv.chassis_number, cv.fuel_type, cv.transmission, cv.mileage,
			   cv.last_service_date, cv.next_service_date, cv.insurance_expiry, cv.registration_expiry,
			   cv.notes, cv.is_active, cv.created_at, cv.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone"
		FROM customer_vehicles cv
		LEFT JOIN customers c ON cv.customer_id = c.id
		WHERE cv.vehicle_number = ?
	`
	
	var vehicle models.CustomerVehicle
	err := r.db.Get(&vehicle, query, vehicleNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle by number: %w", err)
	}
	
	return &vehicle, nil
}

func (r *vehicleRepository) GetByCustomerID(customerID int64) ([]models.CustomerVehicle, error) {
	query := `
		SELECT id, customer_id, vehicle_number, brand, model, year, color,
			   engine_number, chassis_number, fuel_type, transmission, mileage,
			   last_service_date, next_service_date, insurance_expiry, registration_expiry,
			   notes, is_active, created_at, updated_at
		FROM customer_vehicles
		WHERE customer_id = ? AND is_active = true
		ORDER BY created_at DESC
	`
	
	var vehicles []models.CustomerVehicle
	err := r.db.Select(&vehicles, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles by customer: %w", err)
	}
	
	return vehicles, nil
}

func (r *vehicleRepository) Update(id int64, vehicle *models.CustomerVehicle) error {
	query := `
		UPDATE customer_vehicles 
		SET brand = :brand, model = :model, year = :year, color = :color,
			engine_number = :engine_number, chassis_number = :chassis_number, 
			fuel_type = :fuel_type, transmission = :transmission, mileage = :mileage,
			last_service_date = :last_service_date, next_service_date = :next_service_date,
			insurance_expiry = :insurance_expiry, registration_expiry = :registration_expiry,
			notes = :notes, is_active = :is_active
		WHERE id = :id
	`
	
	vehicle.ID = id
	_, err := r.db.NamedExec(query, vehicle)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	
	return nil
}

func (r *vehicleRepository) Delete(id int64) error {
	query := `UPDATE customer_vehicles SET is_active = false WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}
	
	return nil
}

func (r *vehicleRepository) List(offset, limit int, customerID *int64, search string) ([]models.CustomerVehicle, int64, error) {
	whereClause := "WHERE cv.is_active = true"
	args := []interface{}{}
	
	if customerID != nil {
		whereClause += " AND cv.customer_id = ?"
		args = append(args, *customerID)
	}
	
	if search != "" {
		whereClause += " AND (cv.vehicle_number LIKE ? OR cv.brand LIKE ? OR cv.model LIKE ? OR c.name LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM customer_vehicles cv 
		LEFT JOIN customers c ON cv.customer_id = c.id 
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
	}
	
	// Get vehicles
	query := fmt.Sprintf(`
		SELECT cv.id, cv.customer_id, cv.vehicle_number, cv.brand, cv.model, cv.year, cv.color,
			   cv.engine_number, cv.chassis_number, cv.fuel_type, cv.transmission, cv.mileage,
			   cv.last_service_date, cv.next_service_date, cv.insurance_expiry, cv.registration_expiry,
			   cv.notes, cv.is_active, cv.created_at, cv.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone"
		FROM customer_vehicles cv
		LEFT JOIN customers c ON cv.customer_id = c.id
		%s
		ORDER BY cv.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var vehicles []models.CustomerVehicle
	err = r.db.Select(&vehicles, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}
	
	return vehicles, total, nil
}
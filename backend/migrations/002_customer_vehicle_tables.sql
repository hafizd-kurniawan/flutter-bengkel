-- Customer & Vehicle Management Tables (PostgreSQL with Soft Delete)

-- Customers table
CREATE TABLE customers (
    customer_id BIGSERIAL PRIMARY KEY,
    customer_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    date_of_birth DATE,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    customer_type VARCHAR(20) CHECK (customer_type IN ('individual', 'corporate')) DEFAULT 'individual',
    loyalty_points INTEGER DEFAULT 0,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER
);

-- Customer vehicles table
CREATE TABLE customer_vehicles (
    vehicle_id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    vehicle_number VARCHAR(20) NOT NULL UNIQUE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    color VARCHAR(50),
    engine_number VARCHAR(100),
    chassis_number VARCHAR(100),
    fuel_type VARCHAR(20) CHECK (fuel_type IN ('gasoline', 'diesel', 'electric', 'hybrid')) DEFAULT 'gasoline',
    transmission VARCHAR(20) CHECK (transmission IN ('manual', 'automatic', 'cvt')) DEFAULT 'manual',
    mileage BIGINT DEFAULT 0,
    last_service_date DATE,
    next_service_date DATE,
    insurance_expiry DATE,
    registration_expiry DATE,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_customers_customer_code ON customers(customer_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_phone ON customers(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_email ON customers(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_is_active ON customers(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_customer_id ON customer_vehicles(customer_id);
CREATE INDEX idx_customer_vehicles_vehicle_number ON customer_vehicles(vehicle_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_brand_model ON customer_vehicles(brand, model) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_is_active ON customer_vehicles(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);
CREATE INDEX idx_customer_vehicles_deleted_at ON customer_vehicles(deleted_at);
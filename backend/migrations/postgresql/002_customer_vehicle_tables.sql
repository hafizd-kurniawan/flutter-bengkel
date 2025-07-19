-- Customer & Vehicle Management Tables (PostgreSQL)

-- Customers table with loyalty tracking and audit
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    customer_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    date_of_birth DATE,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    customer_type VARCHAR(50) DEFAULT 'regular' CHECK (customer_type IN ('regular', 'premium', 'vip')),
    loyalty_points INTEGER DEFAULT 0,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Customer vehicles with comprehensive details and audit
CREATE TABLE customer_vehicles (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    vehicle_number VARCHAR(20) NOT NULL UNIQUE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    color VARCHAR(50),
    engine_number VARCHAR(100),
    chassis_number VARCHAR(100),
    fuel_type VARCHAR(20) CHECK (fuel_type IN ('gasoline', 'diesel', 'electric', 'hybrid')),
    transmission VARCHAR(20) CHECK (transmission IN ('manual', 'automatic', 'cvt')),
    mileage BIGINT DEFAULT 0,
    last_service_date DATE,
    next_service_date DATE,
    insurance_expiry DATE,
    registration_expiry DATE,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Create indexes for better performance
CREATE INDEX idx_customers_customer_code ON customers(customer_code);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_is_active ON customers(is_active);
CREATE INDEX idx_customers_created_by ON customers(created_by);
CREATE INDEX idx_customer_vehicles_customer_id ON customer_vehicles(customer_id);
CREATE INDEX idx_customer_vehicles_vehicle_number ON customer_vehicles(vehicle_number);
CREATE INDEX idx_customer_vehicles_brand_model ON customer_vehicles(brand, model);
CREATE INDEX idx_customer_vehicles_is_active ON customer_vehicles(is_active);
CREATE INDEX idx_customer_vehicles_created_by ON customer_vehicles(created_by);

-- Create triggers for updated_at columns
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_customer_vehicles_updated_at BEFORE UPDATE ON customer_vehicles
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
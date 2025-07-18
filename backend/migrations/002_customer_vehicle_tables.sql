-- Customer & Vehicle Management Tables with PostgreSQL + Soft Delete

-- Define custom types
CREATE TYPE gender_type AS ENUM ('male', 'female', 'other');
CREATE TYPE customer_type AS ENUM ('individual', 'corporate');
CREATE TYPE fuel_type AS ENUM ('gasoline', 'diesel', 'electric', 'hybrid');
CREATE TYPE transmission_type AS ENUM ('manual', 'automatic', 'cvt');

-- Customers table
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    date_of_birth DATE,
    gender gender_type,
    customer_type customer_type DEFAULT 'individual',
    loyalty_points INT DEFAULT 0,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_by UUID NULL REFERENCES users(id)
);

-- Customer vehicles table
CREATE TABLE customer_vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    vehicle_number VARCHAR(20) NOT NULL UNIQUE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    color VARCHAR(50),
    engine_number VARCHAR(100),
    chassis_number VARCHAR(100),
    fuel_type fuel_type DEFAULT 'gasoline',
    transmission transmission_type DEFAULT 'manual',
    mileage BIGINT DEFAULT 0,
    last_service_date DATE,
    next_service_date DATE,
    insurance_expiry DATE,
    registration_expiry DATE,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    deleted_by UUID NULL REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- Create indexes with soft delete consideration
CREATE INDEX idx_customers_customer_code ON customers(customer_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_phone ON customers(phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_email ON customers(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_is_active ON customers(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

CREATE INDEX idx_customer_vehicles_customer_id ON customer_vehicles(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_vehicle_number ON customer_vehicles(vehicle_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_brand_model ON customer_vehicles(brand, model) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_is_active ON customer_vehicles(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_customer_vehicles_deleted_at ON customer_vehicles(deleted_at);

-- Create triggers for updated_at
CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customer_vehicles_updated_at BEFORE UPDATE ON customer_vehicles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
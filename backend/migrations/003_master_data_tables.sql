-- Master Data & Inventory Tables (PostgreSQL with Soft Delete)

-- Unit types for products
CREATE TABLE unit_types (
    unit_type_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    abbreviation VARCHAR(10) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER
);

-- Product categories
CREATE TABLE categories (
    category_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id BIGINT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (parent_id) REFERENCES categories(category_id) ON DELETE SET NULL
);

-- Suppliers
CREATE TABLE suppliers (
    supplier_id BIGSERIAL PRIMARY KEY,
    supplier_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    contact_person VARCHAR(255),
    payment_terms TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER
);

-- Products table with inventory tracking
CREATE TABLE products (
    product_id BIGSERIAL PRIMARY KEY,
    product_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id BIGINT NOT NULL,
    unit_type_id BIGINT NOT NULL,
    supplier_id BIGINT,
    cost_price DECIMAL(15,2) DEFAULT 0,
    selling_price DECIMAL(15,2) DEFAULT 0,
    stock_quantity INTEGER DEFAULT 0,
    min_stock_level INTEGER DEFAULT 0,
    max_stock_level INTEGER DEFAULT 0,
    has_serial_number BOOLEAN DEFAULT FALSE,
    is_service BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories(category_id),
    FOREIGN KEY (unit_type_id) REFERENCES unit_types(unit_type_id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(supplier_id) ON DELETE SET NULL
);

-- Product serial numbers for tracking
CREATE TABLE product_serial_numbers (
    serial_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    serial_number VARCHAR(255) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('available', 'sold', 'reserved', 'damaged')) DEFAULT 'available',
    purchase_date DATE,
    sale_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE,
    UNIQUE (product_id, serial_number)
);

-- Service categories
CREATE TABLE service_categories (
    service_category_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER
);

-- Services table
CREATE TABLE services (
    service_id BIGSERIAL PRIMARY KEY,
    service_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id BIGINT NOT NULL,
    standard_price DECIMAL(15,2) DEFAULT 0,
    estimated_duration INTEGER DEFAULT 0, -- in minutes
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (category_id) REFERENCES service_categories(service_category_id)
);

-- Create indexes
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_is_active ON categories(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

CREATE INDEX idx_suppliers_supplier_code ON suppliers(supplier_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_is_active ON suppliers(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_suppliers_deleted_at ON suppliers(deleted_at);

CREATE INDEX idx_products_product_code ON products(product_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_supplier_id ON products(supplier_id);
CREATE INDEX idx_products_is_active ON products(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_stock_level ON products(stock_quantity, min_stock_level) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_deleted_at ON products(deleted_at);

CREATE INDEX idx_product_serial_numbers_serial ON product_serial_numbers(serial_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_status ON product_serial_numbers(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_serial_numbers_deleted_at ON product_serial_numbers(deleted_at);

CREATE INDEX idx_services_service_code ON services(service_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_services_category_id ON services(category_id);
CREATE INDEX idx_services_is_active ON services(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_services_deleted_at ON services(deleted_at);

CREATE INDEX idx_service_categories_is_active ON service_categories(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_service_categories_deleted_at ON service_categories(deleted_at);

CREATE INDEX idx_unit_types_deleted_at ON unit_types(deleted_at);
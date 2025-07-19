-- Vehicle Trading Tables (PostgreSQL with Soft Delete)
-- Complete implementation for vehicle buying, selling, and inventory management

-- Vehicle purchase record when buying from customers
CREATE TABLE vehicle_purchases (
    purchase_id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    purchase_date DATE NOT NULL,
    purchase_price DECIMAL(15, 2) NOT NULL,
    payment_method VARCHAR(50) DEFAULT 'cash',
    notes TEXT,
    status VARCHAR(20) CHECK (status IN ('pending', 'completed', 'cancelled')) DEFAULT 'completed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id)
);

-- Vehicle inventory for trading
CREATE TABLE vehicle_inventory (
    inventory_id BIGSERIAL PRIMARY KEY,
    vehicle_purchase_id BIGINT NOT NULL,
    plate_number VARCHAR(20) UNIQUE NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL,
    production_year INTEGER NOT NULL,
    chassis_number VARCHAR(100) UNIQUE NOT NULL,
    engine_number VARCHAR(100) UNIQUE NOT NULL,
    color VARCHAR(50) NOT NULL,
    mileage INTEGER DEFAULT 0,
    condition_rating INTEGER CHECK (condition_rating >= 1 AND condition_rating <= 5) DEFAULT 5,
    purchase_price DECIMAL(15, 2) NOT NULL,
    estimated_selling_price DECIMAL(15, 2) NOT NULL,
    actual_selling_price DECIMAL(15, 2),
    status VARCHAR(20) CHECK (status IN ('Available', 'Reserved', 'Sold', 'Under_Maintenance')) DEFAULT 'Available',
    vehicle_photos TEXT[], -- JSON array of photo URLs
    condition_notes TEXT,
    selling_date DATE,
    profit_margin DECIMAL(15, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (vehicle_purchase_id) REFERENCES vehicle_purchases(purchase_id)
);

-- Vehicle sales record when selling to customers
CREATE TABLE vehicle_sales (
    sale_id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    sales_person_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    sale_date DATE NOT NULL,
    selling_price DECIMAL(15, 2) NOT NULL,
    commission_rate DECIMAL(5, 2) DEFAULT 0,
    commission_amount DECIMAL(15, 2) DEFAULT 0,
    payment_type VARCHAR(20) CHECK (payment_type IN ('cash', 'credit', 'trade_in', 'financing')) NOT NULL,
    down_payment DECIMAL(15, 2) DEFAULT 0,
    financing_amount DECIMAL(15, 2) DEFAULT 0,
    financing_bank VARCHAR(100),
    financing_term_months INTEGER,
    status VARCHAR(20) CHECK (status IN ('Pending', 'Completed', 'Cancelled')) DEFAULT 'Completed',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (inventory_id) REFERENCES vehicle_inventory(inventory_id),
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    FOREIGN KEY (sales_person_id) REFERENCES users(user_id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id)
);

-- Vehicle condition assessments
CREATE TABLE vehicle_condition_assessments (
    assessment_id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL,
    assessor_id BIGINT NOT NULL,
    assessment_date DATE NOT NULL,
    exterior_condition INTEGER CHECK (exterior_condition >= 1 AND exterior_condition <= 5),
    interior_condition INTEGER CHECK (interior_condition >= 1 AND interior_condition <= 5),
    engine_condition INTEGER CHECK (engine_condition >= 1 AND engine_condition <= 5),
    transmission_condition INTEGER CHECK (transmission_condition >= 1 AND transmission_condition <= 5),
    tire_condition INTEGER CHECK (tire_condition >= 1 AND tire_condition <= 5),
    electrical_condition INTEGER CHECK (electrical_condition >= 1 AND electrical_condition <= 5),
    overall_rating INTEGER CHECK (overall_rating >= 1 AND overall_rating <= 5),
    assessment_notes TEXT,
    recommended_repairs TEXT,
    estimated_repair_cost DECIMAL(15, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (inventory_id) REFERENCES vehicle_inventory(inventory_id),
    FOREIGN KEY (assessor_id) REFERENCES users(user_id)
);

-- Vehicle photos
CREATE TABLE vehicle_photos (
    photo_id BIGSERIAL PRIMARY KEY,
    inventory_id BIGINT NOT NULL,
    photo_url VARCHAR(500) NOT NULL,
    photo_type VARCHAR(50) CHECK (photo_type IN ('exterior_front', 'exterior_back', 'exterior_left', 'exterior_right', 'interior_front', 'interior_back', 'engine', 'dashboard', 'trunk', 'other')),
    description TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (inventory_id) REFERENCES vehicle_inventory(inventory_id)
);

-- Commission tracking
CREATE TABLE sales_commissions (
    commission_id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL,
    sales_person_id BIGINT NOT NULL,
    commission_rate DECIMAL(5, 2) NOT NULL,
    commission_amount DECIMAL(15, 2) NOT NULL,
    payment_status VARCHAR(20) CHECK (payment_status IN ('pending', 'paid', 'cancelled')) DEFAULT 'pending',
    payment_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (sale_id) REFERENCES vehicle_sales(sale_id),
    FOREIGN KEY (sales_person_id) REFERENCES users(user_id)
);

-- Create comprehensive indexes for performance
CREATE INDEX idx_vehicle_purchases_customer_id ON vehicle_purchases(customer_id);
CREATE INDEX idx_vehicle_purchases_outlet_id ON vehicle_purchases(outlet_id);
CREATE INDEX idx_vehicle_purchases_purchase_date ON vehicle_purchases(purchase_date);
CREATE INDEX idx_vehicle_purchases_status ON vehicle_purchases(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_purchases_deleted_at ON vehicle_purchases(deleted_at);

CREATE INDEX idx_vehicle_inventory_status ON vehicle_inventory(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_inventory_brand_model ON vehicle_inventory(brand, model) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_inventory_price_range ON vehicle_inventory(estimated_selling_price) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_inventory_production_year ON vehicle_inventory(production_year) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_inventory_condition_rating ON vehicle_inventory(condition_rating) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_inventory_plate_number ON vehicle_inventory(plate_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_inventory_deleted_at ON vehicle_inventory(deleted_at);

CREATE INDEX idx_vehicle_sales_customer_id ON vehicle_sales(customer_id);
CREATE INDEX idx_vehicle_sales_sales_person_id ON vehicle_sales(sales_person_id);
CREATE INDEX idx_vehicle_sales_outlet_id ON vehicle_sales(outlet_id);
CREATE INDEX idx_vehicle_sales_sale_date ON vehicle_sales(sale_date);
CREATE INDEX idx_vehicle_sales_status ON vehicle_sales(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_sales_deleted_at ON vehicle_sales(deleted_at);

CREATE INDEX idx_vehicle_condition_assessments_inventory_id ON vehicle_condition_assessments(inventory_id);
CREATE INDEX idx_vehicle_condition_assessments_assessor_id ON vehicle_condition_assessments(assessor_id);
CREATE INDEX idx_vehicle_condition_assessments_deleted_at ON vehicle_condition_assessments(deleted_at);

CREATE INDEX idx_vehicle_photos_inventory_id ON vehicle_photos(inventory_id);
CREATE INDEX idx_vehicle_photos_photo_type ON vehicle_photos(photo_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_photos_is_primary ON vehicle_photos(is_primary) WHERE deleted_at IS NULL;
CREATE INDEX idx_vehicle_photos_deleted_at ON vehicle_photos(deleted_at);

CREATE INDEX idx_sales_commissions_sale_id ON sales_commissions(sale_id);
CREATE INDEX idx_sales_commissions_sales_person_id ON sales_commissions(sales_person_id);
CREATE INDEX idx_sales_commissions_payment_status ON sales_commissions(payment_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_sales_commissions_deleted_at ON sales_commissions(deleted_at);
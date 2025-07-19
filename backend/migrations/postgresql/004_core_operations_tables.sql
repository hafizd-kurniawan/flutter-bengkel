-- Core Operations Tables (PostgreSQL) - Enhanced for POS & Service Management

-- Service jobs with queue management and auto-assignment
CREATE TABLE service_jobs (
    id BIGSERIAL PRIMARY KEY,
    job_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    vehicle_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    technician_id BIGINT,
    queue_number INTEGER NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')),
    problem_description TEXT NOT NULL,
    estimated_completion TIMESTAMP,
    actual_completion TIMESTAMP,
    total_amount DECIMAL(15,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    final_amount DECIMAL(15,2) DEFAULT 0,
    warranty_period_days INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (vehicle_id) REFERENCES customer_vehicles(id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (technician_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Service job details (services and products used)
CREATE TABLE service_details (
    id BIGSERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    product_id BIGINT,
    service_id BIGINT,
    quantity DECIMAL(10,3) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL,
    CHECK (product_id IS NOT NULL OR service_id IS NOT NULL)
);

-- Service job history for tracking progress and audit
CREATE TABLE service_job_histories (
    id BIGSERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    previous_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Transactions table for all business operations
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    transaction_number VARCHAR(50) NOT NULL UNIQUE,
    transaction_type VARCHAR(30) NOT NULL CHECK (transaction_type IN ('service', 'sparepart_sale', 'vehicle_purchase', 'vehicle_sale')),
    customer_id BIGINT,
    outlet_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    service_job_id BIGINT,
    subtotal_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'partial', 'paid', 'cancelled')),
    notes TEXT,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE SET NULL
);

-- Transaction details
CREATE TABLE transaction_details (
    id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    product_id BIGINT,
    service_id BIGINT,
    description TEXT,
    quantity DECIMAL(10,3) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL
);

-- Purchase orders for inventory management
CREATE TABLE purchase_orders (
    id BIGSERIAL PRIMARY KEY,
    po_number VARCHAR(50) NOT NULL UNIQUE,
    supplier_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'confirmed', 'received', 'cancelled')),
    order_date DATE NOT NULL,
    expected_delivery_date DATE,
    actual_delivery_date DATE,
    subtotal_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Purchase order details
CREATE TABLE purchase_order_details (
    id BIGSERIAL PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity DECIMAL(10,3) NOT NULL,
    unit_cost DECIMAL(15,2) NOT NULL,
    total_cost DECIMAL(15,2) NOT NULL,
    received_quantity DECIMAL(10,3) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Enhanced vehicle purchases for vehicle trading business
CREATE TABLE vehicle_purchases (
    id BIGSERIAL PRIMARY KEY,
    purchase_number VARCHAR(50) NOT NULL UNIQUE,
    vehicle_number VARCHAR(20) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    color VARCHAR(50),
    mileage BIGINT,
    condition_rating VARCHAR(20) DEFAULT 'good' CHECK (condition_rating IN ('excellent', 'good', 'fair', 'poor')),
    purchase_price DECIMAL(15,2) NOT NULL,
    estimated_selling_price DECIMAL(15,2),
    actual_selling_price DECIMAL(15,2),
    seller_name VARCHAR(255),
    seller_phone VARCHAR(20),
    seller_address TEXT,
    status VARCHAR(20) DEFAULT 'purchased' CHECK (status IN ('purchased', 'available', 'sold', 'under_repair')),
    purchase_date DATE NOT NULL,
    sale_date DATE,
    outlet_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    notes TEXT,
    -- New fields for enhanced vehicle trading
    service_required BOOLEAN DEFAULT TRUE,
    service_completion_date TIMESTAMP NULL,
    selling_price DECIMAL(15,2) NULL,
    sale_status VARCHAR(20) DEFAULT 'Available' CHECK (sale_status IN ('Available', 'Reserved', 'Sold')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Vehicle sales tracking for profit calculation
CREATE TABLE vehicle_sales (
    id BIGSERIAL PRIMARY KEY,
    sale_id SERIAL,
    vehicle_purchase_id BIGINT REFERENCES vehicle_purchases(id),
    customer_id BIGINT REFERENCES customers(id),
    sales_user_id BIGINT REFERENCES users(id),
    sale_date DATE NOT NULL,
    selling_price DECIMAL(15,2) NOT NULL,
    profit_amount DECIMAL(15,2) NOT NULL,
    commission_amount DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by BIGINT REFERENCES users(id)
);

-- Create indexes for better performance
CREATE INDEX idx_service_jobs_job_number ON service_jobs(job_number);
CREATE INDEX idx_service_jobs_customer_id ON service_jobs(customer_id);
CREATE INDEX idx_service_jobs_vehicle_id ON service_jobs(vehicle_id);
CREATE INDEX idx_service_jobs_technician_id ON service_jobs(technician_id);
CREATE INDEX idx_service_jobs_status ON service_jobs(status);
CREATE INDEX idx_service_jobs_queue ON service_jobs(outlet_id, queue_number);
CREATE INDEX idx_service_jobs_created_by ON service_jobs(created_by);
CREATE INDEX idx_service_details_service_job_id ON service_details(service_job_id);
CREATE INDEX idx_transactions_transaction_number ON transactions(transaction_number);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_customer_id ON transactions(customer_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_purchase_orders_po_number ON purchase_orders(po_number);
CREATE INDEX idx_purchase_orders_supplier_id ON purchase_orders(supplier_id);
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX idx_vehicle_purchases_purchase_number ON vehicle_purchases(purchase_number);
CREATE INDEX idx_vehicle_purchases_vehicle_number ON vehicle_purchases(vehicle_number);
CREATE INDEX idx_vehicle_purchases_status ON vehicle_purchases(status);
CREATE INDEX idx_vehicle_purchases_sale_status ON vehicle_purchases(sale_status);
CREATE INDEX idx_vehicle_sales_vehicle_purchase_id ON vehicle_sales(vehicle_purchase_id);
CREATE INDEX idx_vehicle_sales_customer_id ON vehicle_sales(customer_id);
CREATE INDEX idx_vehicle_sales_sales_user_id ON vehicle_sales(sales_user_id);

-- Create triggers for updated_at columns
CREATE TRIGGER update_service_jobs_updated_at BEFORE UPDATE ON service_jobs
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_purchase_orders_updated_at BEFORE UPDATE ON purchase_orders
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_vehicle_purchases_updated_at BEFORE UPDATE ON vehicle_purchases
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_vehicle_sales_updated_at BEFORE UPDATE ON vehicle_sales
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
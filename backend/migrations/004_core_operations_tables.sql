-- Core Operations Tables (PostgreSQL with Soft Delete)

-- Service jobs with queue management
CREATE TABLE service_jobs (
    job_id BIGSERIAL PRIMARY KEY,
    job_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    vehicle_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    technician_id BIGINT,
    queue_number INTEGER NOT NULL,
    priority VARCHAR(20) CHECK (priority IN ('low', 'normal', 'high', 'urgent')) DEFAULT 'normal',
    status VARCHAR(20) CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')) DEFAULT 'pending',
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
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    FOREIGN KEY (vehicle_id) REFERENCES customer_vehicles(vehicle_id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id),
    FOREIGN KEY (technician_id) REFERENCES users(user_id) ON DELETE SET NULL
);

-- Service job details (services and products used)
CREATE TABLE service_details (
    detail_id BIGSERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    product_id BIGINT,
    service_id BIGINT,
    quantity DECIMAL(10,3) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(job_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE SET NULL,
    CHECK (product_id IS NOT NULL OR service_id IS NOT NULL)
);

-- Service job history for tracking progress
CREATE TABLE service_job_histories (
    history_id BIGSERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    previous_status VARCHAR(20) CHECK (previous_status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')),
    new_status VARCHAR(20) CHECK (new_status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(job_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Transactions table for all business operations
CREATE TABLE transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    transaction_number VARCHAR(50) NOT NULL UNIQUE,
    transaction_type VARCHAR(20) CHECK (transaction_type IN ('service', 'sparepart_sale', 'vehicle_purchase', 'vehicle_sale')) NOT NULL,
    customer_id BIGINT,
    outlet_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    service_job_id BIGINT,
    subtotal_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    payment_status VARCHAR(20) CHECK (payment_status IN ('pending', 'partial', 'paid', 'cancelled')) DEFAULT 'pending',
    notes TEXT,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE SET NULL,
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(job_id) ON DELETE SET NULL
);

-- Transaction details
CREATE TABLE transaction_details (
    detail_id BIGSERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    product_id BIGINT,
    service_id BIGINT,
    description TEXT,
    quantity DECIMAL(10,3) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(service_id) ON DELETE SET NULL
);

-- Purchase orders for inventory management
CREATE TABLE purchase_orders (
    po_id BIGSERIAL PRIMARY KEY,
    po_number VARCHAR(50) NOT NULL UNIQUE,
    supplier_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    status VARCHAR(20) CHECK (status IN ('draft', 'sent', 'confirmed', 'received', 'cancelled')) DEFAULT 'draft',
    order_date DATE NOT NULL,
    expected_delivery_date DATE,
    actual_delivery_date DATE,
    subtotal_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    tax_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(supplier_id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Purchase order details
CREATE TABLE purchase_order_details (
    detail_id BIGSERIAL PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity DECIMAL(10,3) NOT NULL,
    unit_cost DECIMAL(15,2) NOT NULL,
    total_cost DECIMAL(15,2) NOT NULL,
    received_quantity DECIMAL(10,3) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(po_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id)
);

-- Create indexes for better performance
CREATE INDEX idx_service_jobs_job_number ON service_jobs(job_number);
CREATE INDEX idx_service_jobs_customer_id ON service_jobs(customer_id);
CREATE INDEX idx_service_jobs_vehicle_id ON service_jobs(vehicle_id);
CREATE INDEX idx_service_jobs_technician_id ON service_jobs(technician_id);
CREATE INDEX idx_service_jobs_status ON service_jobs(status);
CREATE INDEX idx_service_jobs_queue ON service_jobs(outlet_id, queue_number);
CREATE INDEX idx_service_details_service_job_id ON service_details(service_job_id);
CREATE INDEX idx_transactions_transaction_number ON transactions(transaction_number);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_customer_id ON transactions(customer_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_purchase_orders_po_number ON purchase_orders(po_number);
CREATE INDEX idx_purchase_orders_supplier_id ON purchase_orders(supplier_id);
CREATE INDEX idx_purchase_orders_status ON purchase_orders(status);
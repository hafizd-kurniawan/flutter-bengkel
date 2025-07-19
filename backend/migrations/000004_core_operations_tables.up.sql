-- Core Operations Tables

-- Service jobs with queue management
CREATE TABLE service_jobs (
    id SERIAL PRIMARY KEY,
    job_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    vehicle_id BIGINT NOT NULL,
    outlet_id BIGINT NOT NULL,
    technician_id BIGINT,
    queue_number INTEGER NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')),
    problem_description TEXT NOT NULL,
    estimated_completion TIMESTAMP WITH TIME ZONE,
    actual_completion TIMESTAMP WITH TIME ZONE,
    total_amount DECIMAL(15,2) DEFAULT 0,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    tax_amount DECIMAL(15,2) DEFAULT 0,
    final_amount DECIMAL(15,2) DEFAULT 0,
    warranty_period_days INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (vehicle_id) REFERENCES customer_vehicles(id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (technician_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Service job details (services and products used)
CREATE TABLE service_details (
    id SERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    product_id BIGINT,
    service_id BIGINT,
    quantity DECIMAL(10,3) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL,
    CHECK (product_id IS NOT NULL OR service_id IS NOT NULL)
);

-- Service job history for tracking progress
CREATE TABLE service_job_histories (
    id SERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    previous_status VARCHAR(20) CHECK (previous_status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')),
    new_status VARCHAR(20) NOT NULL CHECK (new_status IN ('pending', 'in_progress', 'completed', 'cancelled', 'on_hold')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Transactions table for all business operations
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
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
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE SET NULL
);

-- Transaction details
CREATE TABLE transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    product_id BIGINT,
    service_id BIGINT,
    description TEXT,
    quantity DECIMAL(10,3) NOT NULL DEFAULT 1,
    unit_price DECIMAL(15,2) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE SET NULL
);

-- Purchase orders for inventory management
CREATE TABLE purchase_orders (
    id SERIAL PRIMARY KEY,
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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Purchase order details
CREATE TABLE purchase_order_details (
    id SERIAL PRIMARY KEY,
    purchase_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity DECIMAL(10,3) NOT NULL,
    unit_cost DECIMAL(15,2) NOT NULL,
    total_cost DECIMAL(15,2) NOT NULL,
    received_quantity DECIMAL(10,3) DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- Vehicle purchases (for vehicle trading business)
CREATE TABLE vehicle_purchases (
    id SERIAL PRIMARY KEY,
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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
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
CREATE INDEX idx_vehicle_purchases_purchase_number ON vehicle_purchases(purchase_number);
CREATE INDEX idx_vehicle_purchases_vehicle_number ON vehicle_purchases(vehicle_number);
CREATE INDEX idx_vehicle_purchases_status ON vehicle_purchases(status);
-- Multi-Payment System & Extended Features

-- Transaction Payments (Multi-Payment Support)
CREATE TABLE transaction_payments (
    id SERIAL PRIMARY KEY,
    transaction_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_order INTEGER NOT NULL,
    reference_number VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Service Job Complete Schema Update
CREATE TABLE service_job_queue (
    id SERIAL PRIMARY KEY,
    outlet_id BIGINT NOT NULL,
    current_queue_number INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    UNIQUE(outlet_id)
);

-- Vehicle Trading Enhancement
CREATE TABLE vehicle_sales (
    id SERIAL PRIMARY KEY,
    vehicle_purchase_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    sales_user_id BIGINT NOT NULL,
    sale_date DATE NOT NULL,
    selling_price DECIMAL(15,2) NOT NULL,
    profit_amount DECIMAL(15,2) NOT NULL,
    commission_amount DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_purchase_id) REFERENCES vehicle_purchases(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (sales_user_id) REFERENCES users(id)
);

-- Workshop Configuration
CREATE TABLE workshop_settings (
    id SERIAL PRIMARY KEY,
    outlet_id BIGINT NOT NULL,
    setting_key VARCHAR(100) NOT NULL,
    setting_value TEXT,
    data_type VARCHAR(20) DEFAULT 'string' CHECK (data_type IN ('string', 'integer', 'decimal', 'boolean', 'json')),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    UNIQUE(outlet_id, setting_key)
);

-- Service Technician Assignment
CREATE TABLE technician_assignments (
    id SERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    technician_id BIGINT NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    assigned_by BIGINT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id),
    FOREIGN KEY (technician_id) REFERENCES users(id),
    FOREIGN KEY (assigned_by) REFERENCES users(id)
);

-- Customer Communication Log
CREATE TABLE customer_communications (
    id SERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    service_job_id BIGINT,
    communication_type VARCHAR(20) NOT NULL CHECK (communication_type IN ('sms', 'whatsapp', 'call', 'email', 'in_person')),
    subject VARCHAR(255),
    message TEXT,
    sent_by BIGINT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delivery_status VARCHAR(20) DEFAULT 'sent' CHECK (delivery_status IN ('pending', 'sent', 'delivered', 'failed')),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id) ON DELETE SET NULL,
    FOREIGN KEY (sent_by) REFERENCES users(id)
);

-- Product Movement History
CREATE TABLE product_movements (
    id SERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('purchase', 'sale', 'adjustment', 'transfer', 'damage', 'return')),
    quantity_change INTEGER NOT NULL,
    unit_cost DECIMAL(15,2),
    reference_id BIGINT,
    reference_type VARCHAR(50),
    notes TEXT,
    moved_by BIGINT NOT NULL,
    moved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (moved_by) REFERENCES users(id)
);

-- Service Performance Metrics
CREATE TABLE service_metrics (
    id SERIAL PRIMARY KEY,
    service_job_id BIGINT NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(15,4),
    metric_unit VARCHAR(20),
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    recorded_by BIGINT NOT NULL,
    FOREIGN KEY (service_job_id) REFERENCES service_jobs(id),
    FOREIGN KEY (recorded_by) REFERENCES users(id)
);

-- Create indexes for multi-payment system
CREATE INDEX idx_transaction_payments_transaction_id ON transaction_payments(transaction_id);
CREATE INDEX idx_transaction_payments_payment_order ON transaction_payments(transaction_id, payment_order);
CREATE INDEX idx_vehicle_sales_vehicle_purchase_id ON vehicle_sales(vehicle_purchase_id);
CREATE INDEX idx_vehicle_sales_customer_id ON vehicle_sales(customer_id);
CREATE INDEX idx_vehicle_sales_sale_date ON vehicle_sales(sale_date);
CREATE INDEX idx_workshop_settings_outlet_key ON workshop_settings(outlet_id, setting_key);
CREATE INDEX idx_technician_assignments_service_job ON technician_assignments(service_job_id);
CREATE INDEX idx_technician_assignments_technician ON technician_assignments(technician_id);
CREATE INDEX idx_customer_communications_customer ON customer_communications(customer_id);
CREATE INDEX idx_customer_communications_service_job ON customer_communications(service_job_id);
CREATE INDEX idx_product_movements_product ON product_movements(product_id);
CREATE INDEX idx_product_movements_type_date ON product_movements(movement_type, moved_at);
CREATE INDEX idx_service_metrics_job_metric ON service_metrics(service_job_id, metric_name);
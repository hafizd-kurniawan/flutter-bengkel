-- Financial & Reporting Tables (PostgreSQL) - Enhanced with Multi-Payment Support

-- Payment methods
CREATE TABLE payment_methods (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(30) NOT NULL CHECK (type IN ('cash', 'bank_transfer', 'credit_card', 'debit_card', 'e_wallet', 'check')),
    account_number VARCHAR(100),
    bank_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced payments table
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    transaction_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Multi-payment support table (NEW)
CREATE TABLE transaction_payments (
    id BIGSERIAL PRIMARY KEY,
    payment_id SERIAL,
    transaction_id BIGINT REFERENCES transactions(id),
    payment_method_id BIGINT REFERENCES payment_methods(id),
    amount DECIMAL(15,2) NOT NULL,
    payment_order INTEGER NOT NULL, -- 1st payment, 2nd payment, etc
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Accounts Payable (money we owe to suppliers)
CREATE TABLE accounts_payables (
    id BIGSERIAL PRIMARY KEY,
    ap_number VARCHAR(50) NOT NULL UNIQUE,
    supplier_id BIGINT NOT NULL,
    purchase_order_id BIGINT,
    amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'outstanding' CHECK (status IN ('outstanding', 'partial', 'paid', 'overdue')),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL
);

-- Payable payments
CREATE TABLE payable_payments (
    id BIGSERIAL PRIMARY KEY,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    accounts_payable_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accounts_payable_id) REFERENCES accounts_payables(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Accounts Receivable (money customers owe us) - Enhanced with Kasir approval
CREATE TABLE accounts_receivables (
    id BIGSERIAL PRIMARY KEY,
    ar_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    transaction_id BIGINT,
    amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'outstanding' CHECK (status IN ('outstanding', 'partial', 'paid', 'overdue')),
    notes TEXT,
    -- New fields for kasir approval
    approved_by BIGINT, -- Kasir who approved payment
    approval_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    FOREIGN KEY (approved_by) REFERENCES users(id)
);

-- Receivable payments with kasir approval tracking
CREATE TABLE receivable_payments (
    id BIGSERIAL PRIMARY KEY,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    accounts_receivable_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    -- Kasir approval fields
    approved_by BIGINT, -- Kasir who approved this payment
    approval_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accounts_receivable_id) REFERENCES accounts_receivables(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
);

-- Cash flows
CREATE TABLE cash_flows (
    id BIGSERIAL PRIMARY KEY,
    outlet_id BIGINT NOT NULL,
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('income', 'expense')),
    category VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_method_id BIGINT NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    description TEXT,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Reports
CREATE TABLE reports (
    id BIGSERIAL PRIMARY KEY,
    report_type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    parameters JSONB,
    generated_by BIGINT NOT NULL,
    file_path VARCHAR(500),
    status VARCHAR(20) DEFAULT 'generating' CHECK (status IN ('generating', 'completed', 'failed')),
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (generated_by) REFERENCES users(id)
);

-- Promotions
CREATE TABLE promotions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    promotion_type VARCHAR(30) NOT NULL CHECK (promotion_type IN ('discount_percentage', 'discount_amount', 'buy_x_get_y', 'free_service')),
    discount_value DECIMAL(15,2) DEFAULT 0,
    minimum_purchase DECIMAL(15,2) DEFAULT 0,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    usage_limit INTEGER,
    usage_count INTEGER DEFAULT 0,
    applicable_to VARCHAR(20) DEFAULT 'all' CHECK (applicable_to IN ('all', 'service', 'sparepart', 'vehicle')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Commission tracking table (NEW)
CREATE TABLE commission_tracking (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    transaction_id BIGINT,
    vehicle_sale_id BIGINT,
    commission_type VARCHAR(30) NOT NULL CHECK (commission_type IN ('sale', 'service', 'vehicle_trading')),
    base_amount DECIMAL(15,2) NOT NULL,
    commission_rate DECIMAL(5,2) NOT NULL, -- Percentage
    commission_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'paid')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_by BIGINT,
    approved_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (vehicle_sale_id) REFERENCES vehicle_sales(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
);

-- Create indexes for financial tables
CREATE INDEX idx_payment_methods_type ON payment_methods(type);
CREATE INDEX idx_payment_methods_is_active ON payment_methods(is_active);
CREATE INDEX idx_payments_payment_number ON payments(payment_number);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);
CREATE INDEX idx_transaction_payments_transaction_id ON transaction_payments(transaction_id);
CREATE INDEX idx_transaction_payments_payment_method_id ON transaction_payments(payment_method_id);
CREATE INDEX idx_accounts_payables_supplier_id ON accounts_payables(supplier_id);
CREATE INDEX idx_accounts_payables_status ON accounts_payables(status);
CREATE INDEX idx_accounts_payables_due_date ON accounts_payables(due_date);
CREATE INDEX idx_accounts_receivables_customer_id ON accounts_receivables(customer_id);
CREATE INDEX idx_accounts_receivables_status ON accounts_receivables(status);
CREATE INDEX idx_accounts_receivables_due_date ON accounts_receivables(due_date);
CREATE INDEX idx_accounts_receivables_approved_by ON accounts_receivables(approved_by);
CREATE INDEX idx_receivable_payments_approved_by ON receivable_payments(approved_by);
CREATE INDEX idx_cash_flows_outlet_id ON cash_flows(outlet_id);
CREATE INDEX idx_cash_flows_transaction_date ON cash_flows(transaction_date);
CREATE INDEX idx_cash_flows_type ON cash_flows(transaction_type);
CREATE INDEX idx_reports_generated_by ON reports(generated_by);
CREATE INDEX idx_reports_type ON reports(report_type);
CREATE INDEX idx_promotions_dates ON promotions(start_date, end_date);
CREATE INDEX idx_promotions_is_active ON promotions(is_active);
CREATE INDEX idx_commission_tracking_user_id ON commission_tracking(user_id);
CREATE INDEX idx_commission_tracking_status ON commission_tracking(status);
CREATE INDEX idx_commission_tracking_type ON commission_tracking(commission_type);

-- Create triggers for updated_at columns
CREATE TRIGGER update_payment_methods_updated_at BEFORE UPDATE ON payment_methods
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_accounts_payables_updated_at BEFORE UPDATE ON accounts_payables
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_accounts_receivables_updated_at BEFORE UPDATE ON accounts_receivables
FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_promotions_updated_at BEFORE UPDATE ON promotions
FOR EACH ROW EXECUTE FUNCTION update_modified_column();
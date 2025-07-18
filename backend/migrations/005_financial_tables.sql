-- Financial & Reporting Tables

-- Payment methods
CREATE TABLE payment_methods (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    type ENUM('cash', 'bank_transfer', 'credit_card', 'debit_card', 'e_wallet', 'check') NOT NULL,
    account_number VARCHAR(100),
    bank_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Payments
CREATE TABLE payments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    transaction_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Accounts Payable (money we owe to suppliers)
CREATE TABLE accounts_payables (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ap_number VARCHAR(50) NOT NULL UNIQUE,
    supplier_id BIGINT NOT NULL,
    purchase_order_id BIGINT,
    amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status ENUM('outstanding', 'partial', 'paid', 'overdue') DEFAULT 'outstanding',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL
);

-- Payable payments
CREATE TABLE payable_payments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    accounts_payable_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accounts_payable_id) REFERENCES accounts_payables(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Accounts Receivable (money customers owe us)
CREATE TABLE accounts_receivables (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ar_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    transaction_id BIGINT,
    amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status ENUM('outstanding', 'partial', 'paid', 'overdue') DEFAULT 'outstanding',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL
);

-- Receivable payments
CREATE TABLE receivable_payments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    accounts_receivable_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (accounts_receivable_id) REFERENCES accounts_receivables(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Cash flows
CREATE TABLE cash_flows (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    outlet_id BIGINT NOT NULL,
    transaction_type ENUM('income', 'expense') NOT NULL,
    category VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_method_id BIGINT NOT NULL,
    reference_id BIGINT,
    reference_type VARCHAR(50),
    description TEXT,
    transaction_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (outlet_id) REFERENCES outlets(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

-- Reports
CREATE TABLE reports (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    report_type VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    parameters JSON,
    generated_by BIGINT NOT NULL,
    file_path VARCHAR(500),
    status ENUM('generating', 'completed', 'failed') DEFAULT 'generating',
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (generated_by) REFERENCES users(id)
);

-- Promotions
CREATE TABLE promotions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    promotion_type ENUM('discount_percentage', 'discount_amount', 'buy_x_get_y', 'free_service') NOT NULL,
    discount_value DECIMAL(15,2) DEFAULT 0,
    minimum_purchase DECIMAL(15,2) DEFAULT 0,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    usage_limit INT,
    usage_count INT DEFAULT 0,
    applicable_to ENUM('all', 'service', 'sparepart', 'vehicle') DEFAULT 'all',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create indexes for financial tables
CREATE INDEX idx_payments_payment_number ON payments(payment_number);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);
CREATE INDEX idx_accounts_payables_supplier_id ON accounts_payables(supplier_id);
CREATE INDEX idx_accounts_payables_status ON accounts_payables(status);
CREATE INDEX idx_accounts_payables_due_date ON accounts_payables(due_date);
CREATE INDEX idx_accounts_receivables_customer_id ON accounts_receivables(customer_id);
CREATE INDEX idx_accounts_receivables_status ON accounts_receivables(status);
CREATE INDEX idx_accounts_receivables_due_date ON accounts_receivables(due_date);
CREATE INDEX idx_cash_flows_outlet_id ON cash_flows(outlet_id);
CREATE INDEX idx_cash_flows_transaction_date ON cash_flows(transaction_date);
CREATE INDEX idx_cash_flows_type ON cash_flows(transaction_type);
CREATE INDEX idx_reports_generated_by ON reports(generated_by);
CREATE INDEX idx_reports_type ON reports(report_type);
CREATE INDEX idx_promotions_dates ON promotions(start_date, end_date);
CREATE INDEX idx_promotions_is_active ON promotions(is_active);
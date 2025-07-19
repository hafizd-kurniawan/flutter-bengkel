-- Financial & Reporting Tables (PostgreSQL with Soft Delete)

-- Payment methods
CREATE TABLE payment_methods (
    method_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) CHECK (type IN ('cash', 'bank_transfer', 'credit_card', 'debit_card', 'e_wallet', 'check')) NOT NULL,
    account_number VARCHAR(100),
    bank_name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER
);

-- Payments
CREATE TABLE payments (
    payment_id BIGSERIAL PRIMARY KEY,
    payment_number VARCHAR(50) NOT NULL UNIQUE,
    transaction_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(method_id)
);

-- Accounts Payable (money we owe to suppliers)
CREATE TABLE accounts_payables (
    payable_id BIGSERIAL PRIMARY KEY,
    ap_number VARCHAR(50) NOT NULL UNIQUE,
    supplier_id BIGINT NOT NULL,
    purchase_order_id BIGINT,
    amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) CHECK (status IN ('outstanding', 'partial', 'paid', 'overdue')) DEFAULT 'outstanding',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(supplier_id),
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(po_id) ON DELETE SET NULL
);

-- Payable payments
CREATE TABLE payable_payments (
    payment_id BIGSERIAL PRIMARY KEY,
    accounts_payable_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (accounts_payable_id) REFERENCES accounts_payables(payable_id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(method_id)
);

-- Accounts Receivable (money customers owe us)
CREATE TABLE accounts_receivables (
    receivable_id BIGSERIAL PRIMARY KEY,
    ar_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id BIGINT NOT NULL,
    transaction_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    paid_amount DECIMAL(15,2) DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) CHECK (status IN ('outstanding', 'partial', 'paid', 'overdue')) DEFAULT 'outstanding',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

-- Receivable payments
CREATE TABLE receivable_payments (
    payment_id BIGSERIAL PRIMARY KEY,
    accounts_receivable_id BIGINT NOT NULL,
    payment_method_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (accounts_receivable_id) REFERENCES accounts_receivables(receivable_id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(method_id)
);

-- Cash flows for tracking all money movements
CREATE TABLE cash_flows (
    flow_id BIGSERIAL PRIMARY KEY,
    outlet_id BIGINT NOT NULL,
    flow_type VARCHAR(20) CHECK (flow_type IN ('inflow', 'outflow')) NOT NULL,
    category VARCHAR(50) NOT NULL, -- 'sales', 'purchase', 'expense', 'investment', etc.
    amount DECIMAL(15,2) NOT NULL,
    description TEXT NOT NULL,
    reference_type VARCHAR(50), -- 'transaction', 'payment', 'expense', etc.
    reference_id BIGINT,
    transaction_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id)
);

-- Commissions for vehicle trading
CREATE TABLE commissions (
    commission_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    transaction_id BIGINT NOT NULL,
    commission_type VARCHAR(20) CHECK (commission_type IN ('sale', 'purchase', 'service')) NOT NULL,
    commission_rate DECIMAL(5,2) NOT NULL, -- Percentage
    commission_amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'paid', 'cancelled')) DEFAULT 'pending',
    paid_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

-- Daily cash summary for outlets
CREATE TABLE daily_cash_summaries (
    summary_id BIGSERIAL PRIMARY KEY,
    outlet_id BIGINT NOT NULL,
    summary_date DATE NOT NULL,
    opening_balance DECIMAL(15,2) DEFAULT 0,
    total_inflow DECIMAL(15,2) DEFAULT 0,
    total_outflow DECIMAL(15,2) DEFAULT 0,
    closing_balance DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INTEGER,
    FOREIGN KEY (outlet_id) REFERENCES outlets(outlet_id),
    UNIQUE(outlet_id, summary_date)
);

-- Create indexes for financial performance
CREATE INDEX idx_payments_payment_number ON payments(payment_number);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_date ON payments(payment_date);
CREATE INDEX idx_accounts_payables_supplier_id ON accounts_payables(supplier_id);
CREATE INDEX idx_accounts_payables_status ON accounts_payables(status);
CREATE INDEX idx_accounts_payables_due_date ON accounts_payables(due_date);
CREATE INDEX idx_accounts_receivables_customer_id ON accounts_receivables(customer_id);
CREATE INDEX idx_accounts_receivables_status ON accounts_receivables(status);
CREATE INDEX idx_accounts_receivables_due_date ON accounts_receivables(due_date);
CREATE INDEX idx_cash_flows_outlet_id ON cash_flows(outlet_id);
CREATE INDEX idx_cash_flows_type ON cash_flows(flow_type);
CREATE INDEX idx_cash_flows_date ON cash_flows(transaction_date);
CREATE INDEX idx_commissions_user_id ON commissions(user_id);
CREATE INDEX idx_commissions_status ON commissions(status);
CREATE INDEX idx_daily_summaries_outlet_date ON daily_cash_summaries(outlet_id, summary_date);
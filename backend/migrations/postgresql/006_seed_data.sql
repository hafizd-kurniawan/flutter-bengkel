-- Seed Data for Workshop Management System (PostgreSQL)

-- Insert default roles
INSERT INTO roles (name, description) VALUES
('Admin', 'System administrator with full access'),
('Manager', 'Workshop manager with management access'),
('Kasir', 'Cashier with POS and service management access'),
('Sales', 'Sales team with vehicle trading access'),
('Technician', 'Mechanic/technician with service access'),
('Customer Service', 'Customer service representative');

-- Insert permissions
INSERT INTO permissions (name, description, resource, action) VALUES
-- User management
('create_users', 'Create new users', 'users', 'create'),
('read_users', 'View user information', 'users', 'read'),
('update_users', 'Update user information', 'users', 'update'),
('delete_users', 'Delete users', 'users', 'delete'),

-- Customer management
('create_customers', 'Create new customers', 'customers', 'create'),
('read_customers', 'View customer information', 'customers', 'read'),
('update_customers', 'Update customer information', 'customers', 'update'),
('delete_customers', 'Delete customers', 'customers', 'delete'),

-- Vehicle management
('create_vehicles', 'Create new vehicles', 'vehicles', 'create'),
('read_vehicles', 'View vehicle information', 'vehicles', 'read'),
('update_vehicles', 'Update vehicle information', 'vehicles', 'update'),
('delete_vehicles', 'Delete vehicles', 'vehicles', 'delete'),

-- Service job management
('create_service_jobs', 'Create new service jobs', 'service_jobs', 'create'),
('read_service_jobs', 'View service jobs', 'service_jobs', 'read'),
('update_service_jobs', 'Update service jobs', 'service_jobs', 'update'),
('delete_service_jobs', 'Delete service jobs', 'service_jobs', 'delete'),
('assign_technicians', 'Assign technicians to jobs', 'service_jobs', 'assign'),

-- Product and inventory
('create_products', 'Create new products', 'products', 'create'),
('read_products', 'View product information', 'products', 'read'),
('update_products', 'Update product information', 'products', 'update'),
('delete_products', 'Delete products', 'products', 'delete'),
('manage_stock', 'Manage product stock levels', 'products', 'stock'),

-- Transactions and POS
('create_transactions', 'Create new transactions', 'transactions', 'create'),
('read_transactions', 'View transactions', 'transactions', 'read'),
('update_transactions', 'Update transactions', 'transactions', 'update'),
('delete_transactions', 'Delete transactions', 'transactions', 'delete'),
('process_payments', 'Process payments', 'payments', 'create'),
('approve_receivables', 'Approve receivable payments', 'receivables', 'approve'),

-- Vehicle trading
('create_vehicle_purchases', 'Purchase vehicles', 'vehicle_purchases', 'create'),
('read_vehicle_purchases', 'View vehicle inventory', 'vehicle_purchases', 'read'),
('update_vehicle_prices', 'Update vehicle selling prices', 'vehicle_purchases', 'update_price'),
('manage_vehicle_sales', 'Manage vehicle sales', 'vehicle_sales', 'manage'),

-- Reports
('view_reports', 'View reports', 'reports', 'read'),
('generate_reports', 'Generate reports', 'reports', 'create'),
('kasir_personal_reports', 'View personal kasir reports', 'reports', 'kasir_personal');

-- Assign permissions to roles

-- Admin gets all permissions
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- Manager gets most permissions except user deletion
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE name != 'delete_users';

-- Kasir permissions (POS-centric)
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE name IN (
    'create_customers', 'read_customers', 'update_customers',
    'create_vehicles', 'read_vehicles', 'update_vehicles',
    'create_service_jobs', 'read_service_jobs', 'update_service_jobs',
    'read_products', 'create_transactions', 'read_transactions', 'update_transactions',
    'process_payments', 'approve_receivables', 'kasir_personal_reports'
);

-- Sales permissions (Vehicle trading focused)
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 4, id FROM permissions WHERE name IN (
    'read_customers', 'read_vehicles', 'read_vehicle_purchases', 
    'update_vehicle_prices', 'manage_vehicle_sales', 'view_reports'
);

-- Technician permissions
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 5, id FROM permissions WHERE name IN (
    'read_customers', 'read_vehicles', 'read_service_jobs', 'update_service_jobs', 'read_products'
);

-- Customer Service permissions
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 6, id FROM permissions WHERE name IN (
    'create_customers', 'read_customers', 'update_customers',
    'create_vehicles', 'read_vehicles', 'update_vehicles',
    'read_service_jobs', 'read_transactions'
);

-- Insert default outlet
INSERT INTO outlets (name, address, phone, email) VALUES
('Main Workshop', 'Jl. Raya Utama No. 123, Jakarta', '021-12345678', 'main@bengkel.com');

-- Insert default admin user (password: admin123)
INSERT INTO users (username, email, password_hash, full_name, phone, role_id, outlet_id, created_by) VALUES
('admin', 'admin@bengkel.com', '$2a$10$YnNE8QfzN9p7fR1lW3RjWuZTz.lGtF8rH0G5bC1iZ2lF3rN9.vz6G', 'System Administrator', '081234567890', 1, 1, 1);

-- Insert default payment methods
INSERT INTO payment_methods (name, type, is_active) VALUES
('Cash', 'cash', TRUE),
('Bank Transfer BCA', 'bank_transfer', TRUE),
('Bank Transfer Mandiri', 'bank_transfer', TRUE),
('Credit Card Visa', 'credit_card', TRUE),
('Credit Card Mastercard', 'credit_card', TRUE),
('Debit Card', 'debit_card', TRUE),
('GoPay', 'e_wallet', TRUE),
('OVO', 'e_wallet', TRUE),
('DANA', 'e_wallet', TRUE);

-- Insert default unit types
INSERT INTO unit_types (name, abbreviation, description) VALUES
('Pieces', 'pcs', 'Individual pieces or units'),
('Liter', 'ltr', 'Volume measurement in liters'),
('Kilogram', 'kg', 'Weight measurement in kilograms'),
('Meter', 'm', 'Length measurement in meters'),
('Set', 'set', 'Complete set of items'),
('Bottle', 'btl', 'Items sold in bottles'),
('Box', 'box', 'Items sold in boxes'),
('Service Hour', 'hr', 'Service time in hours');

-- Insert default service categories
INSERT INTO service_categories (name, description) VALUES
('Oil Change', 'Regular engine oil and filter changes'),
('Brake Service', 'Brake inspection, repair and replacement'),
('Engine Service', 'Engine maintenance and repair'),
('Transmission Service', 'Transmission maintenance and repair'),
('Electrical Service', 'Electrical system diagnosis and repair'),
('Air Conditioning', 'AC maintenance and repair'),
('Body Work', 'Body repair and painting services'),
('Tire Service', 'Tire installation, balancing, and alignment');

-- Insert sample services
INSERT INTO services (service_code, name, description, category_id, standard_price, estimated_duration) VALUES
('SVC001', 'Engine Oil Change', 'Regular engine oil and filter replacement', 1, 150000, 30),
('SVC002', 'Brake Pad Replacement', 'Front or rear brake pad replacement', 2, 300000, 60),
('SVC003', 'Engine Tune-Up', 'Complete engine tuning and adjustment', 3, 500000, 120),
('SVC004', 'Transmission Service', 'Transmission fluid change and inspection', 4, 400000, 90),
('SVC005', 'AC Service', 'Air conditioning cleaning and refrigerant refill', 6, 200000, 45),
('SVC006', 'Wheel Alignment', 'Front wheel alignment service', 8, 100000, 30);

-- Insert default product categories
INSERT INTO categories (name, description) VALUES
('Engine Oil', 'Various types of engine oils'),
('Filters', 'Air filters, oil filters, fuel filters'),
('Brake Parts', 'Brake pads, brake discs, brake fluid'),
('Engine Parts', 'Various engine components'),
('Electrical Parts', 'Electrical components and accessories'),
('Body Parts', 'Body panels and accessories'),
('Tires', 'Various tire brands and sizes'),
('Tools', 'Workshop tools and equipment');

-- Insert sample suppliers
INSERT INTO suppliers (supplier_code, name, phone, email, address, city) VALUES
('SUP001', 'PT Auto Parts Indonesia', '021-11111111', 'sales@autoparts.co.id', 'Jl. Industri No. 1', 'Jakarta'),
('SUP002', 'CV Sparepart Mobil', '021-22222222', 'info@sparepartmobil.co.id', 'Jl. Otomotif No. 2', 'Jakarta'),
('SUP003', 'Toko Ban Sejahtera', '021-33333333', 'ban@sejahtera.co.id', 'Jl. Ban Raya No. 3', 'Jakarta');

-- Insert sample products
INSERT INTO products (product_code, name, description, category_id, unit_type_id, supplier_id, cost_price, selling_price, stock_quantity, min_stock_level, max_stock_level) VALUES
('PRD001', 'Engine Oil 5W-30 4L', 'Synthetic engine oil 5W-30 4 liter', 1, 2, 1, 120000, 150000, 50, 10, 100),
('PRD002', 'Oil Filter Standard', 'Standard oil filter for most vehicles', 2, 1, 1, 25000, 35000, 100, 20, 200),
('PRD003', 'Air Filter', 'Engine air filter replacement', 2, 1, 1, 40000, 55000, 75, 15, 150),
('PRD004', 'Brake Pad Set Front', 'Front brake pad set for sedan', 3, 5, 2, 200000, 250000, 30, 5, 50),
('PRD005', 'Brake Pad Set Rear', 'Rear brake pad set for sedan', 3, 5, 2, 150000, 200000, 30, 5, 50),
('PRD006', 'Spark Plug Set', '4-piece spark plug set', 4, 5, 1, 80000, 120000, 40, 8, 80);

-- Insert sample customers
INSERT INTO customers (customer_code, name, phone, email, address, city, customer_type, created_by) VALUES
('CUST001', 'John Doe', '081234567890', 'john.doe@email.com', 'Jl. Customer No. 1', 'Jakarta', 'regular', 1),
('CUST002', 'Jane Smith', '081234567891', 'jane.smith@email.com', 'Jl. Customer No. 2', 'Jakarta', 'premium', 1),
('CUST003', 'PT Corporate Client', '021-99999999', 'fleet@corporate.co.id', 'Jl. Corporate No. 1', 'Jakarta', 'vip', 1);

-- Insert sample vehicles
INSERT INTO customer_vehicles (customer_id, vehicle_number, brand, model, year, color, fuel_type, transmission, mileage, created_by) VALUES
(1, 'B1234ABC', 'Toyota', 'Avanza', 2020, 'Silver', 'gasoline', 'manual', 25000, 1),
(2, 'B5678DEF', 'Honda', 'Civic', 2019, 'White', 'gasoline', 'automatic', 35000, 1),
(3, 'B9999GHI', 'Mitsubishi', 'Pajero', 2021, 'Black', 'diesel', 'automatic', 15000, 1);

-- Commit the transaction
COMMIT;
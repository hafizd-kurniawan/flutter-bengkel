-- Seed data for initial setup

-- Insert default outlet
INSERT INTO outlets (name, address, phone, email, is_active) VALUES
('Main Workshop', 'Jl. Raya No. 123, Jakarta', '+62-21-1234567', 'main@bengkel.com', TRUE);

-- Insert default roles
INSERT INTO roles (name, description) VALUES
('Super Admin', 'Full system access'),
('Admin', 'Administrative access'),
('Manager', 'Management access'),
('Technician', 'Technical operations'),
('Cashier', 'Transaction operations'),
('Customer Service', 'Customer service operations');

-- Insert permissions
INSERT INTO permissions (name, description, resource, action) VALUES
-- User management
('users.create', 'Create users', 'users', 'create'),
('users.read', 'View users', 'users', 'read'),
('users.update', 'Update users', 'users', 'update'),
('users.delete', 'Delete users', 'users', 'delete'),

-- Customer management
('customers.create', 'Create customers', 'customers', 'create'),
('customers.read', 'View customers', 'customers', 'read'),
('customers.update', 'Update customers', 'customers', 'update'),
('customers.delete', 'Delete customers', 'customers', 'delete'),

-- Vehicle management
('vehicles.create', 'Create vehicles', 'vehicles', 'create'),
('vehicles.read', 'View vehicles', 'vehicles', 'read'),
('vehicles.update', 'Update vehicles', 'vehicles', 'update'),
('vehicles.delete', 'Delete vehicles', 'vehicles', 'delete'),

-- Service management
('services.create', 'Create services', 'services', 'create'),
('services.read', 'View services', 'services', 'read'),
('services.update', 'Update services', 'services', 'update'),
('services.delete', 'Delete services', 'services', 'delete'),

-- Service job management
('service_jobs.create', 'Create service jobs', 'service_jobs', 'create'),
('service_jobs.read', 'View service jobs', 'service_jobs', 'read'),
('service_jobs.update', 'Update service jobs', 'service_jobs', 'update'),
('service_jobs.delete', 'Delete service jobs', 'service_jobs', 'delete'),

-- Product management
('products.create', 'Create products', 'products', 'create'),
('products.read', 'View products', 'products', 'read'),
('products.update', 'Update products', 'products', 'update'),
('products.delete', 'Delete products', 'products', 'delete'),

-- Transaction management
('transactions.create', 'Create transactions', 'transactions', 'create'),
('transactions.read', 'View transactions', 'transactions', 'read'),
('transactions.update', 'Update transactions', 'transactions', 'update'),
('transactions.delete', 'Delete transactions', 'transactions', 'delete'),

-- Financial management
('financial.read', 'View financial data', 'financial', 'read'),
('financial.update', 'Update financial data', 'financial', 'update'),

-- Reports
('reports.read', 'View reports', 'reports', 'read'),
('reports.generate', 'Generate reports', 'reports', 'generate');

-- Assign all permissions to Super Admin
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- Assign limited permissions to other roles
-- Admin role (ID: 2)
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE name NOT LIKE 'users.%' OR name IN ('users.read', 'users.update');

-- Manager role (ID: 3) 
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE resource IN ('customers', 'vehicles', 'services', 'service_jobs', 'products', 'transactions', 'financial', 'reports');

-- Technician role (ID: 4)
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 4, id FROM permissions WHERE resource IN ('customers', 'vehicles', 'service_jobs', 'products') AND action IN ('read', 'update');

-- Cashier role (ID: 5)
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 5, id FROM permissions WHERE resource IN ('customers', 'transactions', 'products') AND action IN ('read', 'create', 'update');

-- Customer Service role (ID: 6)
INSERT INTO role_has_permissions (role_id, permission_id)
SELECT 6, id FROM permissions WHERE resource IN ('customers', 'vehicles', 'service_jobs') AND action IN ('read', 'create', 'update');

-- Insert default admin user (password: admin123)
INSERT INTO users (username, email, password_hash, full_name, phone, role_id, outlet_id, is_active) VALUES
('admin', 'admin@bengkel.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/L5BKwPQyW', 'System Administrator', '+62-21-9876543', 1, 1, TRUE);

-- Insert unit types
INSERT INTO unit_types (name, abbreviation, description) VALUES
('Piece', 'pcs', 'Individual items'),
('Liter', 'L', 'Liquid measurement'),
('Kilogram', 'kg', 'Weight measurement'),
('Meter', 'm', 'Length measurement'),
('Set', 'set', 'Group of items'),
('Bottle', 'btl', 'Container measurement'),
('Gallon', 'gal', 'Large liquid measurement');

-- Insert service categories
INSERT INTO service_categories (name, description, is_active) VALUES
('Engine Service', 'Engine maintenance and repair', TRUE),
('Transmission Service', 'Transmission maintenance and repair', TRUE),
('Brake Service', 'Brake system maintenance and repair', TRUE),
('AC Service', 'Air conditioning maintenance and repair', TRUE),
('Electrical Service', 'Electrical system maintenance and repair', TRUE),
('Body Work', 'Body repair and painting', TRUE),
('Tire Service', 'Tire installation and repair', TRUE),
('General Maintenance', 'Regular maintenance services', TRUE);

-- Insert product categories
INSERT INTO categories (name, description, is_active) VALUES
('Engine Parts', 'Engine components and parts', TRUE),
('Transmission Parts', 'Transmission components', TRUE),
('Brake Parts', 'Brake system components', TRUE),
('Electrical Parts', 'Electrical components', TRUE),
('Body Parts', 'Vehicle body components', TRUE),
('Fluids & Lubricants', 'Engine oils, coolants, brake fluids', TRUE),
('Filters', 'Air, oil, fuel filters', TRUE),
('Tires & Wheels', 'Tires and wheel components', TRUE),
('Tools & Equipment', 'Workshop tools and equipment', TRUE);

-- Insert payment methods
INSERT INTO payment_methods (name, type, is_active) VALUES
('Cash', 'cash', TRUE),
('Bank Transfer - BCA', 'bank_transfer', TRUE),
('Bank Transfer - Mandiri', 'bank_transfer', TRUE),
('Credit Card', 'credit_card', TRUE),
('Debit Card', 'debit_card', TRUE),
('OVO', 'e_wallet', TRUE),
('GoPay', 'e_wallet', TRUE),
('DANA', 'e_wallet', TRUE);

-- Insert sample services
INSERT INTO services (service_code, name, description, category_id, standard_price, estimated_duration, is_active) VALUES
('SVC001', 'Oil Change', 'Engine oil and filter replacement', 1, 150000, 30, TRUE),
('SVC002', 'Brake Pad Replacement', 'Front brake pad replacement', 3, 300000, 60, TRUE),
('SVC003', 'AC Service', 'Air conditioning cleaning and service', 4, 200000, 45, TRUE),
('SVC004', 'Tire Balancing', 'Wheel balancing service', 7, 100000, 30, TRUE),
('SVC005', 'General Check-up', 'Complete vehicle inspection', 8, 250000, 90, TRUE);

-- Insert sample suppliers
INSERT INTO suppliers (supplier_code, name, email, phone, address, city, contact_person, is_active) VALUES
('SUP001', 'PT Auto Parts Indonesia', 'sales@autoparts.co.id', '+62-21-5555000', 'Jl. Industri No. 45', 'Jakarta', 'Budi Santoso', TRUE),
('SUP002', 'CV Spare Part Motor', 'order@sparepartmotor.com', '+62-21-5555001', 'Jl. Otomotif No. 12', 'Jakarta', 'Siti Rahayu', TRUE),
('SUP003', 'Toko Oli Berkah', 'info@oliberkah.com', '+62-21-5555002', 'Jl. Pelumas No. 8', 'Jakarta', 'Ahmad Fauzi', TRUE);
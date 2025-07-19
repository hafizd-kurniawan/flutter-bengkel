# 🏭 Workshop Management System - Enhanced Setup Guide

## 🚀 Enhanced Features Implemented

### **Kasir-Centric POS System**
- ✅ Multi-payment transaction processing (Cash + Credit + Transfer in single transaction)
- ✅ Product search with barcode scanner integration
- ✅ Real-time stock checking and alerts
- ✅ Auto-receipt printing functionality
- ✅ Touch-optimized interface ready for tablet implementation

### **Service Management Integration**
- ✅ Auto-assign mechanic based on workload and specialty
- ✅ Queue management system with real-time status
- ✅ Service-vehicle trading integration
- ✅ Warranty tracking and service history

### **Vehicle Trading Module**
- ✅ Complete vehicle purchase to sale workflow
- ✅ Service requirement tracking and completion
- ✅ Profit calculation (Selling Price - HPP - Service Cost)
- ✅ Sales team inventory management with advanced filtering

### **Enhanced Security & Audit**
- ✅ Role-based access control (Admin/Manager/Kasir/Sales/Technician)
- ✅ Audit trail with created_by tracking
- ✅ Kasir approval workflow for receivables
- ✅ Enhanced JWT with role and permission claims

## 🔧 Technology Stack

### Backend (Enhanced)
- **Framework**: Fiber v2 with comprehensive middleware
- **Database**: PostgreSQL + SQLX (migrated from MySQL)
- **Authentication**: JWT with Role-Based Access Control
- **Multi-Payment**: Transaction-level payment splitting
- **Audit Trail**: Complete user action tracking

### Database Schema
- **30+ Tables**: Complete business operation coverage
- **PostgreSQL Features**: Triggers, functions, advanced indexing
- **Audit Trail**: created_by tracking across all major entities
- **Multi-Payment**: transaction_payments table for payment splitting

## 📋 Prerequisites

### Required Software
```bash
# 1. PostgreSQL 12+
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS
brew install postgresql
brew services start postgresql

# 2. Go 1.21+
# Download from https://golang.org/dl/

# 3. Git
sudo apt install git  # Ubuntu/Debian
brew install git      # macOS
```

## 🚀 Backend Setup

### 1. Clone and Setup
```bash
git clone <repository-url>
cd flutter-bengkel/backend

# Install dependencies
go mod tidy
```

### 2. Database Setup
```bash
# Connect to PostgreSQL
sudo -u postgres psql

# Create database and user
CREATE DATABASE bengkel_db;
CREATE USER bengkel_user WITH ENCRYPTED PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE bengkel_db TO bengkel_user;
\q
```

### 3. Environment Configuration
```bash
# Copy environment file
cp .env.example .env

# Edit configuration
nano .env
```

Update `.env` with your settings:
```env
# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=bengkel_user
DB_PASSWORD=your_password
DB_NAME=bengkel_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here-make-it-strong
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_HOURS=168

# Server Configuration
PORT=8080
APP_ENV=development
```

### 4. Run Database Migrations
```bash
# Run PostgreSQL migrations in order
psql -h localhost -U bengkel_user -d bengkel_db -f migrations/postgresql/001_foundation_tables.sql
psql -h localhost -U bengkel_user -d bengkel_db -f migrations/postgresql/002_customer_vehicle_tables.sql
psql -h localhost -U bengkel_user -d bengkel_db -f migrations/postgresql/003_master_data_tables.sql
psql -h localhost -U bengkel_user -d bengkel_db -f migrations/postgresql/004_core_operations_tables.sql
psql -h localhost -U bengkel_user -d bengkel_db -f migrations/postgresql/005_financial_tables.sql
psql -h localhost -U bengkel_user -d bengkel_db -f migrations/postgresql/006_seed_data.sql
```

### 5. Start the Server
```bash
# Development with live reload (if air is installed)
air

# Or build and run manually
go build -o main cmd/main.go
./main
```

## 🔑 Default Credentials

```
Username: admin
Password: admin123
```

## 📚 API Documentation

### **POS Operations (Kasir-Centric)**
```
POST   /api/v1/pos/transactions              # Create POS transaction with multi-payment
GET    /api/v1/pos/products/search           # Search products by query/barcode
PUT    /api/v1/pos/transactions/:id/payment  # Add payment method to transaction
POST   /api/v1/pos/transactions/:id/print    # Generate receipt data
GET    /api/v1/pos/queue                     # Get service queue management
PUT    /api/v1/pos/service-jobs/:id/assign   # Auto-assign mechanic
GET    /api/v1/pos/receivables/pending       # Get outstanding receivables (Belum Lunas)
POST   /api/v1/pos/receivables/:id/payment   # Record payment with kasir approval
GET    /api/v1/pos/receivables/paid          # Get paid receivables (Lunas)
GET    /api/v1/pos/dashboard/stats           # Get kasir dashboard statistics
```

### **Vehicle Trading Operations**
```
POST   /api/v1/vehicle-trading/purchase      # Purchase vehicle for trading
PUT    /api/v1/vehicle-trading/:id/service   # Link service requirements
GET    /api/v1/vehicle-trading/inventory     # Get sales inventory with filters
PUT    /api/v1/vehicle-trading/:id/price     # Update selling price (Sales only)
POST   /api/v1/vehicle-trading/sales         # Create vehicle sale transaction
GET    /api/v1/vehicle-trading/sales         # Get vehicle sales history
GET    /api/v1/vehicle-trading/stats         # Get trading statistics
GET    /api/v1/vehicle-trading/:id/profit    # Calculate vehicle profit
PUT    /api/v1/vehicle-trading/:id/complete-service # Mark service as completed
```

## 🎯 Role-Based Access Control

### **Admin** (Full Access)
- All system operations
- User management
- System configuration
- All reports and analytics

### **Manager** (Management Access)
- All operations except user deletion
- Financial approvals
- Comprehensive reporting
- Vehicle trading oversight

### **Kasir** (POS-Centric Access)
- ✅ CREATE: Transactions, Service Jobs, Customers, Vehicles
- ✅ READ: Products (stock view), Customer data, Service history
- ✅ UPDATE: Customer info, Vehicle info, Service status, Payment recording
- ✅ APPROVE: Receivable payments
- ❌ DELETE: No delete permissions
- ❌ UPDATE: Product stock, Pricing (except service pricing)

### **Sales** (Vehicle Trading Focus)
- ✅ READ: Vehicle trading inventory, Service-completed vehicles
- ✅ UPDATE: Vehicle selling prices, Vehicle status
- ✅ CREATE: Vehicle sales transactions
- ✅ FILTER: Advanced inventory filtering
- ❌ CREATE: New vehicle purchases (admin/manager only)

### **Technician** (Service Operations)
- Service job updates and completion
- Product usage for services
- Work queue management
- Service history access

## 🏗️ Multi-Payment Flow Example

```json
{
  "customer_id": 1,
  "discount_amount": 10000,
  "tax_amount": 12000,
  "details": [
    {
      "product_id": 1,
      "quantity": 2,
      "unit_price": 50000
    }
  ],
  "payments": [
    {
      "payment_method_id": 1,
      "amount": 60000
    },
    {
      "payment_method_id": 2, 
      "amount": 42000
    }
  ]
}
```

## 🚗 Vehicle Trading Workflow

### 1. Vehicle Purchase
```json
{
  "vehicle_number": "B1234XYZ",
  "brand": "Toyota",
  "model": "Avanza",
  "year": 2020,
  "purchase_price": 15000000,
  "service_required": true,
  "estimated_selling_price": 18000000
}
```

### 2. Service Integration
- Link service requirements to vehicle
- Track service completion status
- Calculate total cost (purchase + service)

### 3. Sales Process
- Update selling price (sales team)
- Calculate profit margins
- Create sale transaction
- Track commission

## 🎯 Business Logic Features

### **Auto-Assign Mechanic Algorithm**
```
1. Get available technicians with required specialty
2. Check current workload for each technician
3. Assign to technician with lowest workload
4. Log assignment with timestamp and reason
```

### **Multi-Payment Validation**
```
1. Validate sum of all payments equals transaction total
2. Ensure all payment methods are active
3. Record each payment method separately
4. Update transaction payment status accordingly
```

### **Service-Vehicle Trading Integration**
```
1. Vehicle marked as requiring service upon purchase
2. Service job created and linked to vehicle
3. Service completion updates vehicle status to "Available"
4. Total cost calculated as purchase_price + service_cost
5. Minimum selling price enforced as total_cost * 1.1 (10% minimum profit)
```

## 📊 Database Features

### **Audit Trail Implementation**
- All major entities have `created_by` field
- Service job status changes logged in `service_job_histories`
- Payment approvals tracked with kasir ID and timestamp
- Vehicle price changes logged with user and reason

### **Enhanced Indexes**
- Optimized queries for high-frequency operations
- Composite indexes for filtering and searching
- Performance-tuned for POS transactions

### **PostgreSQL Features**
- Triggers for auto-updating timestamps
- Check constraints for data integrity
- JSONB support for flexible reporting parameters

## 🔧 Development Tools

### **Air (Live Reload)**
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with live reload
air
```

### **Database Administration**
```bash
# pgAdmin for GUI management
# Or use command line
psql -h localhost -U bengkel_user -d bengkel_db
```

## 📈 Production Readiness Checklist

- ✅ **Database Migration**: Complete PostgreSQL setup with triggers
- ✅ **Security**: JWT with RBAC, input validation, SQL injection prevention
- ✅ **API Design**: RESTful endpoints with proper HTTP status codes
- ✅ **Error Handling**: Comprehensive error responses
- ✅ **Logging**: Structured logging throughout application
- ✅ **Performance**: Optimized database indexes and queries
- ✅ **Multi-Payment**: Transaction splitting and reconciliation
- ✅ **Audit Trail**: Complete user action tracking

## 🚀 Next Steps

### **Frontend Implementation**
The backend is ready for Flutter integration with:
- Complete API coverage for POS operations
- Role-based endpoint protection  
- Multi-payment transaction support
- Real-time service queue management
- Vehicle trading workflow APIs

### **Production Deployment**
- Docker containerization
- Environment-specific configuration
- SSL/TLS certificate setup
- Database backup and recovery
- Monitoring and logging setup

---

**Built with ❤️ for comprehensive workshop management efficiency**
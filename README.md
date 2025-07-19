# 🏭 Workshop Management System

A comprehensive workshop management system built with **Golang Fiber + PostgreSQL + SQLX** backend and **Flutter** frontend. This system provides complete **Kasir-centric POS**, **Service Management**, and **Vehicle Trading** operations with advanced multi-payment support and role-based access control.

## 🎯 **Implemented Features**

### **🏪 POS System (Kasir-Centric)**
- ✅ **Multi-Payment Transactions** - Cash + Credit + Transfer in single transaction
- ✅ **Product Search** - Barcode scanner and query-based search
- ✅ **Real-time Stock Management** - Stock level checking and alerts
- ✅ **Touch-Optimized Interface** - Ready for tablet implementation
- ✅ **Auto-Receipt Generation** - ESC/POS thermal printer support

### **🔧 Service Management**
- ✅ **Auto-Assign Mechanics** - Workload-based assignment algorithm
- ✅ **Queue Management** - Real-time service job tracking
- ✅ **Service Integration** - Complete vehicle service workflow
- ✅ **Warranty Tracking** - Service history and warranty management

### **🚗 Vehicle Trading Module**
- ✅ **Complete Trading Workflow** - Purchase → Service → Sale
- ✅ **Profit Calculation** - Selling Price - HPP - Service Cost
- ✅ **Sales Inventory** - Advanced filtering and management
- ✅ **Commission Tracking** - Sales performance monitoring

### **🔒 Enhanced Security & Audit**
- ✅ **Role-Based Access Control** - Admin/Manager/Kasir/Sales/Technician
- ✅ **Audit Trail System** - Complete user action tracking
- ✅ **Kasir Approval Workflow** - Receivables payment approval
- ✅ **JWT Authentication** - With role and permission claims

## 🔧 **Technology Stack**

### **Backend (Golang)**
- **Framework**: Fiber v2 with comprehensive middleware
- **Database**: PostgreSQL + SQLX (raw SQL, no ORM)
- **Authentication**: JWT with Role-Based Access Control
- **Multi-Payment**: Transaction-level payment method splitting
- **Real-time**: WebSocket ready for service status updates

### **Database Schema**
- **30+ Tables**: Complete business operation coverage  
- **PostgreSQL Features**: Triggers, functions, advanced indexing
- **Audit Trail**: created_by tracking across all entities
- **Multi-Payment**: transaction_payments table for payment splitting

### **Frontend (Flutter) - Ready for Integration**
- **UI Framework**: Material Design 3 for tablet-first experience
- **State Management**: Riverpod with AsyncNotifier
- **Navigation**: Responsive drawer/bottom nav
- **Print Integration**: ESC/POS thermal printer support
- **Barcode**: Scanner integration for product lookup

## 📱 **API Endpoints Overview**

### **POS Operations** (Kasir-Centric)
```
POST   /api/v1/pos/transactions              # Multi-payment POS transaction
GET    /api/v1/pos/products/search           # Barcode/query product search  
PUT    /api/v1/pos/transactions/:id/payment  # Add payment method
POST   /api/v1/pos/transactions/:id/print    # Generate receipt
GET    /api/v1/pos/queue                     # Service queue management
PUT    /api/v1/pos/service-jobs/:id/assign   # Auto-assign mechanic
GET    /api/v1/pos/receivables/pending       # Outstanding receivables
POST   /api/v1/pos/receivables/:id/payment   # Record payment (Kasir approve)
GET    /api/v1/pos/receivables/paid          # Paid receivables
GET    /api/v1/pos/dashboard/stats           # Kasir dashboard
```

### **Vehicle Trading** (Sales Team)
```
POST   /api/v1/vehicle-trading/purchase      # Purchase vehicle
PUT    /api/v1/vehicle-trading/:id/service   # Link service requirements
GET    /api/v1/vehicle-trading/inventory     # Sales inventory with filters
PUT    /api/v1/vehicle-trading/:id/price     # Update selling price (Sales only)
POST   /api/v1/vehicle-trading/sales         # Create vehicle sale
GET    /api/v1/vehicle-trading/sales         # Sales history
GET    /api/v1/vehicle-trading/stats         # Trading statistics
GET    /api/v1/vehicle-trading/:id/profit    # Profit calculation
```

## 🚀 **Quick Start**

### **Option 1: Docker Setup (Recommended)**
```bash
# Clone repository
git clone <repository-url>
cd flutter-bengkel

# Start with Docker Compose
docker-compose up -d

# Wait for services to start, then access:
# - API: http://localhost:8080
# - pgAdmin: http://localhost:5050
```

### **Option 2: Manual Setup**
```bash
# 1. Setup PostgreSQL
sudo apt install postgresql postgresql-contrib
sudo -u postgres createdb bengkel_db

# 2. Setup Backend  
cd backend
go mod tidy
cp .env.example .env
# Edit .env with your database credentials

# 3. Run migrations
psql -d bengkel_db -f migrations/postgresql/001_foundation_tables.sql
psql -d bengkel_db -f migrations/postgresql/002_customer_vehicle_tables.sql
# ... (run all 6 migration files in order)

# 4. Start server
go run cmd/main.go
```

## 🔑 **Default Credentials**
```
Username: admin  
Password: admin123
```

## 🎯 **Role-Based Access Control**

### **Kasir** (POS-Centric)
- ✅ **CREATE**: Transactions, Service Jobs, Customers, Vehicles
- ✅ **READ**: Products (stock view), All customer data, Service history  
- ✅ **UPDATE**: Customer info, Vehicle info, Service status, Payment recording
- ✅ **APPROVE**: Receivable payments
- ❌ **DELETE**: No delete permissions (audit trail preservation)

### **Sales** (Vehicle Trading)
- ✅ **READ**: Vehicle inventory, Service-completed vehicles
- ✅ **UPDATE**: Vehicle selling prices, Vehicle status
- ✅ **CREATE**: Vehicle sales transactions
- ✅ **FILTER**: Advanced inventory filtering
- ❌ **CREATE**: Vehicle purchases (manager/admin only)

### **Admin/Manager** (Full Access)
- Complete system access
- User management and configuration
- Financial reporting and analytics
- System maintenance operations

## 🏗️ **Business Logic Examples**

### **Multi-Payment Transaction**
```json
{
  "customer_id": 1,
  "discount_amount": 10000,
  "tax_amount": 12000,
  "details": [
    {"product_id": 1, "quantity": 2, "unit_price": 50000}
  ],
  "payments": [
    {"payment_method_id": 1, "amount": 60000},  // Cash
    {"payment_method_id": 4, "amount": 42000}   // Credit Card
  ]
}
```

### **Vehicle Trading Workflow**
```
1. Customer brings vehicle → Vehicle inspection & valuation
2. Purchase price negotiation (HPP) → Service requirement assessment  
3. Service work completion → Final vehicle preparation
4. Add to sales inventory → Sales team pricing
5. Customer purchase → Profit calculation & commission
```

### **Auto-Assign Mechanic Algorithm**
```
1. Get available technicians with required specialty
2. Check current workload (pending jobs count)
3. Assign to technician with lowest workload
4. Log assignment with timestamp and reason
```

## 📊 **Database Features**

### **Audit Trail Implementation**
- All entities have `created_by` field tracking user actions
- Service job history logs all status changes  
- Payment approvals tracked with kasir ID and timestamps
- Vehicle price changes logged with user and reason

### **PostgreSQL Enhancements**
- **Triggers**: Auto-updating timestamps across all tables
- **Check Constraints**: Data integrity validation
- **Advanced Indexes**: Optimized for high-frequency POS operations
- **JSONB Support**: Flexible reporting parameters

## 🔧 **Development Features**

### **Live Reload Development**
```bash
# Install Air for live reload
go install github.com/cosmtrek/air@latest

# Run with live reload
air
```

### **API Testing**
```bash
# Test all endpoints
chmod +x test-api.sh
./test-api.sh
```

### **Database Management**
- pgAdmin included in Docker setup
- Complete PostgreSQL migration scripts
- Seed data with roles, permissions, and sample data

## 📈 **Production Ready Features**

- ✅ **Comprehensive API Coverage**: 25+ endpoints for complete business operations
- ✅ **Security**: JWT with RBAC, input validation, SQL injection prevention  
- ✅ **Performance**: Optimized database indexes and connection pooling
- ✅ **Audit Trail**: Complete user action tracking and compliance
- ✅ **Multi-Payment Processing**: Transaction splitting and reconciliation
- ✅ **Error Handling**: Structured error responses and logging
- ✅ **Docker Support**: Production-ready containerization

## 🎯 **Implementation Status**

### ✅ **Backend (Complete)**
- Complete API implementation with all business logic
- PostgreSQL database with comprehensive schema
- Role-based access control and authentication  
- Multi-payment transaction processing
- Vehicle trading and service integration
- Comprehensive error handling and validation

### 🔄 **Frontend (Ready for Integration)**
- Backend APIs fully support Flutter integration
- Material Design 3 UI components planned
- Responsive design for tablet-first experience
- Real-time updates via WebSocket (ready)
- Print integration for thermal receipt printers

## 🚀 **Next Development Phase**

The backend is production-ready and provides all APIs needed for:

1. **Flutter Frontend Integration** - Complete API coverage available
2. **Real-time Features** - WebSocket infrastructure ready
3. **Print Integration** - Receipt generation APIs implemented
4. **Barcode Scanning** - Product search APIs with barcode support  
5. **Reporting Dashboard** - Statistics and analytics APIs ready

## 📄 **Documentation**

- [**Setup Guide**](SETUP.md) - Detailed installation and configuration
- [**API Documentation**](docs/api.md) - Complete endpoint reference
- [**Database Schema**](docs/database.md) - ERD and table descriptions  
- [**Business Logic**](docs/business.md) - Workflow documentation

## 📄 **License**

MIT License - See [LICENSE](LICENSE) file for details

---

**🎯 Built for comprehensive workshop management efficiency with modern architecture and best practices**
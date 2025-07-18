# Workshop Management System Backend

A comprehensive workshop management system built with Go Fiber and SQLX.

## Features

- **Multi-Business Support**: Service workshop, sparepart sales, and vehicle trading
- **User Management**: Role-based access control with JWT authentication
- **Customer Management**: Complete customer and vehicle registration
- **Service Jobs**: Queue management and service tracking
- **Inventory Management**: Product catalog with stock tracking
- **Financial Management**: Transaction processing and payment tracking
- **Reporting**: Business analytics and reporting

## Tech Stack

- **Framework**: Fiber v2 (High-performance HTTP framework)
- **Database**: SQLX with MySQL (Raw SQL queries)
- **Authentication**: JWT with RBAC system
- **Documentation**: Swagger/OpenAPI
- **Live Reload**: Air for development

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up your database configuration in `.env`:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. Run database migrations:
   ```bash
   go run cmd/main.go
   ```

## Development

Run with live reload:
```bash
air
```

## API Documentation

Once the server is running, visit:
- Swagger UI: http://localhost:8080/swagger/
- Health check: http://localhost:8080/health

## Database Schema

The system includes 25+ tables covering:
- Foundation & Security (users, roles, permissions, outlets)
- Customer & Vehicle Management
- Master Data & Inventory
- Core Operations (service jobs, transactions)
- Financial & Reporting

## Default Credentials

- Username: `admin`
- Password: `admin123`

## Project Structure

```
backend/
├── cmd/                    # Application entry points
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection & migrations
│   ├── models/            # Data models
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Authentication & middleware
│   ├── services/          # Business logic
│   ├── repositories/      # SQLX data access
│   └── utils/             # Utility functions
├── migrations/            # Database migrations
└── docs/                 # API documentation
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Logout

### Core Resources
- `/api/v1/users` - User management
- `/api/v1/customers` - Customer management
- `/api/v1/vehicles` - Vehicle management
- `/api/v1/services` - Service catalog
- `/api/v1/products` - Product inventory
- `/api/v1/service-jobs` - Service job management
- `/api/v1/transactions` - Transaction handling
- `/api/v1/payments` - Payment processing

### Master Data
- `/api/v1/master-data/service-categories`
- `/api/v1/master-data/product-categories`
- `/api/v1/master-data/suppliers`
- `/api/v1/master-data/unit-types`
- `/api/v1/master-data/payment-methods`

## License

MIT License
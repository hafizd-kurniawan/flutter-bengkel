package repositories

import (
	"github.com/jmoiron/sqlx"
)

// Repositories contains all repositories
type Repositories struct {
	User         UserRepository
	Role         RoleRepository
	Permission   PermissionRepository
	Outlet       OutletRepository
	Customer     CustomerRepository
	Vehicle      VehicleRepository
	Service      ServiceRepository
	Product      ProductRepository
	ServiceJob   ServiceJobRepository
	Transaction  TransactionRepository
	Payment      PaymentRepository
}

// New creates a new repositories instance
func New(db *sqlx.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Role:         NewRoleRepository(db),
		Permission:   NewPermissionRepository(db),
		Outlet:       NewOutletRepository(db),
		Customer:     NewCustomerRepository(db),
		Vehicle:      NewVehicleRepository(db),
		Service:      NewServiceRepository(db),
		Product:      NewProductRepository(db),
		ServiceJob:   NewServiceJobRepository(db),
		Transaction:  NewTransactionRepository(db),
		Payment:      NewPaymentRepository(db),
	}
}
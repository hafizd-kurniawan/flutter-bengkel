package services

import (
	"flutter-bengkel/internal/config"
	"flutter-bengkel/internal/repositories"
)

// Services contains all application services
type Services struct {
	Auth           AuthService
	User           UserService
	Customer       CustomerService
	Vehicle        VehicleService
	Service        ServiceService
	Product        ProductService
	ServiceJob     ServiceJobService
	Transaction    TransactionService
	Payment        PaymentService
	Receivable     ReceivableService     // NEW
	VehicleTrading VehicleTradingService // NEW
	Dashboard      DashboardService      // NEW
}

// New creates a new services instance
func New(repos *repositories.Repositories, cfg *config.Config) *Services {
	return &Services{
		Auth:           NewAuthService(repos, cfg),
		User:           NewUserService(repos),
		Customer:       NewCustomerService(repos),
		Vehicle:        NewVehicleService(repos),
		Service:        NewServiceService(repos),
		Product:        NewProductService(repos),
		ServiceJob:     NewServiceJobService(repos),
		Transaction:    NewTransactionService(repos),
		Payment:        NewPaymentService(repos),
		Receivable:     NewReceivableService(repos),     // NEW
		VehicleTrading: NewVehicleTradingService(repos), // NEW
		Dashboard:      NewDashboardService(repos),      // NEW
	}
}
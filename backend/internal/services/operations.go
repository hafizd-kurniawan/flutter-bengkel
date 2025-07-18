package services

import (
	"errors"
	"time"

	"flutter-bengkel/internal/models"
	"flutter-bengkel/internal/repositories"
)

// ServiceJob Service
type ServiceJobService interface {
	Create(req *models.CreateServiceJobRequest, outletID int64, userID int64) (*models.ServiceJob, error)
	GetByID(id int64) (*models.ServiceJob, error)
	Update(id int64, req *models.UpdateServiceJobRequest) (*models.ServiceJob, error)
	UpdateStatus(id int64, status string, userID int64, notes string) error
	Delete(id int64) error
	List(page, limit int, outletID *int64, status string, search string) ([]models.ServiceJob, *models.PaginationMeta, error)
	AddDetail(serviceJobID int64, detail *models.ServiceDetail) (*models.ServiceDetail, error)
	GetDetails(serviceJobID int64) ([]models.ServiceDetail, error)
	UpdateDetail(detailID int64, detail *models.ServiceDetail) error
	DeleteDetail(detailID int64) error
	CalculateTotal(serviceJobID int64) error
}

type serviceJobService struct {
	repos *repositories.Repositories
}

func NewServiceJobService(repos *repositories.Repositories) ServiceJobService {
	return &serviceJobService{repos: repos}
}

func (s *serviceJobService) Create(req *models.CreateServiceJobRequest, outletID int64, userID int64) (*models.ServiceJob, error) {
	// Validate customer exists
	if _, err := s.repos.Customer.GetByID(req.CustomerID); err != nil {
		return nil, errors.New("customer not found")
	}

	// Validate vehicle exists and belongs to customer
	vehicle, err := s.repos.Vehicle.GetByID(req.VehicleID)
	if err != nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.CustomerID != req.CustomerID {
		return nil, errors.New("vehicle does not belong to customer")
	}

	// Validate technician if provided
	if req.TechnicianID != nil {
		if _, err := s.repos.User.GetByID(*req.TechnicianID); err != nil {
			return nil, errors.New("technician not found")
		}
	}

	// Generate job number
	jobNumber, err := s.repos.ServiceJob.GenerateJobNumber()
	if err != nil {
		return nil, err
	}

	// Get next queue number
	queueNumber, err := s.repos.ServiceJob.GetNextQueueNumber(outletID)
	if err != nil {
		return nil, err
	}

	// Set defaults
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	serviceJob := &models.ServiceJob{
		JobNumber:          jobNumber,
		CustomerID:         req.CustomerID,
		VehicleID:          req.VehicleID,
		OutletID:           outletID,
		TechnicianID:       req.TechnicianID,
		QueueNumber:        queueNumber,
		Priority:           priority,
		Status:             "pending",
		ProblemDescription: req.ProblemDescription,
		WarrantyPeriodDays: req.WarrantyPeriodDays,
		Notes:              req.Notes,
	}

	if err := s.repos.ServiceJob.Create(serviceJob); err != nil {
		return nil, err
	}

	return s.repos.ServiceJob.GetByID(serviceJob.ID)
}

func (s *serviceJobService) GetByID(id int64) (*models.ServiceJob, error) {
	serviceJob, err := s.repos.ServiceJob.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get details
	details, err := s.repos.ServiceJob.GetDetails(id)
	if err == nil {
		serviceJob.Details = details
	}

	return serviceJob, nil
}

func (s *serviceJobService) Update(id int64, req *models.UpdateServiceJobRequest) (*models.ServiceJob, error) {
	// Get existing service job
	existingServiceJob, err := s.repos.ServiceJob.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Validate technician if provided
	if req.TechnicianID != nil {
		if _, err := s.repos.User.GetByID(*req.TechnicianID); err != nil {
			return nil, errors.New("technician not found")
		}
	}

	// Update fields
	serviceJob := &models.ServiceJob{
		TechnicianID:        req.TechnicianID,
		Priority:            req.Priority,
		Status:              req.Status,
		EstimatedCompletion: req.EstimatedCompletion,
		ActualCompletion:    req.ActualCompletion,
		WarrantyPeriodDays:  req.WarrantyPeriodDays,
		Notes:               req.Notes,
		// Keep existing totals - these should be calculated separately
		TotalAmount:    existingServiceJob.TotalAmount,
		DiscountAmount: existingServiceJob.DiscountAmount,
		TaxAmount:      existingServiceJob.TaxAmount,
		FinalAmount:    existingServiceJob.FinalAmount,
	}

	if err := s.repos.ServiceJob.Update(id, serviceJob); err != nil {
		return nil, err
	}

	return s.repos.ServiceJob.GetByID(id)
}

func (s *serviceJobService) UpdateStatus(id int64, status string, userID int64, notes string) error {
	// Validate status
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled", "on_hold"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid status")
	}

	return s.repos.ServiceJob.UpdateStatus(id, status, userID, notes)
}

func (s *serviceJobService) Delete(id int64) error {
	return s.repos.ServiceJob.Delete(id)
}

func (s *serviceJobService) List(page, limit int, outletID *int64, status string, search string) ([]models.ServiceJob, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	serviceJobs, total, err := s.repos.ServiceJob.List(offset, limit, outletID, status, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return serviceJobs, meta, nil
}

func (s *serviceJobService) AddDetail(serviceJobID int64, detail *models.ServiceDetail) (*models.ServiceDetail, error) {
	// Validate service job exists
	if _, err := s.repos.ServiceJob.GetByID(serviceJobID); err != nil {
		return nil, errors.New("service job not found")
	}

	// Validate product or service exists
	if detail.ProductID != nil {
		if _, err := s.repos.Product.GetByID(*detail.ProductID); err != nil {
			return nil, errors.New("product not found")
		}
	}
	if detail.ServiceID != nil {
		if _, err := s.repos.Service.GetByID(*detail.ServiceID); err != nil {
			return nil, errors.New("service not found")
		}
	}

	// Either product or service must be specified
	if detail.ProductID == nil && detail.ServiceID == nil {
		return nil, errors.New("either product or service must be specified")
	}

	detail.ServiceJobID = serviceJobID
	detail.TotalPrice = detail.Quantity * detail.UnitPrice

	if err := s.repos.ServiceJob.AddDetail(detail); err != nil {
		return nil, err
	}

	// Recalculate totals
	if err := s.CalculateTotal(serviceJobID); err != nil {
		// Log error but don't fail the operation
	}

	return detail, nil
}

func (s *serviceJobService) GetDetails(serviceJobID int64) ([]models.ServiceDetail, error) {
	return s.repos.ServiceJob.GetDetails(serviceJobID)
}

func (s *serviceJobService) UpdateDetail(detailID int64, detail *models.ServiceDetail) error {
	detail.TotalPrice = detail.Quantity * detail.UnitPrice

	if err := s.repos.ServiceJob.UpdateDetail(detailID, detail); err != nil {
		return err
	}

	// Get the service job ID and recalculate totals
	// This would require getting the detail first to get the service job ID
	// For now, we'll skip the automatic recalculation

	return nil
}

func (s *serviceJobService) DeleteDetail(detailID int64) error {
	return s.repos.ServiceJob.DeleteDetail(detailID)
}

func (s *serviceJobService) CalculateTotal(serviceJobID int64) error {
	// Get all details
	details, err := s.repos.ServiceJob.GetDetails(serviceJobID)
	if err != nil {
		return err
	}

	// Calculate total
	var totalAmount float64
	for _, detail := range details {
		totalAmount += detail.TotalPrice
	}

	// Get current service job to preserve discount and tax
	serviceJob, err := s.repos.ServiceJob.GetByID(serviceJobID)
	if err != nil {
		return err
	}

	// Update totals
	serviceJob.TotalAmount = totalAmount
	serviceJob.FinalAmount = totalAmount - serviceJob.DiscountAmount + serviceJob.TaxAmount

	return s.repos.ServiceJob.Update(serviceJobID, serviceJob)
}

// Transaction Service
type TransactionService interface {
	Create(req *models.CreateTransactionRequest, outletID int64, userID int64) (*models.Transaction, error)
	GetByID(id int64) (*models.Transaction, error)
	Update(id int64, req *models.Transaction) (*models.Transaction, error)
	Delete(id int64) error
	List(page, limit int, outletID *int64, transactionType string, search string) ([]models.Transaction, *models.PaginationMeta, error)
	UpdatePaymentStatus(id int64, status string) error
}

type transactionService struct {
	repos *repositories.Repositories
}

func NewTransactionService(repos *repositories.Repositories) TransactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) Create(req *models.CreateTransactionRequest, outletID int64, userID int64) (*models.Transaction, error) {
	// Validate customer if provided
	if req.CustomerID != nil {
		if _, err := s.repos.Customer.GetByID(*req.CustomerID); err != nil {
			return nil, errors.New("customer not found")
		}
	}

	// Validate service job if provided
	if req.ServiceJobID != nil {
		if _, err := s.repos.ServiceJob.GetByID(*req.ServiceJobID); err != nil {
			return nil, errors.New("service job not found")
		}
	}

	// Generate transaction number
	transactionNumber, err := s.repos.Transaction.GenerateTransactionNumber(req.TransactionType)
	if err != nil {
		return nil, err
	}

	// Calculate subtotal from details
	var subtotalAmount float64
	for _, detail := range req.Details {
		// Validate product or service
		if detail.ProductID != nil {
			if _, err := s.repos.Product.GetByID(*detail.ProductID); err != nil {
				return nil, errors.New("product not found")
			}
		}
		if detail.ServiceID != nil {
			if _, err := s.repos.Service.GetByID(*detail.ServiceID); err != nil {
				return nil, errors.New("service not found")
			}
		}

		totalPrice := detail.Quantity * detail.UnitPrice
		subtotalAmount += totalPrice
	}

	// Calculate total amount
	totalAmount := subtotalAmount - req.DiscountAmount + req.TaxAmount

	transaction := &models.Transaction{
		TransactionNumber: transactionNumber,
		TransactionType:   req.TransactionType,
		CustomerID:        req.CustomerID,
		OutletID:          outletID,
		UserID:            userID,
		ServiceJobID:      req.ServiceJobID,
		SubtotalAmount:    subtotalAmount,
		DiscountAmount:    req.DiscountAmount,
		TaxAmount:         req.TaxAmount,
		TotalAmount:       totalAmount,
		PaymentStatus:     "pending",
		Notes:             req.Notes,
		TransactionDate:   time.Now(),
	}

	if err := s.repos.Transaction.Create(transaction); err != nil {
		return nil, err
	}

	// Add transaction details
	for _, detailReq := range req.Details {
		detail := &models.TransactionDetail{
			TransactionID: transaction.ID,
			ProductID:     detailReq.ProductID,
			ServiceID:     detailReq.ServiceID,
			Description:   detailReq.Description,
			Quantity:      detailReq.Quantity,
			UnitPrice:     detailReq.UnitPrice,
			TotalPrice:    detailReq.Quantity * detailReq.UnitPrice,
		}

		if err := s.repos.Transaction.AddDetail(detail); err != nil {
			return nil, err
		}
	}

	return s.repos.Transaction.GetByID(transaction.ID)
}

func (s *transactionService) GetByID(id int64) (*models.Transaction, error) {
	transaction, err := s.repos.Transaction.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get details
	details, err := s.repos.Transaction.GetDetails(id)
	if err == nil {
		transaction.Details = details
	}

	// Get payments
	payments, err := s.repos.Payment.GetByTransactionID(id)
	if err == nil {
		transaction.Payments = payments
	}

	return transaction, nil
}

func (s *transactionService) Update(id int64, req *models.Transaction) (*models.Transaction, error) {
	if err := s.repos.Transaction.Update(id, req); err != nil {
		return nil, err
	}

	return s.repos.Transaction.GetByID(id)
}

func (s *transactionService) Delete(id int64) error {
	return s.repos.Transaction.Delete(id)
}

func (s *transactionService) List(page, limit int, outletID *int64, transactionType string, search string) ([]models.Transaction, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	transactions, total, err := s.repos.Transaction.List(offset, limit, outletID, transactionType, search)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return transactions, meta, nil
}

func (s *transactionService) UpdatePaymentStatus(id int64, status string) error {
	validStatuses := []string{"pending", "partial", "paid", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("invalid payment status")
	}

	return s.repos.Transaction.UpdatePaymentStatus(id, status)
}

// Payment Service
type PaymentService interface {
	Create(req *models.CreatePaymentRequest) (*models.Payment, error)
	GetByID(id int64) (*models.Payment, error)
	Delete(id int64) error
	List(page, limit int, transactionID *int64) ([]models.Payment, *models.PaginationMeta, error)
	GetByTransactionID(transactionID int64) ([]models.Payment, error)
	ListPaymentMethods() ([]models.PaymentMethod, error)
}

type paymentService struct {
	repos *repositories.Repositories
}

func NewPaymentService(repos *repositories.Repositories) PaymentService {
	return &paymentService{repos: repos}
}

func (s *paymentService) Create(req *models.CreatePaymentRequest) (*models.Payment, error) {
	// Validate transaction exists
	transaction, err := s.repos.Transaction.GetByID(req.TransactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Validate payment amount doesn't exceed remaining amount
	existingPayments, err := s.repos.Payment.GetByTransactionID(req.TransactionID)
	if err != nil {
		return nil, err
	}

	var totalPaid float64
	for _, payment := range existingPayments {
		totalPaid += payment.Amount
	}

	if totalPaid+req.Amount > transaction.TotalAmount {
		return nil, errors.New("payment amount exceeds remaining amount")
	}

	// Generate payment number
	paymentNumber, err := s.repos.Payment.GeneratePaymentNumber()
	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		PaymentNumber:   paymentNumber,
		TransactionID:   req.TransactionID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
		PaymentDate:     time.Now(),
		ReferenceNumber: req.ReferenceNumber,
		Notes:           req.Notes,
	}

	if err := s.repos.Payment.Create(payment); err != nil {
		return nil, err
	}

	// Update transaction payment status
	newTotalPaid := totalPaid + req.Amount
	var status string
	if newTotalPaid >= transaction.TotalAmount {
		status = "paid"
	} else if newTotalPaid > 0 {
		status = "partial"
	} else {
		status = "pending"
	}

	if err := s.repos.Transaction.UpdatePaymentStatus(req.TransactionID, status); err != nil {
		// Log error but don't fail the payment creation
	}

	return s.repos.Payment.GetByID(payment.ID)
}

func (s *paymentService) GetByID(id int64) (*models.Payment, error) {
	return s.repos.Payment.GetByID(id)
}

func (s *paymentService) Delete(id int64) error {
	return s.repos.Payment.Delete(id)
}

func (s *paymentService) List(page, limit int, transactionID *int64) ([]models.Payment, *models.PaginationMeta, error) {
	offset := (page - 1) * limit
	payments, total, err := s.repos.Payment.List(offset, limit, transactionID)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	meta := &models.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return payments, meta, nil
}

func (s *paymentService) GetByTransactionID(transactionID int64) ([]models.Payment, error) {
	return s.repos.Payment.GetByTransactionID(transactionID)
}

func (s *paymentService) ListPaymentMethods() ([]models.PaymentMethod, error) {
	return s.repos.Payment.ListPaymentMethods()
}
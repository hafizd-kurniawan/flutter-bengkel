package repositories

import (
	"fmt"

	"flutter-bengkel/internal/models"

	"github.com/jmoiron/sqlx"
)

// ServiceJob Repository
type ServiceJobRepository interface {
	Create(serviceJob *models.ServiceJob) error
	GetByID(id int64) (*models.ServiceJob, error)
	GetByJobNumber(jobNumber string) (*models.ServiceJob, error)
	Update(id int64, serviceJob *models.ServiceJob) error
	UpdateStatus(id int64, status string, userID int64, notes string) error
	Delete(id int64) error
	List(offset, limit int, outletID *int64, status string, search string) ([]models.ServiceJob, int64, error)
	GetNextQueueNumber(outletID int64) (int, error)
	GenerateJobNumber() (string, error)
	AddDetail(detail *models.ServiceDetail) error
	GetDetails(serviceJobID int64) ([]models.ServiceDetail, error)
	UpdateDetail(id int64, detail *models.ServiceDetail) error
	DeleteDetail(id int64) error
}

type serviceJobRepository struct {
	db *sqlx.DB
}

func NewServiceJobRepository(db *sqlx.DB) ServiceJobRepository {
	return &serviceJobRepository{db: db}
}

func (r *serviceJobRepository) Create(serviceJob *models.ServiceJob) error {
	query := `
		INSERT INTO service_jobs (job_number, customer_id, vehicle_id, outlet_id, technician_id,
								 queue_number, priority, status, problem_description, 
								 estimated_completion, total_amount, discount_amount, 
								 tax_amount, final_amount, warranty_period_days, notes)
		VALUES (:job_number, :customer_id, :vehicle_id, :outlet_id, :technician_id,
				:queue_number, :priority, :status, :problem_description, 
				:estimated_completion, :total_amount, :discount_amount, 
				:tax_amount, :final_amount, :warranty_period_days, :notes)
	`
	
	result, err := r.db.NamedExec(query, serviceJob)
	if err != nil {
		return fmt.Errorf("failed to create service job: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get service job ID: %w", err)
	}
	
	serviceJob.ID = id
	return nil
}

func (r *serviceJobRepository) GetByID(id int64) (*models.ServiceJob, error) {
	query := `
		SELECT sj.id, sj.job_number, sj.customer_id, sj.vehicle_id, sj.outlet_id, sj.technician_id,
			   sj.queue_number, sj.priority, sj.status, sj.problem_description, 
			   sj.estimated_completion, sj.actual_completion, sj.total_amount, sj.discount_amount, 
			   sj.tax_amount, sj.final_amount, sj.warranty_period_days, sj.notes, 
			   sj.created_at, sj.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone",
			   cv.id as "vehicle.id", cv.vehicle_number as "vehicle.vehicle_number",
			   cv.brand as "vehicle.brand", cv.model as "vehicle.model",
			   o.id as "outlet.id", o.name as "outlet.name",
			   u.id as "technician.id", u.full_name as "technician.full_name"
		FROM service_jobs sj
		LEFT JOIN customers c ON sj.customer_id = c.id
		LEFT JOIN customer_vehicles cv ON sj.vehicle_id = cv.id
		LEFT JOIN outlets o ON sj.outlet_id = o.id
		LEFT JOIN users u ON sj.technician_id = u.id
		WHERE sj.id = ?
	`
	
	var serviceJob models.ServiceJob
	err := r.db.Get(&serviceJob, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service job: %w", err)
	}
	
	return &serviceJob, nil
}

func (r *serviceJobRepository) GetByJobNumber(jobNumber string) (*models.ServiceJob, error) {
	query := `
		SELECT sj.id, sj.job_number, sj.customer_id, sj.vehicle_id, sj.outlet_id, sj.technician_id,
			   sj.queue_number, sj.priority, sj.status, sj.problem_description, 
			   sj.estimated_completion, sj.actual_completion, sj.total_amount, sj.discount_amount, 
			   sj.tax_amount, sj.final_amount, sj.warranty_period_days, sj.notes, 
			   sj.created_at, sj.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone",
			   cv.id as "vehicle.id", cv.vehicle_number as "vehicle.vehicle_number",
			   cv.brand as "vehicle.brand", cv.model as "vehicle.model",
			   o.id as "outlet.id", o.name as "outlet.name",
			   u.id as "technician.id", u.full_name as "technician.full_name"
		FROM service_jobs sj
		LEFT JOIN customers c ON sj.customer_id = c.id
		LEFT JOIN customer_vehicles cv ON sj.vehicle_id = cv.id
		LEFT JOIN outlets o ON sj.outlet_id = o.id
		LEFT JOIN users u ON sj.technician_id = u.id
		WHERE sj.job_number = ?
	`
	
	var serviceJob models.ServiceJob
	err := r.db.Get(&serviceJob, query, jobNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get service job by number: %w", err)
	}
	
	return &serviceJob, nil
}

func (r *serviceJobRepository) Update(id int64, serviceJob *models.ServiceJob) error {
	query := `
		UPDATE service_jobs 
		SET technician_id = :technician_id, priority = :priority, status = :status,
			estimated_completion = :estimated_completion, actual_completion = :actual_completion,
			total_amount = :total_amount, discount_amount = :discount_amount, 
			tax_amount = :tax_amount, final_amount = :final_amount, 
			warranty_period_days = :warranty_period_days, notes = :notes
		WHERE id = :id
	`
	
	serviceJob.ID = id
	_, err := r.db.NamedExec(query, serviceJob)
	if err != nil {
		return fmt.Errorf("failed to update service job: %w", err)
	}
	
	return nil
}

func (r *serviceJobRepository) UpdateStatus(id int64, status string, userID int64, notes string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Get current status
	var currentStatus string
	err = tx.Get(&currentStatus, "SELECT status FROM service_jobs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to get current status: %w", err)
	}
	
	// Update status
	_, err = tx.Exec("UPDATE service_jobs SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	
	// Add history record
	_, err = tx.Exec(`
		INSERT INTO service_job_histories (service_job_id, user_id, previous_status, new_status, notes)
		VALUES (?, ?, ?, ?, ?)
	`, id, userID, currentStatus, status, notes)
	if err != nil {
		return fmt.Errorf("failed to create history: %w", err)
	}
	
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (r *serviceJobRepository) Delete(id int64) error {
	query := `UPDATE service_jobs SET status = 'cancelled' WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service job: %w", err)
	}
	
	return nil
}

func (r *serviceJobRepository) List(offset, limit int, outletID *int64, status string, search string) ([]models.ServiceJob, int64, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	
	if outletID != nil {
		whereClause += " AND sj.outlet_id = ?"
		args = append(args, *outletID)
	}
	
	if status != "" {
		whereClause += " AND sj.status = ?"
		args = append(args, status)
	}
	
	if search != "" {
		whereClause += " AND (sj.job_number LIKE ? OR c.name LIKE ? OR cv.vehicle_number LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}
	
	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM service_jobs sj 
		LEFT JOIN customers c ON sj.customer_id = c.id
		LEFT JOIN customer_vehicles cv ON sj.vehicle_id = cv.id
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count service jobs: %w", err)
	}
	
	// Get service jobs
	query := fmt.Sprintf(`
		SELECT sj.id, sj.job_number, sj.customer_id, sj.vehicle_id, sj.outlet_id, sj.technician_id,
			   sj.queue_number, sj.priority, sj.status, sj.problem_description, 
			   sj.estimated_completion, sj.actual_completion, sj.total_amount, sj.discount_amount, 
			   sj.tax_amount, sj.final_amount, sj.warranty_period_days, sj.notes, 
			   sj.created_at, sj.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone",
			   cv.id as "vehicle.id", cv.vehicle_number as "vehicle.vehicle_number",
			   cv.brand as "vehicle.brand", cv.model as "vehicle.model",
			   o.id as "outlet.id", o.name as "outlet.name",
			   u.id as "technician.id", u.full_name as "technician.full_name"
		FROM service_jobs sj
		LEFT JOIN customers c ON sj.customer_id = c.id
		LEFT JOIN customer_vehicles cv ON sj.vehicle_id = cv.id
		LEFT JOIN outlets o ON sj.outlet_id = o.id
		LEFT JOIN users u ON sj.technician_id = u.id
		%s
		ORDER BY sj.queue_number ASC, sj.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var serviceJobs []models.ServiceJob
	err = r.db.Select(&serviceJobs, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list service jobs: %w", err)
	}
	
	return serviceJobs, total, nil
}

func (r *serviceJobRepository) GetNextQueueNumber(outletID int64) (int, error) {
	query := `
		SELECT COALESCE(MAX(queue_number), 0) + 1 
		FROM service_jobs 
		WHERE outlet_id = ? AND DATE(created_at) = CURDATE() AND status != 'cancelled'
	`
	
	var queueNumber int
	err := r.db.Get(&queueNumber, query, outletID)
	if err != nil {
		return 0, fmt.Errorf("failed to get next queue number: %w", err)
	}
	
	return queueNumber, nil
}

func (r *serviceJobRepository) GenerateJobNumber() (string, error) {
	query := `SELECT COUNT(*) FROM service_jobs WHERE job_number LIKE CONCAT('SJ', DATE_FORMAT(NOW(), '%Y%m%d'), '%')`
	
	var count int
	err := r.db.Get(&count, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate job number: %w", err)
	}
	
	return fmt.Sprintf("SJ%s%04d", 
		fmt.Sprintf("%04d%02d%02d", 2024, 1, 1), // This should use current date
		count+1), nil
}

func (r *serviceJobRepository) AddDetail(detail *models.ServiceDetail) error {
	query := `
		INSERT INTO service_details (service_job_id, product_id, service_id, quantity, unit_price, total_price, notes)
		VALUES (:service_job_id, :product_id, :service_id, :quantity, :unit_price, :total_price, :notes)
	`
	
	result, err := r.db.NamedExec(query, detail)
	if err != nil {
		return fmt.Errorf("failed to add service detail: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get service detail ID: %w", err)
	}
	
	detail.ID = id
	return nil
}

func (r *serviceJobRepository) GetDetails(serviceJobID int64) ([]models.ServiceDetail, error) {
	query := `
		SELECT sd.id, sd.service_job_id, sd.product_id, sd.service_id, sd.quantity, 
			   sd.unit_price, sd.total_price, sd.notes, sd.created_at,
			   p.id as "product.id", p.product_code as "product.product_code", 
			   p.name as "product.name",
			   s.id as "service.id", s.service_code as "service.service_code", 
			   s.name as "service.name"
		FROM service_details sd
		LEFT JOIN products p ON sd.product_id = p.id
		LEFT JOIN services s ON sd.service_id = s.id
		WHERE sd.service_job_id = ?
		ORDER BY sd.created_at
	`
	
	var details []models.ServiceDetail
	err := r.db.Select(&details, query, serviceJobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service details: %w", err)
	}
	
	return details, nil
}

func (r *serviceJobRepository) UpdateDetail(id int64, detail *models.ServiceDetail) error {
	query := `
		UPDATE service_details 
		SET quantity = :quantity, unit_price = :unit_price, total_price = :total_price, notes = :notes
		WHERE id = :id
	`
	
	detail.ID = id
	_, err := r.db.NamedExec(query, detail)
	if err != nil {
		return fmt.Errorf("failed to update service detail: %w", err)
	}
	
	return nil
}

func (r *serviceJobRepository) DeleteDetail(id int64) error {
	query := `DELETE FROM service_details WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service detail: %w", err)
	}
	
	return nil
}

// Transaction Repository
type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id int64) (*models.Transaction, error)
	GetByTransactionNumber(transactionNumber string) (*models.Transaction, error)
	Update(id int64, transaction *models.Transaction) error
	Delete(id int64) error
	List(offset, limit int, outletID *int64, transactionType string, search string) ([]models.Transaction, int64, error)
	GenerateTransactionNumber(transactionType string) (string, error)
	AddDetail(detail *models.TransactionDetail) error
	GetDetails(transactionID int64) ([]models.TransactionDetail, error)
	UpdateDetail(id int64, detail *models.TransactionDetail) error
	DeleteDetail(id int64) error
	UpdatePaymentStatus(id int64, status string) error
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (transaction_number, transaction_type, customer_id, outlet_id, 
								  user_id, service_job_id, subtotal_amount, discount_amount, 
								  tax_amount, total_amount, payment_status, notes, transaction_date)
		VALUES (:transaction_number, :transaction_type, :customer_id, :outlet_id, 
				:user_id, :service_job_id, :subtotal_amount, :discount_amount, 
				:tax_amount, :total_amount, :payment_status, :notes, :transaction_date)
	`
	
	result, err := r.db.NamedExec(query, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get transaction ID: %w", err)
	}
	
	transaction.ID = id
	return nil
}

func (r *transactionRepository) GetByID(id int64) (*models.Transaction, error) {
	query := `
		SELECT t.id, t.transaction_number, t.transaction_type, t.customer_id, t.outlet_id, 
			   t.user_id, t.service_job_id, t.subtotal_amount, t.discount_amount, 
			   t.tax_amount, t.total_amount, t.payment_status, t.notes, 
			   t.transaction_date, t.created_at, t.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone",
			   o.id as "outlet.id", o.name as "outlet.name",
			   u.id as "user.id", u.full_name as "user.full_name",
			   sj.id as "service_job.id", sj.job_number as "service_job.job_number"
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		LEFT JOIN outlets o ON t.outlet_id = o.id
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN service_jobs sj ON t.service_job_id = sj.id
		WHERE t.id = ?
	`
	
	var transaction models.Transaction
	err := r.db.Get(&transaction, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	
	return &transaction, nil
}

func (r *transactionRepository) GetByTransactionNumber(transactionNumber string) (*models.Transaction, error) {
	query := `
		SELECT t.id, t.transaction_number, t.transaction_type, t.customer_id, t.outlet_id, 
			   t.user_id, t.service_job_id, t.subtotal_amount, t.discount_amount, 
			   t.tax_amount, t.total_amount, t.payment_status, t.notes, 
			   t.transaction_date, t.created_at, t.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone",
			   o.id as "outlet.id", o.name as "outlet.name",
			   u.id as "user.id", u.full_name as "user.full_name",
			   sj.id as "service_job.id", sj.job_number as "service_job.job_number"
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		LEFT JOIN outlets o ON t.outlet_id = o.id
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN service_jobs sj ON t.service_job_id = sj.id
		WHERE t.transaction_number = ?
	`
	
	var transaction models.Transaction
	err := r.db.Get(&transaction, query, transactionNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by number: %w", err)
	}
	
	return &transaction, nil
}

func (r *transactionRepository) Update(id int64, transaction *models.Transaction) error {
	query := `
		UPDATE transactions 
		SET subtotal_amount = :subtotal_amount, discount_amount = :discount_amount, 
			tax_amount = :tax_amount, total_amount = :total_amount, 
			payment_status = :payment_status, notes = :notes
		WHERE id = :id
	`
	
	transaction.ID = id
	_, err := r.db.NamedExec(query, transaction)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	
	return nil
}

func (r *transactionRepository) Delete(id int64) error {
	query := `UPDATE transactions SET payment_status = 'cancelled' WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	
	return nil
}

func (r *transactionRepository) List(offset, limit int, outletID *int64, transactionType string, search string) ([]models.Transaction, int64, error) {
	whereClause := "WHERE t.payment_status != 'cancelled'"
	args := []interface{}{}
	
	if outletID != nil {
		whereClause += " AND t.outlet_id = ?"
		args = append(args, *outletID)
	}
	
	if transactionType != "" {
		whereClause += " AND t.transaction_type = ?"
		args = append(args, transactionType)
	}
	
	if search != "" {
		whereClause += " AND (t.transaction_number LIKE ? OR c.name LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}
	
	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM transactions t 
		LEFT JOIN customers c ON t.customer_id = c.id
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}
	
	// Get transactions
	query := fmt.Sprintf(`
		SELECT t.id, t.transaction_number, t.transaction_type, t.customer_id, t.outlet_id, 
			   t.user_id, t.service_job_id, t.subtotal_amount, t.discount_amount, 
			   t.tax_amount, t.total_amount, t.payment_status, t.notes, 
			   t.transaction_date, t.created_at, t.updated_at,
			   c.id as "customer.id", c.customer_code as "customer.customer_code", 
			   c.name as "customer.name", c.phone as "customer.phone",
			   o.id as "outlet.id", o.name as "outlet.name",
			   u.id as "user.id", u.full_name as "user.full_name",
			   sj.id as "service_job.id", sj.job_number as "service_job.job_number"
		FROM transactions t
		LEFT JOIN customers c ON t.customer_id = c.id
		LEFT JOIN outlets o ON t.outlet_id = o.id
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN service_jobs sj ON t.service_job_id = sj.id
		%s
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var transactions []models.Transaction
	err = r.db.Select(&transactions, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list transactions: %w", err)
	}
	
	return transactions, total, nil
}

func (r *transactionRepository) GenerateTransactionNumber(transactionType string) (string, error) {
	var prefix string
	switch transactionType {
	case "service":
		prefix = "TXS"
	case "sparepart_sale":
		prefix = "TXP"
	case "vehicle_purchase":
		prefix = "TXV"
	case "vehicle_sale":
		prefix = "TXS"
	default:
		prefix = "TXG"
	}
	
	query := fmt.Sprintf("SELECT COUNT(*) FROM transactions WHERE transaction_number LIKE '%s%%'", prefix)
	
	var count int
	err := r.db.Get(&count, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate transaction number: %w", err)
	}
	
	return fmt.Sprintf("%s%08d", prefix, count+1), nil
}

func (r *transactionRepository) AddDetail(detail *models.TransactionDetail) error {
	query := `
		INSERT INTO transaction_details (transaction_id, product_id, service_id, description, 
										quantity, unit_price, total_price)
		VALUES (:transaction_id, :product_id, :service_id, :description, 
				:quantity, :unit_price, :total_price)
	`
	
	result, err := r.db.NamedExec(query, detail)
	if err != nil {
		return fmt.Errorf("failed to add transaction detail: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get transaction detail ID: %w", err)
	}
	
	detail.ID = id
	return nil
}

func (r *transactionRepository) GetDetails(transactionID int64) ([]models.TransactionDetail, error) {
	query := `
		SELECT td.id, td.transaction_id, td.product_id, td.service_id, td.description,
			   td.quantity, td.unit_price, td.total_price, td.created_at,
			   p.id as "product.id", p.product_code as "product.product_code", 
			   p.name as "product.name",
			   s.id as "service.id", s.service_code as "service.service_code", 
			   s.name as "service.name"
		FROM transaction_details td
		LEFT JOIN products p ON td.product_id = p.id
		LEFT JOIN services s ON td.service_id = s.id
		WHERE td.transaction_id = ?
		ORDER BY td.created_at
	`
	
	var details []models.TransactionDetail
	err := r.db.Select(&details, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction details: %w", err)
	}
	
	return details, nil
}

func (r *transactionRepository) UpdateDetail(id int64, detail *models.TransactionDetail) error {
	query := `
		UPDATE transaction_details 
		SET description = :description, quantity = :quantity, unit_price = :unit_price, 
			total_price = :total_price
		WHERE id = :id
	`
	
	detail.ID = id
	_, err := r.db.NamedExec(query, detail)
	if err != nil {
		return fmt.Errorf("failed to update transaction detail: %w", err)
	}
	
	return nil
}

func (r *transactionRepository) DeleteDetail(id int64) error {
	query := `DELETE FROM transaction_details WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction detail: %w", err)
	}
	
	return nil
}

func (r *transactionRepository) UpdatePaymentStatus(id int64, status string) error {
	query := `UPDATE transactions SET payment_status = ? WHERE id = ?`
	
	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	
	return nil
}

// Payment Repository
type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id int64) (*models.Payment, error)
	GetByTransactionID(transactionID int64) ([]models.Payment, error)
	Delete(id int64) error
	List(offset, limit int, transactionID *int64) ([]models.Payment, int64, error)
	ListPaymentMethods() ([]models.PaymentMethod, error)
	GeneratePaymentNumber() (string, error)
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	query := `
		INSERT INTO payments (payment_number, transaction_id, payment_method_id, amount, 
							  payment_date, reference_number, notes)
		VALUES (:payment_number, :transaction_id, :payment_method_id, :amount, 
				:payment_date, :reference_number, :notes)
	`
	
	result, err := r.db.NamedExec(query, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get payment ID: %w", err)
	}
	
	payment.ID = id
	return nil
}

func (r *paymentRepository) GetByID(id int64) (*models.Payment, error) {
	query := `
		SELECT p.id, p.payment_number, p.transaction_id, p.payment_method_id, p.amount,
			   p.payment_date, p.reference_number, p.notes, p.created_at,
			   pm.id as "payment_method.id", pm.name as "payment_method.name", 
			   pm.type as "payment_method.type"
		FROM payments p
		LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
		WHERE p.id = ?
	`
	
	var payment models.Payment
	err := r.db.Get(&payment, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	
	return &payment, nil
}

func (r *paymentRepository) GetByTransactionID(transactionID int64) ([]models.Payment, error) {
	query := `
		SELECT p.id, p.payment_number, p.transaction_id, p.payment_method_id, p.amount,
			   p.payment_date, p.reference_number, p.notes, p.created_at,
			   pm.id as "payment_method.id", pm.name as "payment_method.name", 
			   pm.type as "payment_method.type"
		FROM payments p
		LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
		WHERE p.transaction_id = ?
		ORDER BY p.payment_date DESC
	`
	
	var payments []models.Payment
	err := r.db.Select(&payments, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by transaction: %w", err)
	}
	
	return payments, nil
}

func (r *paymentRepository) Delete(id int64) error {
	query := `DELETE FROM payments WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	
	return nil
}

func (r *paymentRepository) List(offset, limit int, transactionID *int64) ([]models.Payment, int64, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	
	if transactionID != nil {
		whereClause += " AND p.transaction_id = ?"
		args = append(args, *transactionID)
	}
	
	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM payments p 
		%s
	`, whereClause)
	
	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}
	
	// Get payments
	query := fmt.Sprintf(`
		SELECT p.id, p.payment_number, p.transaction_id, p.payment_method_id, p.amount,
			   p.payment_date, p.reference_number, p.notes, p.created_at,
			   pm.id as "payment_method.id", pm.name as "payment_method.name", 
			   pm.type as "payment_method.type"
		FROM payments p
		LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
		%s
		ORDER BY p.payment_date DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	
	args = append(args, limit, offset)
	
	var payments []models.Payment
	err = r.db.Select(&payments, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	
	return payments, total, nil
}

func (r *paymentRepository) ListPaymentMethods() ([]models.PaymentMethod, error) {
	query := `
		SELECT id, name, type, account_number, bank_name, is_active, created_at, updated_at
		FROM payment_methods 
		WHERE is_active = true 
		ORDER BY name
	`
	
	var paymentMethods []models.PaymentMethod
	err := r.db.Select(&paymentMethods, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}
	
	return paymentMethods, nil
}

func (r *paymentRepository) GeneratePaymentNumber() (string, error) {
	query := `SELECT COUNT(*) FROM payments WHERE payment_number LIKE 'PAY%'`
	
	var count int
	err := r.db.Get(&count, query)
	if err != nil {
		return "", fmt.Errorf("failed to generate payment number: %w", err)
	}
	
	return fmt.Sprintf("PAY%08d", count+1), nil
}
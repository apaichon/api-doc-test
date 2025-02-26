package payment

import (
	"api/internal/db"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// PaymentRepo represents the repository for payment operations
type PaymentRepo struct {
	DB *db.DB
}

// NewPaymentRepo creates a new instance of PaymentRepo
func NewPaymentRepo() *PaymentRepo {
	db := db.NewDB()
	return &PaymentRepo{DB: db}
}

// GetPayments fetches payments with pagination support
func (pr *PaymentRepo) GetPayments(params db.PaginationParams) (*db.PaginationResponse, error) {
	baseQuery := "SELECT * FROM payments"
	countQuery := "SELECT COUNT(*) FROM payments"

	// Get total count
	var total int64
	row, err := pr.DB.QueryRow(countQuery)
	if err != nil {
		return nil, fmt.Errorf("error counting payments: %w", err)
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error scanning total count: %w", err)
	}

	// Build pagination query
	paginationQuery, args, err := db.BuildPaginationQuery(params)
	if err != nil {
		return nil, fmt.Errorf("error building pagination query: %w", err)
	}

	// Execute paginated query
	rows, err := pr.DB.Query(baseQuery+paginationQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %w", err)
	}
	defer rows.Close()

	var payments []*Payment
	var lastID string

	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.PaymentID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentDate,
			&payment.PayTo,
			&payment.Note,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning payment: %w", err)
		}
		payments = append(payments, &payment)
		lastID = fmt.Sprintf("%v", payment.PaymentID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %w", err)
	}

	// Calculate pagination metadata
	hasMore := false
	totalPages := 0
	if params.Limit > 0 {
		hasMore = len(payments) == params.Limit
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}

	response := &db.PaginationResponse{
		Data:       payments,
		Total:      total,
		HasMore:    hasMore,
		TotalPages: totalPages,
	}

	// Add cursor/page specific fields
	switch params.Type {
	case db.CursorPagination:
		if hasMore {
			response.NextCursor = lastID
		}
	case db.PagePagination:
		response.Page = params.Page
	}

	return response, nil
}

// Get PaymentByID retrieves a payment by its ID from the database
func (pr *PaymentRepo) GetPaymentByID(id string) (*Payment, error) {
	var payment Payment
	row, err := pr.DB.QueryRow(`
		SELECT 
			payment_id, id, amount, payment_method, payment_date,
			pay_to, note, status, description, currency,
			created_at, updated_at 
		FROM payments 
		WHERE payment_id = ? OR id = ?`,
		id, id,
	)

	if err != nil {
		return nil, err
	}

	err = row.Scan(
		&payment.PaymentID,
		&payment.ID,
		&payment.Amount,
		&payment.PaymentMethod,
		&payment.PaymentDate,
		&payment.PayTo,
		&payment.Note,
		&payment.Status,
		&payment.Description,
		&payment.Currency,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// Insert Payment inserts a new payment into the database
func (pr *PaymentRepo) InsertPayment(payment *Payment) (string, error) {
	if payment.PaymentID == "" {
		payment.PaymentID = uuid.New().String()
	}
	if payment.ID == "" {
		payment.ID = payment.PaymentID
	}

	_, err := pr.DB.Insert(`
		INSERT INTO payments (
			payment_id, id, amount, payment_method, payment_date, 
			pay_to, note, status, description, currency, 
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		payment.PaymentID,
		payment.ID,
		payment.Amount,
		payment.PaymentMethod,
		payment.PaymentDate,
		payment.PayTo,
		payment.Note,
		payment.Status,
		payment.Description,
		payment.Currency,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("Error inserting payment: %v\n", err)
		return "", err
	}
	return payment.PaymentID, nil
}

// Update Payment updates an existing payment in the database
func (pr *PaymentRepo) UpdatePayment(payment *Payment) (string, error) {
	_, err := pr.DB.Update(`
		UPDATE payments SET 
			amount=?, payment_method=?, payment_date=?, 
			pay_to=?, note=?, status=?, description=?, 
			currency=?, updated_at=? 
		WHERE payment_id=?`,
		payment.Amount,
		payment.PaymentMethod,
		payment.PaymentDate,
		payment.PayTo,
		payment.Note,
		payment.Status,
		payment.Description,
		payment.Currency,
		payment.UpdatedAt,
		payment.PaymentID,
	)
	if err != nil {
		return "", err
	}
	return payment.PaymentID, nil
}

// Delete Payment deletes a payment from the database
func (pr *PaymentRepo) DeletePayment(id string) (string, error) {
	_, err := pr.DB.Delete("DELETE FROM payments WHERE payment_id=?", id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// SearchPaymentWithCursorPagination searches payments with cursor pagination
func (pr *PaymentRepo) SearchPaymentWithCursorPagination(search string, params db.PaginationParams) (*db.PaginationResponse, error) {
	baseQuery := `SELECT * FROM payments 
		WHERE payment_method LIKE ? OR pay_to LIKE ? OR note LIKE ?`
	countQuery := `SELECT COUNT(*) FROM payments 
		WHERE payment_method LIKE ? OR pay_to LIKE ? OR note LIKE ?`

	searchTerm := "%" + search + "%"
	searchArgs := []interface{}{searchTerm, searchTerm, searchTerm}

	// Get total count
	var total int64
	row, err := pr.DB.QueryRow(countQuery, searchArgs...)
	fmt.Printf("countQuery:%v, searchArgs:%v, err:%v", countQuery, searchArgs, err)

	if err != nil {
		return nil, fmt.Errorf("error counting payments: %w", err)
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error scanning total count: %w", err)
	}

	// Build pagination query
	paginationQuery, paginationArgs, err := db.BuildPaginationQuery(params)
	fmt.Printf("paginationQuery:%v, paginationArgs:%v, err:%v", paginationQuery, paginationArgs, err)
	if err != nil {
		return nil, fmt.Errorf("error building pagination: %w", err)
	}

	// Combine search and pagination args
	args := append(searchArgs, paginationArgs...)

	// Execute search with pagination
	rows, err := pr.DB.Query(baseQuery+paginationQuery, args...)
	fmt.Printf("baseQuery:%v, paginationQuery:%v, args:%v, err:%v", baseQuery, paginationQuery, args, err)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %w", err)
	}
	defer rows.Close()

	var payments []*Payment
	var lastID string

	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.PaymentID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentDate,
			&payment.PayTo,
			&payment.Note,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning payment: %w", err)
		}
		payments = append(payments, &payment)
		lastID = fmt.Sprintf("%v", payment.PaymentID)
	}

	hasMore := len(payments) == params.Limit
	response := &db.PaginationResponse{
		Data:       payments,
		Total:      total,
		HasMore:    hasMore,
		NextCursor: lastID,
	}

	return response, nil
}

// SearchPaymentWithOffsetPagination searches payments with offset pagination
func (pr *PaymentRepo) SearchPaymentWithOffsetPagination(search string, params db.PaginationParams) (*db.PaginationResponse, error) {
	baseQuery := `SELECT * FROM payments 
		WHERE payment_method LIKE ? OR pay_to LIKE ? OR note LIKE ?`
	countQuery := `SELECT COUNT(*) FROM payments 
		WHERE payment_method LIKE ? OR pay_to LIKE ? OR note LIKE ?`

	searchTerm := "%" + search + "%"
	searchArgs := []interface{}{searchTerm, searchTerm, searchTerm}

	// Get total count
	var total int64
	row, err := pr.DB.QueryRow(countQuery, searchArgs...)
	if err != nil {
		return nil, fmt.Errorf("error counting payments: %w", err)
	}
	err = row.Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("error scanning total count: %w", err)
	}

	// Build pagination query
	paginationQuery, paginationArgs, err := db.BuildPaginationQuery(params)
	if err != nil {
		return nil, fmt.Errorf("error building pagination: %w", err)
	}

	// Combine search and pagination args
	args := append(searchArgs, paginationArgs...)

	// Execute search with pagination
	rows, err := pr.DB.Query(baseQuery+paginationQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying payments: %w", err)
	}
	defer rows.Close()

	var payments []*Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.PaymentID,
			&payment.Amount,
			&payment.PaymentMethod,
			&payment.PaymentDate,
			&payment.PayTo,
			&payment.Note,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning payment: %w", err)
		}
		payments = append(payments, &payment)
	}

	hasMore := params.Offset+len(payments) < int(total)
	response := &db.PaginationResponse{
		Data:    payments,
		Total:   total,
		HasMore: hasMore,
	}

	return response, nil
}

package loan

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Status represents the loan application status
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusReviewing Status = "REVIEWING"
	StatusApproved  Status = "APPROVED"
	StatusRejected  Status = "REJECTED"
	StatusDisbursed Status = "DISBURSED"
	StatusCompleted Status = "COMPLETED"
	StatusDefaulted Status = "DEFAULTED"
)

// PaymentStatus represents the status of a payment period
type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "PENDING"
	PaymentPaid       PaymentStatus = "PAID"
	PaymentOverdue    PaymentStatus = "OVERDUE"
	PaymentIncomplete PaymentStatus = "INCOMPLETE"
)

// Evidence represents supporting documents
type Evidence struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

// LoanApplication represents a loan application
type LoanApplication struct {
	ID              string          `json:"id"`
	ApplicantID     string          `json:"applicant_id"`
	Amount          float64         `json:"amount"`
	Term            int             `json:"term"`
	Purpose         string          `json:"purpose"`
	Status          Status          `json:"status"`
	Evidence        []Evidence      `json:"evidence"`
	CreditScore     int             `json:"credit_score"`
	InterestRate    float64         `json:"interest_rate"`
	AppliedAt       time.Time       `json:"applied_at"`
	LastUpdatedAt   time.Time       `json:"last_updated_at"`
	ApprovedAt      *time.Time      `json:"approved_at"`
	DisbursedAt     *time.Time      `json:"disbursed_at"`
	PaymentSchedule []PaymentPeriod `json:"payment_schedule"`
}

// PaymentPeriod represents a single payment period
type PaymentPeriod struct {
	ID              string        `json:"id"`
	LoanID          string        `json:"loan_id"`
	DueDate         time.Time     `json:"due_date"`
	Amount          float64       `json:"amount"`
	InterestAmount  float64       `json:"interest_amount"`
	PrincipalAmount float64       `json:"principal_amount"`
	PaidAmount      float64       `json:"paid_amount"`
	FineAmount      float64       `json:"fine_amount"`
	Status          PaymentStatus `json:"status"`
	PaidAt          *time.Time    `json:"paid_at"`
}

// LoanService handles loan-related operations
type LoanService interface {
	ApplyForLoan(application *LoanApplication, evidence []Evidence) error
	ReviewApplication(loanID string) (*LoanApplication, error)
	ApproveLoan(loanID string, interestRate float64) (*LoanApplication, error)
	RejectLoan(loanID string, reason string) error
	DisburseLoan(loanID string) error
	GeneratePaymentSchedule(loanID string) error
	ProcessPayment(loanID string, periodID string, amount float64) error
	CheckPaymentStatus(loanID string, periodID string) error
	UpdateCreditScore(loanID string, creditScore int, interestRate float64) error
}

// CreditService handles credit checking
type CreditService interface {
	CheckCredit(applicantID string) (int, error)
	ValidateIncome(evidence []Evidence) (bool, error)
	CalculateRisk(creditScore int, amount float64) (float64, error)
}

// PaymentService handles payment processing
type PaymentService interface {
	TransferFunds(fromAccount, toAccount string, amount float64) error
	ValidatePayment(paymentID string) error
	CalculateFine(dueDate time.Time, amount float64) float64
}

// DocumentService handles document management
type DocumentService interface {
	StoreEvidence(evidence *Evidence) error
	GenerateInvoice(payment *PaymentPeriod) error
	GenerateStatement(loanID string) error
}

// Implementation of LoanService
type loanService struct {
	db              *sql.DB
	creditService   CreditService
	paymentService  PaymentService
	documentService DocumentService
}

func NewLoanService(db *sql.DB, cs CreditService, ps PaymentService, ds DocumentService) LoanService {
	return &loanService{
		db:              db,
		creditService:   cs,
		paymentService:  ps,
		documentService: ds,
	}
}

// ApplyForLoan handles new loan applications
func (s *loanService) ApplyForLoan(application *LoanApplication, evidence []Evidence) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Validate application data
	if application.Amount <= 0 || application.Term <= 0 {
		return errors.New("invalid loan amount or term")
	}

	// Insert loan application
	now := time.Now()
	application.Status = StatusPending
	application.AppliedAt = now
	application.LastUpdatedAt = now

	// fmt.Println("application", application)

	_, err = tx.Exec(`
		INSERT INTO loan_applications (
			id, applicant_id, amount, term, purpose, status,credit_score,interest_rate,
			applied_at, last_updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		application.ID, application.ApplicantID, application.Amount,
		application.Term, application.Purpose, application.Status,
		application.CreditScore, application.InterestRate,
		application.AppliedAt, application.LastUpdatedAt,
	)
	if err != nil {
		return err
	}

	// Store evidence
	for _, ev := range evidence {
		ev.UploadedAt = now
		_, err = tx.Exec(`
			INSERT INTO evidence (
				id, type, description, url, uploaded_at,
				loan_application_id
			) VALUES (?, ?, ?, ?, ?, ?)`,
			ev.ID, ev.Type, ev.Description, ev.URL,
			ev.UploadedAt, application.ID,
		)
		if err != nil {
			return err
		}
		application.Evidence = append(application.Evidence, ev)
	}

	return tx.Commit()
}

// ReviewApplication reviews loan application and checks credit
func (s *loanService) ReviewApplication(loanID string) (*LoanApplication, error) {
	// Fetch application from database
	application := &LoanApplication{}
	err := s.db.QueryRow(`
		SELECT id, applicant_id, amount, term, purpose, status,
			   credit_score, interest_rate, applied_at, last_updated_at,
			   approved_at, disbursed_at
		FROM loan_applications WHERE id = ?`, loanID,
	).Scan(
		&application.ID, &application.ApplicantID, &application.Amount,
		&application.Term, &application.Purpose, &application.Status,
		&application.CreditScore, &application.InterestRate,
		&application.AppliedAt, &application.LastUpdatedAt,
		&application.ApprovedAt, &application.DisbursedAt,
	)
	if err != nil {
		return nil, err
	}

	// Check credit score
	creditScore, err := s.creditService.CheckCredit(application.ApplicantID)
	if err != nil {
		return nil, err
	}
	application.CreditScore = creditScore

	// Update application status
	now := time.Now()
	_, err = s.db.Exec(`
		UPDATE loan_applications 
		SET status = ?, credit_score = ?, last_updated_at = ?
		WHERE id = ?`,
		StatusReviewing, creditScore, now, loanID,
	)
	if err != nil {
		return nil, err
	}

	application.Status = StatusReviewing
	application.LastUpdatedAt = now

	return application, nil
}

// ApproveLoan approves the loan and sets interest rate
func (s *loanService) ApproveLoan(loanID string, interestRate float64) (*LoanApplication, error) {
	// Fetch and validate application
	application := &LoanApplication{} // In real implementation, fetch from database

	s.db.QueryRow(`
		SELECT id, applicant_id, amount, term, purpose, status,
			   credit_score, interest_rate, applied_at, last_updated_at,
			   approved_at, disbursed_at
		FROM loan_applications WHERE id = ?`, loanID,
	).Scan(
		&application.ID, &application.ApplicantID, &application.Amount, &application.Term, &application.Purpose, &application.Status,
		&application.CreditScore, &application.InterestRate, &application.AppliedAt, &application.LastUpdatedAt,
		&application.ApprovedAt, &application.DisbursedAt,
	)

	fmt.Printf("application: %v %v \n", loanID, application)

	if application.Status != StatusReviewing {
		return nil, errors.New("loan application is not in review status")
	}

	now := time.Now()
	application.Status = StatusApproved
	application.InterestRate = interestRate
	application.ApprovedAt = &now
	application.LastUpdatedAt = now

	_, err := s.db.Exec(`
		UPDATE loan_applications 
		SET status = ?, interest_rate = ?, approved_at = ?, last_updated_at = ?
		WHERE id = ?`,
		StatusApproved, interestRate, now, now, loanID,
	)

	return application, err
}

// DisburseLoan handles the money transfer to borrower's account
func (s *loanService) DisburseLoan(loanID string) error {
	// Fetch application
	application := &LoanApplication{} // In real implementation, fetch from database

	if application.Status != StatusApproved {
		return errors.New("loan is not approved")
	}

	// Transfer funds
	if err := s.paymentService.TransferFunds("LOAN_ACCOUNT", "BORROWER_ACCOUNT", application.Amount); err != nil {
		return err
	}

	now := time.Now()
	application.Status = StatusDisbursed
	application.DisbursedAt = &now
	application.LastUpdatedAt = now

	return nil
}

// ProcessPayment handles loan payments
func (s *loanService) ProcessPayment(loanID string, periodID string, amount float64) error {
	// Fetch payment period
	period := &PaymentPeriod{} // In real implementation, fetch from database

	// Validate payment
	if err := s.paymentService.ValidatePayment(periodID); err != nil {
		return err
	}

	// Calculate fine if payment is late
	if time.Now().After(period.DueDate) {
		fine := s.paymentService.CalculateFine(period.DueDate, period.Amount)
		period.FineAmount = fine
		period.Amount += fine
	}

	// Process payment
	if amount >= period.Amount {
		period.Status = PaymentPaid
		now := time.Now()
		period.PaidAt = &now
	} else {
		period.Status = PaymentIncomplete
		period.PaidAmount = amount
	}

	// Generate invoice/statement
	return s.documentService.GenerateInvoice(period)
}

// CheckPaymentStatus verifies payment status and updates accordingly
func (s *loanService) CheckPaymentStatus(loanID string, periodID string) error {
	// Fetch payment period
	period := &PaymentPeriod{} // In real implementation, fetch from database

	if time.Now().After(period.DueDate) && period.Status == PaymentPending {
		period.Status = PaymentOverdue
		fine := s.paymentService.CalculateFine(period.DueDate, period.Amount)
		period.FineAmount = fine
		period.Amount += fine

		// Generate overdue statement
		return s.documentService.GenerateStatement(loanID)
	}

	return nil
}

// Helper function to calculate monthly payment
func calculateMonthlyPayment(principal float64, annualRate float64, terms int) float64 {
	monthlyRate := annualRate / 12 / 100
	numerator := monthlyRate * (1 + monthlyRate) * float64(terms)
	denominator := (1+monthlyRate)*float64(terms) - 1
	return principal * (numerator / denominator)
}

// GeneratePaymentSchedule generates the payment schedule for a loan
func (s *loanService) GeneratePaymentSchedule(loanID string) error {
	// Fetch loan application
	application := &LoanApplication{} // In real implementation, fetch from database

	monthlyPayment := calculateMonthlyPayment(application.Amount, application.InterestRate, application.Term)

	// Generate payment periods
	for i := 0; i < application.Term; i++ {
		dueDate := application.DisbursedAt.AddDate(0, i+1, 0)
		period := PaymentPeriod{
			LoanID:  loanID,
			DueDate: dueDate,
			Amount:  monthlyPayment,
			Status:  PaymentPending,
		}
		application.PaymentSchedule = append(application.PaymentSchedule, period)
	}

	for _, period := range application.PaymentSchedule {
		_, err := s.db.Exec(`
			INSERT INTO loan_payments (
				id, loan_id, amount, due_date, status, payment_date, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			period.ID, loanID, period.Amount, period.DueDate, period.Status, period.PaidAt, time.Now(), time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// RejectLoan rejects a loan application with a reason
func (s *loanService) RejectLoan(loanID string, reason string) error {
	// Fetch application
	application := &LoanApplication{} // In real implementation, fetch from database

	application.Status = StatusRejected
	application.LastUpdatedAt = time.Now()

	return nil
}

// UpdateCreditScore updates the credit score for a loan application
func (s *loanService) UpdateCreditScore(loanID string, creditScore int, interestRate float64) error {
	// Update the credit score in the database
	query := fmt.Sprintf(`UPDATE loan_applications 
		SET credit_score = ? , interest_rate = ?
		WHERE id = '%s'`, loanID)
	// fmt.Println("query", query)
	_, err := s.db.Exec(query, creditScore, interestRate)
	// fmt.Println("err: updateCreditScore", err)
	return err
}

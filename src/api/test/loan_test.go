package test

import (
	"api/internal/loan"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

// Mock implementations
type mockCreditService struct{}
type mockPaymentService struct{}
type mockDocumentService struct{}

func (m *mockCreditService) CheckCredit(applicantID string) (int, error)           { return 750, nil }
func (m *mockCreditService) ValidateIncome(evidence []loan.Evidence) (bool, error) { return true, nil }
func (m *mockCreditService) CalculateRisk(creditScore int, amount float64) (float64, error) {
	return 0.05, nil
}

func (m *mockPaymentService) TransferFunds(from, to string, amount float64) error     { return nil }
func (m *mockPaymentService) ValidatePayment(paymentID string) error                  { return nil }
func (m *mockPaymentService) CalculateFine(dueDate time.Time, amount float64) float64 { return 0.0 }

func (m *mockDocumentService) StoreEvidence(evidence *loan.Evidence) error       { return nil }
func (m *mockDocumentService) GenerateInvoice(payment *loan.PaymentPeriod) error { return nil }
func (m *mockDocumentService) GenerateStatement(loanID string) error             { return nil }

// Add after imports
const schema = `
CREATE TABLE IF NOT EXISTS loan_applications (
    id TEXT PRIMARY KEY,
    applicant_id TEXT NOT NULL,
    amount REAL NOT NULL,
    term INTEGER NOT NULL,
    purpose TEXT,
    status TEXT NOT NULL,
    credit_score INTEGER,
    interest_rate REAL,
    applied_at TIMESTAMP NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    approved_at TIMESTAMP,
    disbursed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS evidence (
    id TEXT PRIMARY KEY,
    loan_application_id TEXT NOT NULL,
    type TEXT NOT NULL,
    description TEXT,
    url TEXT,
    uploaded_at TIMESTAMP NOT NULL,
    FOREIGN KEY (loan_application_id) REFERENCES loan_applications(id)
);

CREATE TABLE IF NOT EXISTS payment_periods (
    id TEXT PRIMARY KEY,
    loan_id TEXT NOT NULL,
    due_date TIMESTAMP NOT NULL,
    amount REAL NOT NULL,
    interest_amount REAL NOT NULL,
    principal_amount REAL NOT NULL,
    paid_amount REAL DEFAULT 0,
    fine_amount REAL DEFAULT 0,
    status TEXT NOT NULL,
    paid_at TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loan_applications(id)
);`

func setupTestService(t *testing.T) loan.LoanService {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	// Initialize schema
	_, err = db.Exec(schema)
	assert.NoError(t, err)

	return loan.NewLoanService(
		db,
		&mockCreditService{},
		&mockPaymentService{},
		&mockDocumentService{},
	)
}

func TestLoanApplication(t *testing.T) {
	service := setupTestService(t)
	now := time.Now()

	tests := []struct {
		name        string
		application *loan.LoanApplication
		evidence    []loan.Evidence
		wantErr     bool
	}{
		{
			name: "Valid application",
			application: &loan.LoanApplication{
				ID:            "LOAN-002",
				ApplicantID:   "APP-002",
				Amount:        10000,
				Term:          12,
				Purpose:       "Home Improvement",
				Status:        loan.StatusPending,
				CreditScore:   750,
				InterestRate:  0.05,
				AppliedAt:     now,
				LastUpdatedAt: now,
				Evidence: []loan.Evidence{
					{
						ID:          "DOC-002",
						Type:        "INCOME_STATEMENT",
						Description: "Monthly Income",
						URL:         "http://example.com/doc-002",
						UploadedAt:  now,
					},
				},
			},
			evidence: []loan.Evidence{
				{
					ID:          "DOC-002",
					Type:        "INCOME_STATEMENT",
					Description: "Monthly Income",
					URL:         "http://example.com/doc-002",
					UploadedAt:  now,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid amount",
			application: &loan.LoanApplication{
				ID:            "LOAN-003",
				ApplicantID:   "APP-003",
				Amount:        -1000,
				Term:          12,
				Purpose:       "Invalid Loan",
				Status:        loan.StatusPending,
				CreditScore:   0,
				InterestRate:  0,
				AppliedAt:     now,
				LastUpdatedAt: now,
			},
			evidence: []loan.Evidence{},
			wantErr:  true,
		},
		{
			name: "Missing evidence",
			application: &loan.LoanApplication{
				ID:            "LOAN-004",
				ApplicantID:   "APP-004",
				Amount:        5000,
				Term:          12,
				Purpose:       "Personal Loan",
				Status:        loan.StatusPending,
				CreditScore:   700,
				InterestRate:  0.06,
				AppliedAt:     now,
				LastUpdatedAt: now,
			},
			evidence: []loan.Evidence{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ApplyForLoan(tt.application, tt.evidence)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, loan.StatusPending, tt.application.Status)
				assert.False(t, tt.application.AppliedAt.IsZero())
			}
		})
	}
}

func TestLoanApproval(t *testing.T) {
	service := setupTestService(t)
	now := time.Now()

	// Setup initial loan for testing
	initialLoan := &loan.LoanApplication{
		ID:            "LOAN-001",
		ApplicantID:   "APP-001",
		Amount:        10000,
		Term:          12,
		Purpose:       "Home Improvement",
		Status:        loan.StatusPending,
		CreditScore:   750,
		InterestRate:  0.05,
		AppliedAt:     now,
		LastUpdatedAt: now,
		Evidence: []loan.Evidence{
			{
				ID:          "DOC-001",
				Type:        "INCOME_STATEMENT",
				Description: "Monthly Income",
				URL:         "http://example.com/doc-001",
				UploadedAt:  now,
			},
		},
	}
	err := service.ApplyForLoan(initialLoan, []loan.Evidence{
		{
			ID:          "DOC-001",
			Type:        "INCOME_STATEMENT",
			Description: "Monthly Income",
		},
	})
	assert.NoError(t, err)

	// Review the loan first
	reviewed, err := service.ReviewApplication("LOAN-001")
	assert.NoError(t, err)
	assert.Equal(t, loan.StatusReviewing, reviewed.Status)

	tests := []struct {
		name         string
		loanID       string
		interestRate float64
		wantErr      bool
	}{
		{
			name:         "Valid approval",
			loanID:       "LOAN-001",
			interestRate: 5.0,
			wantErr:      false,
		},
		{
			name:         "Invalid rate",
			loanID:       "LOAN-002",
			interestRate: -1.0,
			wantErr:      true,
		},
		{
			name:         "Non-existent loan",
			loanID:       "LOAN-999",
			interestRate: 5.0,
			wantErr:      true,
		},
		{
			name:         "Already approved loan",
			loanID:       "LOAN-001",
			interestRate: 6.0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			application, err := service.ApproveLoan(tt.loanID, tt.interestRate)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if err == nil && application != nil {
				assert.Equal(t, loan.StatusApproved, application.Status)
				assert.Equal(t, tt.interestRate, application.InterestRate)
				assert.NotNil(t, application.ApprovedAt)
			}
		})
	}
}

func TestPaymentProcessing(t *testing.T) {
	service := setupTestService(t)
	now := time.Now()

	initialLoan := &loan.LoanApplication{
		ID:            "LOAN-001",
		ApplicantID:   "APP-001",
		Amount:        12000,
		Term:          12,
		Purpose:       "Home Improvement",
		Status:        loan.StatusPending,
		CreditScore:   750,
		InterestRate:  0.05,
		AppliedAt:     now,
		LastUpdatedAt: now,
		PaymentSchedule: []loan.PaymentPeriod{
			{
				ID:              "PAY-001",
				LoanID:          "LOAN-001",
				DueDate:         now.AddDate(0, 1, 0),
				Amount:          1000,
				InterestAmount:  50,
				PrincipalAmount: 950,
				PaidAmount:      0,
				FineAmount:      0,
				Status:          loan.PaymentPending,
			},
		},
	}
	err := service.ApplyForLoan(initialLoan, []loan.Evidence{
		{
			ID:          "DOC-001",
			Type:        "INCOME_STATEMENT",
			Description: "Monthly Income",
		},
	})
	assert.NoError(t, err)
	_, err = service.ReviewApplication("LOAN-001")
	assert.NoError(t, err)
	_, err = service.ApproveLoan("LOAN-001", 0.05)
	assert.NoError(t, err)

	err = service.GeneratePaymentSchedule("LOAN-001")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		loanID   string
		periodID string
		amount   float64
		wantErr  bool
	}{
		{
			name:     "Full payment",
			loanID:   "LOAN-001",
			periodID: "PAY-001",
			amount:   1000.0,
			wantErr:  false,
		},
		{
			name:     "Partial payment",
			loanID:   "LOAN-001",
			periodID: "PAY-002",
			amount:   500.0,
			wantErr:  false,
		},
		{
			name:     "Invalid amount",
			loanID:   "LOAN-001",
			periodID: "PAY-003",
			amount:   -100.0,
			wantErr:  true,
		},
		{
			name:     "Non-existent loan",
			loanID:   "LOAN-999",
			periodID: "PAY-001",
			amount:   1000.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ProcessPayment(tt.loanID, tt.periodID, tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

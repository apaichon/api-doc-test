package loan

import "time"

type paymentService struct{}

func NewPaymentService() PaymentService {
    return &paymentService{}
}

func (s *paymentService) TransferFunds(fromAccount, toAccount string, amount float64) error {
    return nil
}

func (s *paymentService) ValidatePayment(paymentID string) error {
    return nil
}

func (s *paymentService) CalculateFine(dueDate time.Time, amount float64) float64 {
    return 0.0
}

type LoanPayment struct {
    ID          string    `json:"id" db:"id"`
    LoanID      string    `json:"loanId" db:"loan_id"`
    Amount      float64   `json:"amount" db:"amount"`
    DueDate     time.Time `json:"dueDate" db:"due_date"`
    Status      string    `json:"status" db:"status"`
    PaymentDate *time.Time `json:"paymentDate,omitempty" db:"payment_date"`
    CreatedAt   time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

const (
    PaymentStatusPending = "PENDING"
    PaymentStatusPaid    = "PAID"
    PaymentStatusOverdue = "OVERDUE"
)

// PaymentSchedule represents a collection of payments for a loan
type PaymentSchedule struct {
    LoanID      string    `json:"loanId"`
    Payments    []LoanPayment `json:"payments"`
    TotalAmount float64   `json:"totalAmount"`
} 
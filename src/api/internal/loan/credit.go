package loan

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

// creditService implements CreditService interface
type creditService struct {
	db *sql.DB
}

// NewCreditService creates a new credit service instance
func NewCreditService(db *sql.DB) CreditService {
	return &creditService{db: db}
}

// CheckCredit checks applicant's credit score
func (s *creditService) CheckCredit(applicantID string) (int, error) {
	var creditScore int
	err := s.db.QueryRow(`
		SELECT credit_score 
		FROM credit_scores 
		WHERE applicant_id = ?
		ORDER BY checked_at DESC 
		LIMIT 1`,
		applicantID,
	).Scan(&creditScore)

	if err == sql.ErrNoRows {
		// If no credit score found, return default score
		return 650, nil // Default moderate score
	}
	if err != nil {
		return 0, fmt.Errorf("failed to check credit: %v", err)
	}

	return creditScore, nil
}

// ValidateIncome validates income evidence
func (s *creditService) ValidateIncome(evidence []Evidence) (bool, error) {
	// Validate required documents are present
	var hasIncomeStatement, hasBankStatement bool

	for _, ev := range evidence {
		switch ev.Type {
		case "INCOME_STATEMENT":
			hasIncomeStatement = true
		case "BANK_STATEMENT":
			hasBankStatement = true
		}
	}

	if !hasIncomeStatement || !hasBankStatement {
		return false, fmt.Errorf("missing required documents")
	}

	// In real implementation, would validate document contents
	return true, nil
}

// CalculateRisk calculates risk rate based on credit score and loan amount
func (s *creditService) CalculateRisk(creditScore int, amount float64) (float64, error) {
	// Base rate from configuration
	baseRate := 0.05 // 5%

	// Adjust rate based on credit score
	var rateAdjustment float64
	switch {
	case creditScore >= 800: // Excellent
		rateAdjustment = -0.02
	case creditScore >= 700: // Good
		rateAdjustment = -0.01
	case creditScore >= 650: // Fair
		rateAdjustment = 0
	case creditScore >= 600: // Poor
		rateAdjustment = 0.02
	default: // Bad
		rateAdjustment = 0.04
	}

	// Adjust rate based on loan amount
	if amount > 50000 {
		rateAdjustment += 0.01
	}

	finalRate := baseRate + rateAdjustment
	if finalRate < 0.03 { // Minimum rate
		finalRate = 0.03
	}
	if finalRate > 0.15 { // Maximum rate
		finalRate = 0.15
	}

	return finalRate, nil
}

func (s *creditService) ApplyForLoan(application *LoanApplication, evidence []Evidence) error {
	// Store application in database
	application.Status = "PENDING"
	application.ID = uuid.New().String()
	// Add implementation details as needed
	return nil
}

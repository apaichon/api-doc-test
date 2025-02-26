package payment

import "time"

// Payment represents a payment in the system
type Payment struct {
	ID            string    `json:"id" example:"1"`
	PaymentID     string    `json:"payment_id"` // Keep for backward compatibility
	Amount        float64   `json:"amount" example:"99.99"`
	Currency      string    `json:"currency" example:"USD"`
	PaymentMethod string    `json:"payment_method"`
	PaymentDate   string    `json:"payment_date"`
	PayTo         string    `json:"pay_to"`
	Note          string    `json:"note"`
	Status        string    `json:"status" example:"completed"`
	Description   string    `json:"description" example:"Payment for services"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreditCard struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}

type CreditCardPayment struct {
	PaymentID      string     `json:"payment_id"`
	CreditCardInfo CreditCard `json:"credit_card_info"`
	Amount         float64    `json:"amount"`
	PayTo          string     `json:"pay_to"`
	Note           string     `json:"note"`
	Status         string     `json:"status"`
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
}

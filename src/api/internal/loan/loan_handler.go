package loan

import (
	"encoding/json"
	"net/http"
)

// LoanHandler handles HTTP requests for loan operations
type LoanHandler struct {
	service LoanService
}

// NewLoanHandler creates a new LoanHandler
func NewLoanHandler(service LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

// ApplyForLoan handles the loan application request
func (h *LoanHandler) ApplyForLoan(w http.ResponseWriter, r *http.Request) {
	var application LoanApplication
	if err := json.NewDecoder(r.Body).Decode(&application); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ApplyForLoan(&application, application.Evidence); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(application)
}

// ReviewApplication handles the review application request
func (h *LoanHandler) ReviewApplication(w http.ResponseWriter, r *http.Request) {
	loanID := r.URL.Query().Get("loanID")
	application, err := h.service.ReviewApplication(loanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(application)
}

// ApproveLoan handles the loan approval request
func (h *LoanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	var request struct {
		LoanID       string  `json:"loan_id"`
		InterestRate float64 `json:"interest_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	application, err := h.service.ApproveLoan(request.LoanID, request.InterestRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(application)
}

// RejectLoan handles the loan rejection request
func (h *LoanHandler) RejectLoan(w http.ResponseWriter, r *http.Request) {
	loanID := r.URL.Query().Get("loanID")
	var request struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.RejectLoan(loanID, request.Reason); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateCreditScore handles the credit score update request
func (h *LoanHandler) UpdateCreditScore(w http.ResponseWriter, r *http.Request) {
	loanID := r.URL.Query().Get("loanID")
	var request struct {
		CreditScore  int     `json:"credit_score"`
		InterestRate float64 `json:"interest_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateCreditScore(loanID, request.CreditScore, request.InterestRate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RegisterRoutes registers the loan routes with the given HTTP mux
func (h *LoanHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/loans/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.ApplyForLoan(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/loans/review", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ReviewApplication(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/loans/approve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.ApproveLoan(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/loans/reject", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.RejectLoan(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/loans/updateCreditScore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.UpdateCreditScore(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

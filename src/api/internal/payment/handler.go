package payment

import (
	"api/internal/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type PaymentHandler struct {
	repo *PaymentRepo
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		repo: NewPaymentRepo(),
	}
}

func GetRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value("requestContext").(string)
	fmt.Printf("GetRequestID: Got RequestID: %s\n", requestID)
	if !ok {
		cookie, err := r.Cookie("request_id")
		if err == nil {
			requestID = cookie.Value
		} else {
			requestID = ""
		}
	}
	return requestID
}

// CreatePaymentHandler handles payment creation
func (h *PaymentHandler) CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload", GetRequestID(r))
		return
	}

	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	id, err := h.repo.InsertPayment(&payment)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create payment", GetRequestID(r))
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	}, "Payment created successfully", GetRequestID(r))
}

// UpdatePaymentHandler handles payment updates
func (h *PaymentHandler) UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var payment Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload", GetRequestID(r))
		return
	}

	payment.UpdatedAt = time.Now()
	_, err := h.repo.UpdatePayment(&payment)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update payment", GetRequestID(r))
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{
		"message": "Payment updated successfully",
	}, "Payment updated successfully", GetRequestID(r))
}

// DeletePaymentHandler handles payment deletion
func (h *PaymentHandler) DeletePaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "Payment ID is required", GetRequestID(r))
		return
	}

	_, err := h.repo.DeletePayment(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete payment", GetRequestID(r))
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{
		"message": "Payment deleted successfully",
	}, "Payment deleted successfully", GetRequestID(r))
}

// GetPaymentsHandler handles paginated payment retrieval
func (h *PaymentHandler) GetPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	paginationType := r.URL.Query().Get("pagination_type")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	fmt.Printf("paginationType:%v, limit:%v", paginationType, limit)

	var params db.PaginationParams
	switch paginationType {
	case "cursor":
		params = db.NewPaginationParams(db.CursorPagination)
		params.Cursor = r.URL.Query().Get("cursor")
		params.Limit = limit

	case "offset":
		params = db.NewPaginationParams(db.OffsetPagination)
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		params.Offset = offset
		params.Limit = limit

	default: // page pagination is default
		params = db.NewPaginationParams(db.PagePagination)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		params.Page = page
		params.Limit = limit
	}
	params.KeyID = "payment_id"
	params.SortFields = []string{"payment_id"}

	// Handle search if present
	search := r.URL.Query().Get("search")

	fmt.Printf("search:%v", search)
	var result *db.PaginationResponse
	var err error

	if search != "" {
		switch paginationType {
		case "cursor":
			result, err = h.repo.SearchPaymentWithCursorPagination(search, params)
		case "offset":
			result, err = h.repo.SearchPaymentWithOffsetPagination(search, params)
		default:
			result, err = h.repo.GetPayments(params)
		}
	} else {
		result, err = h.repo.GetPayments(params)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch payments", GetRequestID(r))
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Add these package-level handler functions
var paymentHandler = NewPaymentHandler()

func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentHandler.CreatePaymentHandler(w, r)
}

func UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentHandler.UpdatePaymentHandler(w, r)
}

func DeletePaymentHandler(w http.ResponseWriter, r *http.Request) {
	paymentHandler.DeletePaymentHandler(w, r)
}

func GetPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	paymentHandler.GetPaymentsHandler(w, r)
}

func SearchPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	// Reuse GetPaymentsHandler since it already has search functionality
	GetPaymentsHandler(w, r)
}

package payment

import (
	"api/internal/handler"
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, code int, message string, requestID string) {
	resp := handler.NewErrorResponse(
		code,
		http.StatusText(code),
		"PAYMENT_ERROR",
		message,
		requestID,
	)
	writeJSON(w, code, resp)
}

func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func writeSuccess(w http.ResponseWriter, code int, data interface{}, message string, requestID string) {
	resp := handler.NewSuccessResponse(
		code,
		http.StatusText(code),
		data,
		message,
		requestID,
	)
	writeJSON(w, code, resp)
}

package types

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message" example:"operation successful"`
}

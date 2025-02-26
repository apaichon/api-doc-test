package handler

import "time"

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Status     int          `json:"status"`
	StatusText string       `json:"status_text"`
	Error      *ErrorDetail `json:"error"`
}

// ErrorDetail contains the detailed error information
type ErrorDetail struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"request_id"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewErrorResponse creates a new error response with the current timestamp
func NewErrorResponse(status int, statusText string, code, message, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Status:     status,
		StatusText: statusText,
		Error: &ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
	}
}

// WithDetails adds error details to the response
func (e *ErrorResponse) WithDetails(field, reason string) *ErrorResponse {
	if e.Error.Details == nil {
		e.Error.Details = make(map[string]string)
	}
	e.Error.Details[field] = reason
	return e
}

// Response represents the standard success response structure
type Response struct {
	Status     int         `json:"status"`
	StatusText string      `json:"status_text"`
	Data       interface{} `json:"data,omitempty"`
	Meta       *MetaData   `json:"meta,omitempty"`
}

// MetaData contains additional information about the response
type MetaData struct {
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	Page      int       `json:"page,omitempty"`
	PerPage   int       `json:"per_page,omitempty"`
	Total     int       `json:"total,omitempty"`
}

// NewResponse creates a new success response with the current timestamp
func NewResponse(status int, statusText string, data interface{}, requestID string) *Response {
	return &Response{
		Status:     status,
		StatusText: statusText,
		Data:       data,
		Meta: &MetaData{
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
	}
}

// WithPagination adds pagination information to the response metadata
func (r *Response) WithPagination(page, perPage, total int) *Response {
	r.Meta.Page = page
	r.Meta.PerPage = perPage
	r.Meta.Total = total
	return r
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(status int, statusText string, data interface{}, message string, requestID string) *Response {
	return &Response{
		Status:     status,
		StatusText: statusText,
		Data:       data,
		Meta: &MetaData{
			RequestID: requestID,
			Timestamp: time.Now().UTC(),
		},
	}
}

package middleware

import (
	"api/internal/handler"
	"api/internal/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mssola/user_agent"
)

// ResponseRecorder captures the response
type ResponseRecorder struct {
	http.ResponseWriter
	Body   *bytes.Buffer
	Status int
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}

type commonResponse struct {
	Status     int           `json:"status"`
	StatusText string        `json:"status_text"`
	Error      errorResponse `json:"error"`
}

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func ApiLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &ResponseRecorder{
			ResponseWriter: w,
			Body:           &bytes.Buffer{},
		}

		start := time.Now()
		body, err := io.ReadAll(r.Body)

		r.Body = io.NopCloser(bytes.NewBuffer(body))

		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			if err = json.Unmarshal(body, &body); err != nil {

				w.WriteHeader(http.StatusBadRequest)

				errResp := handler.NewErrorResponse(
					http.StatusBadRequest,
					"Bad Request",
					"INVALID_REQUEST",
					"Invalid request body",
					GetRequestID(r),
				)

				json.NewEncoder(w).Encode(errResp)
				return
			}
		}

		// Call the next handler
		next.ServeHTTP(recorder, r)

		writeApiLog(r, recorder, body, start)
	})
}

func writeApiLog(r *http.Request, w http.ResponseWriter, body []byte, start time.Time) {
	logEntry, err := prepareApiLog(r, w, body, start)
	if err != nil {
		log.Printf("Log Error:%v", err)
		return
	}

	fmt.Printf("Audit Log: %+v\n", logEntry)
	apiLog := logger.GetLogInitializer()

	// Start the log writing Go routine
	go apiLog.WriteApiLogToFile(*logEntry)
}

func getRequestData(r *http.Request, body []byte) string {
	if r.Method == http.MethodGet {
		return r.URL.RawQuery
	}

	// For non-GET requests, return the body as string
	if len(body) > 0 {
		return string(body)
	}

	return ""
}

func prepareApiLog(r *http.Request, w http.ResponseWriter, body []byte, start time.Time) (*logger.ApiLog, error) {

	uaString := r.Header.Get("User-Agent")
	token := r.Header.Get("Authorization")
	ip := r.RemoteAddr

	// Parse the User-Agent
	ua := user_agent.New(uaString)
	// Get the browser name and version
	browserName, browserVersion := ua.Browser()
	// Get the operating system name
	osInfo := ua.OS()
	device := ua.Model()
	userId := getUserIdFromJWT(token)
	status := getStatusCode(w)
	fmt.Printf("Status: %+v\n", status)
	level := getLogLevel(status.Status)

	requestData := getRequestData(r, body)
	requestData = maskSensitiveData(requestData)
	errorMsg := getErrorFromResponse(w)

	logData := &logger.ApiLog{
		Level:                level,
		RequestID:            GetRequestID(r),
		Timestamp:            time.Now(),
		Duration:             time.Since(start),
		Method:               r.Method,
		Path:                 r.URL.Path,
		StatusCode:           status.Status,
		StatusText:           status.StatusText,
		RequestBody:          requestData,
		ClientIP:             ip,
		ClientBrowser:        browserName,
		ClientBrowserVersion: browserVersion,
		ClientOS:             osInfo,
		ClientOSVersion:      ua.OSInfo().Version,
		ClientDevice:         device,
		UserID:               userId,
		Error:                errorMsg,
	}
	return logData, nil
}

func maskSensitiveData(data string) string {
	// Define sensitive keywords
	sensitiveKeywords := []string{"password", "id_card", "credit_card", "ssid", "card_number", "expiry_date", "cvv", "phone", "mobile_no"}

	for _, keyword := range sensitiveKeywords {
		// Mask the sensitive data
		data = strings.ReplaceAll(data, keyword, "****")
	}

	return data
}

func getLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "ERROR"
	case statusCode >= 400:
		return "WARN"
	case statusCode >= 300:
		return "INFO"
	default:
		return "INFO"
	}
}

func getStatusCode(w http.ResponseWriter) commonResponse {
	recorder := w.(*ResponseRecorder)
	if recorder == nil {
		return commonResponse{
			Status:     http.StatusInternalServerError,
			StatusText: "",
			Error: errorResponse{
				Error: struct {
					Message string `json:"message"`
				}{},
			},
		}
	}

	// If we have a status code from WriteHeader, use it
	if recorder.Status != 0 {
		return commonResponse{
			Status:     recorder.Status,
			StatusText: http.StatusText(recorder.Status),
		}
	}

	// Get response body
	respBody := recorder.Body.Bytes()
	if len(respBody) == 0 {
		return commonResponse{
			Status:     http.StatusOK, // Default to 200 if no status was set
			StatusText: "",
			Error: errorResponse{
				Error: struct {
					Message string `json:"message"`
				}{},
			},
		}
	}

	// Try to parse error message
	var resp commonResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return commonResponse{
			Status:     recorder.Status,
			StatusText: "",
			Error: errorResponse{
				Error: struct {
					Message string `json:"message"`
				}{},
			},
		}
	}

	return resp
}

func getErrorFromResponse(w http.ResponseWriter) string {
	// Create a response recorder to capture the response
	recorder := w.(*ResponseRecorder)
	if recorder == nil {
		return ""
	}

	// Get response body
	respBody := recorder.Body.Bytes()
	if len(respBody) == 0 {
		return ""
	}

	// Try to parse error message
	var errResp errorResponse
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		return ""
	}

	return errResp.Error.Message
}

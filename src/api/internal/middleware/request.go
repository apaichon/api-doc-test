package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"github.com/google/uuid"
)

// RequestContext wraps request data that can be shared across middleware
type RequestContext struct {
	RequestID string
	// Add other fields you want to share
	mu     sync.RWMutex
	values map[string]interface{}
}

// Store custom request context in request

const requestContextKey = "requestContext"

func NewRequestContext() *RequestContext {
	return &RequestContext{
		values: make(map[string]interface{}),
	}
}

// Safe getter/setter methods
func (rc *RequestContext) Set(key string, value interface{}) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.values[key] = value
}

func (rc *RequestContext) Get(key string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	val, ok := rc.values[key]
	return val, ok
}

// Helper to get RequestContext from http.Request
func GetRequestContext(r *http.Request) *RequestContext {
	rc := r.Context().Value(requestContextKey).(*RequestContext)

	fmt.Printf("GetRequestContext: Got RequestContext: %v\n", rc)

	return rc
}

func GetRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value(requestContextKey).(string)
	if !ok {
		requestID = ""
	}

	return requestID
}

// Middleware to initialize RequestContext
func RequestContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set initial RequestID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), requestContextKey, requestID)

		fmt.Printf("RequestContextMiddleware: Set RequestID: %s\n", requestID)

		http.SetCookie(w, &http.Cookie{
			Name:    "request_id",
			Value:   requestID,
			Expires: time.Now().Add(1 * time.Minute),
		}) 

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

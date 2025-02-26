package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"api/internal/handler"

	"golang.org/x/time/rate"
)

// Define the rate limiter middleware
func RateLimitMiddleware(limit rate.Limit, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(limit, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, ok := r.Context().Value(requestContextKey).(string)
			if !ok {
				requestID = ""
			}
			fmt.Printf("RateLimitMiddleware: Got RequestID: %s\n", requestID)

			/*cookie, err := r.Cookie("request_id")
			if err == nil {
				requestID = cookie.Value
			}

			fmt.Printf("RateLimitMiddleware: Got RequestID from cookie: %s\n", requestID)
			*/

			// Check if the request exceeds the rate limit
			if !limiter.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				errResp := handler.NewErrorResponse(
					http.StatusTooManyRequests,
					"Rate limit exceeded",
					"RATE_LIMIT_EXCEEDED",
					"Too many requests, please try again later",
					requestID,
				)

				json.NewEncoder(w).Encode(errResp)
				return
			}

			// Proceed to the next handler if the rate limit is not exceeded
			next.ServeHTTP(w, r)
		})
	}
}

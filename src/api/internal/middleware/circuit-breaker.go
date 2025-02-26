package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"api/internal/handler"

	"github.com/graphql-go/graphql"
)

// Higher-order function that wraps a resolver with timeout handling
func CircuitBreakerResolver(resolver graphql.FieldResolveFn, timeout time.Duration) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(p.Context, timeout)
		defer cancel()

		// Create a channel to signal completion or timeout
		done := make(chan struct{}, 1)
		defer close(done)

		var result interface{}
		var resolverErr error

		// Execute the resolver within a goroutine
		go func() {
			result, resolverErr = resolver(p)
			done <- struct{}{}
		}()

		// Wait for either completion or timeout
		select {
		case <-done:
			if resolverErr != nil {
				return nil, resolverErr
			}
			return result, nil
		case <-ctx.Done():
			// Timeout occurred
			return nil, errors.New("the request exceeded the timeout limit")
		}
	}
}

// Middleware function signature
type Middleware func(http.Handler) http.Handler

func CircuitBreakerMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestID, ok := r.Context().Value(requestContextKey).(string)
			if !ok {
				requestID = ""
			}
			fmt.Printf("CircuitBreakerMiddleware: Got RequestID: %s\n", requestID)

			// Create a context with timeout while preserving existing values
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Create new request with the timeout context
			r = r.WithContext(ctx)

			// Create a channel to signal completion or timeout
			done := make(chan struct{}, 1)

			// Execute the handler within a goroutine
			go func() {
				next.ServeHTTP(w, r)
				done <- struct{}{}
			}()

			// Wait for either completion or timeout
			select {
			case <-done:
				return
			case <-ctx.Done():
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)

				errResp := handler.NewErrorResponse(
					http.StatusGatewayTimeout,
					"Request Timeout",
					"TIMEOUT_ERROR",
					"The request exceeded the timeout limit",
					requestID,
				)

				json.NewEncoder(w).Encode(errResp)
			}
		})
	}
}

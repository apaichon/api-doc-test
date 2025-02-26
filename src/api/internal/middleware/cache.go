package middleware

import (
	"api/internal/cache"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type CacheConfig struct {
	TTL            time.Duration
	ExcludePaths   []string
	ExcludeMethods []string
}

func NewCacheConfig() *CacheConfig {
	return &CacheConfig{
		TTL:            5 * time.Minute, // Default TTL
		ExcludePaths:   []string{"/api/login", "/api/logout"},
		ExcludeMethods: []string{"POST", "PUT", "DELETE", "PATCH"},
	}
}

func CacheMiddleware(config *CacheConfig) func(http.Handler) http.Handler {
	redisCache, err := cache.GetRedisInstance()
	if err != nil {
		log.Printf("Failed to get Redis instance: %v", err)
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip caching for excluded paths and methods
			if shouldSkipCache(r, config) {
				next.ServeHTTP(w, r)
				return
			}

			// Generate cache key
			cacheKey := generateCacheKey(r)
			fmt.Printf("cacheKey:%v", cacheKey)

			// Try to get from cache
			if cached, err := redisCache.Get(cacheKey); err == nil {
				var cachedResponse CachedResponse
				fmt.Printf("cached data:%v", cached)
				if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
					fmt.Printf("cached hit:%v", cachedResponse)
					writeResponse(w, &cachedResponse)
					return
				} else {
					log.Printf("Cache unmarshal error: %v", err)
				}
			}

			fmt.Printf("cached miss")

			// Create response recorder
			rec := &ResponseRecorder{
				ResponseWriter: w,
				Body:           &bytes.Buffer{},
				Status:         200, // Default status
			}

			// Process request
			next.ServeHTTP(rec, r)

			// If no status was explicitly set, default to 200
			if rec.Status == 0 {
				rec.Status = http.StatusOK
			}

			fmt.Printf("rec.Status:%v", rec.Status)

			// Cache the response
			if rec.Status >= 200 && rec.Status < 300 {
				fmt.Printf("rec.Status:%v", rec.Status)
				cacheResponse := &CachedResponse{
					Status:  rec.Status,
					Headers: rec.Header(),
					Body:    rec.Body.Bytes(),
				}
				jsonData, err := json.Marshal(cacheResponse)
				if err != nil {
					log.Printf("Failed to marshal cache response: %v", err)
					return
				}
				if err := redisCache.Set(cacheKey, string(jsonData)); err != nil {
					log.Printf("Failed to set cache: %v", err)
				}
			} else {
				// For non-cacheable responses, write directly to client
				writeResponse(w, &CachedResponse{
					Status:  rec.Status,
					Headers: rec.Header(),
					Body:    rec.Body.Bytes(),
				})
			}
		})
	}
}

type CachedResponse struct {
	Status  int
	Headers http.Header
	Body    []byte
}

func shouldSkipCache(r *http.Request, config *CacheConfig) bool {
	// Skip excluded paths
	for _, path := range config.ExcludePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			return true
		}
	}

	// Skip excluded methods
	for _, method := range config.ExcludeMethods {
		if r.Method == method {
			return true
		}
	}

	return false
}

func generateCacheKey(r *http.Request) string {
	// Combine method, path, query params, and headers for cache key
	data := fmt.Sprintf("%s:%s:%s:%s",
		r.Method,
		r.URL.Path,
		r.URL.RawQuery,
		getRelevantHeaders(r),
	)

	// Generate SHA-256 hash
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func getRelevantHeaders(r *http.Request) string {
	// Add headers that affect the response (e.g., Accept, Authorization)
	relevantHeaders := []string{
		r.Header.Get("Accept"),
		r.Header.Get("Authorization"),
	}
	return strings.Join(relevantHeaders, ":")
}

func writeResponse(w http.ResponseWriter, resp *CachedResponse) {
	// Copy headers
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
			fmt.Printf("key:%v, value:%v", key, value)
		}
	}
	w.WriteHeader(resp.Status)
	w.Write(resp.Body)
}

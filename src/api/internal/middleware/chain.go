package middleware

import (
    "net/http"
    "sync"
)

func ChainMiddleware(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// SharedContext holds values that can be shared across middleware
type SharedContext struct {
    mu     sync.RWMutex
    values map[string]interface{}
}

// NewSharedContext creates a new SharedContext instance
func NewSharedContext() *SharedContext {
    return &SharedContext{
        values: make(map[string]interface{}),
    }
}

// Set adds or updates a value in the shared context
func (sc *SharedContext) Set(key string, value interface{}) {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.values[key] = value
}

// Get retrieves a value from the shared context
func (sc *SharedContext) Get(key string) (interface{}, bool) {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    val, ok := sc.values[key]
    return val, ok
}

// MiddlewareHandler wraps http.Handler with shared context
type MiddlewareHandler struct {
    handler http.Handler
    ctx     *SharedContext
}

func (mh *MiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    mh.handler.ServeHTTP(w, r)
}

// ChainMiddleware chains multiple middleware with shared context
func ChainMiddleware2(h http.Handler, middleware ...func(http.Handler, *SharedContext) http.Handler) http.Handler {
    ctx := NewSharedContext()
    for _, m := range middleware {
        h = m(h, ctx)
    }
    return &MiddlewareHandler{handler: h, ctx: ctx}
}
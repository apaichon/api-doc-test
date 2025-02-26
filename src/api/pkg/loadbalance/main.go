package loadbalance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Strategy defines load balancing algorithm
type Strategy int

const (
	RoundRobin Strategy = iota
	LeastConnections
)

// Backend represents a server instance
type Backend struct {
	URL          *url.URL
	Alive        bool
	ActiveConns  int32
	ReverseProxy *httputil.ReverseProxy
	mux          sync.RWMutex
}

// LoadBalancer distributes incoming requests
type LoadBalancer struct {
	backends    []*Backend
	current     uint64
	strategy    Strategy
	healthCheck bool
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(basePort, instances int, strategy Strategy) *LoadBalancer {
	var backends []*Backend

	for i := 0; i < instances; i++ {
		port := basePort + i
		urlStr := fmt.Sprintf("http://127.0.0.1:%d", port)
		url, err := url.Parse(urlStr)
		if err != nil {
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		backend := &Backend{
			URL:          url,
			Alive:        true,
			ReverseProxy: proxy,
		}

		// Custom error handler for the proxy
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			backend.SetAlive(false)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}

		backends = append(backends, backend)
	}

	return &LoadBalancer{
		backends: backends,
		strategy: strategy,
	}
}

// SetAlive updates backend health status
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive checks if backend is healthy
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

// NextBackend selects the next backend based on strategy
func (lb *LoadBalancer) NextBackend() *Backend {
	switch lb.strategy {
	case RoundRobin:
		return lb.roundRobin()
	case LeastConnections:
		return lb.leastConnections()
	default:
		return lb.roundRobin()
	}
}

// Round Robin selection
func (lb *LoadBalancer) roundRobin() *Backend {
	next := atomic.AddUint64(&lb.current, uint64(1))
	len := uint64(len(lb.backends))
	idx := next % len

	// Check next len backends starting from idx
	for i := uint64(0); i < len; i++ {
		backend := lb.backends[(idx+i)%len]
		if backend.IsAlive() {
			return backend
		}
	}
	return nil
}

// Least Connections selection
func (lb *LoadBalancer) leastConnections() *Backend {
	var leastConn int32 = -1
	var selectedBackend *Backend

	for _, backend := range lb.backends {
		if !backend.IsAlive() {
			continue
		}
		conns := atomic.LoadInt32(&backend.ActiveConns)
		if leastConn == -1 || conns < leastConn {
			leastConn = conns
			selectedBackend = backend
		}
	}
	return selectedBackend
}

// ServeHTTP implements the http.Handler interface
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.NextBackend()
	if backend == nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	atomic.AddInt32(&backend.ActiveConns, 1)
	backend.ReverseProxy.ServeHTTP(w, r)
	atomic.AddInt32(&backend.ActiveConns, -1)
}

// HealthCheck performs health checks on backends
func (lb *LoadBalancer) HealthCheck() {
	log.Println("Starting health checks...")
	for _, backend := range lb.backends {
		prevStatus := backend.IsAlive()
		alive := isBackendAlive(backend.URL)
		backend.SetAlive(alive)

		if alive != prevStatus {
			if alive {
				log.Printf("Backend %s is now UP", backend.URL)
			} else {
				log.Printf("Backend %s is now DOWN", backend.URL)
			}
		}
	}
	log.Println("Health checks completed")
}

// Helper function to check if a backend is alive
func isBackendAlive(u *url.URL) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  false,
		},
	}

	// Try to connect to the health check endpoint
	healthURL := fmt.Sprintf("%s/api/health", u.String())
	req, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", u.String(), err)
		return false
	}

	// Add headers that might be required by middleware
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Health check failed for %s: %v", u.String(), err)
		return false
	}
	defer resp.Body.Close()

	// Check for specific health check response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check returned status %d for %s", resp.StatusCode, u.String())
		return false
	}

	// Verify response body
	var healthResponse struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		log.Printf("Failed to decode health response from %s: %v", u.String(), err)
		return false
	}

	return healthResponse.Status == "ok"
}

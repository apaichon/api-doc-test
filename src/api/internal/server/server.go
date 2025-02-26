package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"api/config"
	"api/internal/auth"
	"api/internal/loan"
	"api/internal/middleware"

	_ "api/cmd/server/docs" // Import swagger docs

	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents the HTTP server
type Server struct {
	db     *sql.DB
	config *config.Config
	loan   loan.LoanService
	server *http.Server
	router *http.ServeMux
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database
	dbPath := fmt.Sprintf("%s/%s", "../../data", cfg.DBName)
	fmt.Println("dbPath:", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Initialize services
	creditService := loan.NewCreditService(db)
	paymentService := loan.NewPaymentService()
	documentService := loan.NewDocumentService() // You'll need to create this too
	loanService := loan.NewLoanService(
		db,
		creditService,
		paymentService,
		documentService,
	)

	return &Server{
		db:     db,
		config: cfg,
		loan:   loanService,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// TODO: Add routes and handlers
	port := fmt.Sprintf(":%d", s.config.GraphQLPort)
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(port, nil)
}

// runServer starts a single server instance on the specified port
func (s *Server) Run(ctx context.Context, port int) {
	// Create a new server instance
	server := s.createServer(port)

	// Start server
	go func() {
		log.Printf("Server is running at http://localhost:%v\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server on port %d failed: %v\n", port, err)
		}
	}()

	<-ctx.Done()
	s.shutdownServer()
}

// createServer creates and configures a new HTTP server
func (s *Server) createServer(port int) *http.Server {
	mux := http.NewServeMux()

	// Update Swagger configuration
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// @Summary Health check endpoint
	// @Description Get the health status of the API
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// @Summary Create a new role
	// @Description Create a new role in the system
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param role body auth.Role true "Role information"
	// @Success 201 {object} auth.Role
	// @Failure 400 {object} ErrorResponse
	// @Router /role [post]
	mux.HandleFunc("/api/role", auth.CreateRoleHandler)

	// @Summary User login
	// @Description Authenticate a user and get JWT token
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param credentials body auth.LoginRequest true "Login credentials"
	// @Success 200 {object} auth.LoginResponse
	// @Failure 401 {object} ErrorResponse
	// @Router /login [post]
	mux.HandleFunc("/api/login", auth.LoginHandler)

	// Create and register loan handler
	loanHandler := loan.NewLoanHandler(s.loan)
	loanHandler.RegisterRoutes(mux)

	handler := middleware.ChainMiddleware(
		mux,
		middleware.GzipMiddleware,
		// middleware.CacheMiddleware(middleware.NewCacheConfig()),
		middleware.ApiLogMiddleware,
		middleware.TracingMiddleware,
		middleware.JWTMiddleware([]string{"/swagger", "/api/health", "/api/login", "/api/register", "/api/logout"}),
		middleware.CircuitBreakerMiddleware(10*time.Second),
		middleware.RateLimitMiddleware(1, 10),
		middleware.RequestContextMiddleware,
		middleware.CorsMiddleware,
	)

	return &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// shutdownServer gracefully shuts down a server
func (s *Server) shutdownServer() {
	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server on %s forced to shutdown: %v\n", s.server.Addr, err)
	} else {
		log.Printf("Server on %s stopped gracefully\n", s.server.Addr)
	}
}

func swaggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/doc.json" {
			w.Header().Set("Content-Type", "application/json")
		} else if r.URL.Path == "/swagger/swagger-ui.css" {
			w.Header().Set("Content-Type", "text/css")
		} else if r.URL.Path == "/swagger/swagger-ui-bundle.js" {
			w.Header().Set("Content-Type", "application/javascript")
		}
		next.ServeHTTP(w, r)
	})
}

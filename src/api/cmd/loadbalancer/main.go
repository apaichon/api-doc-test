package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"api/config"
	"api/pkg/loadbalance"
)

var cfg *config.Config

func init() {
	cfg = config.NewConfig()
}

// main is the entry point for the load balancer
func main() {
	var (
		port      int
		instances int
		basePort  int
	)

	// Parse command line arguments
	flag.IntVar(&port, "port", 3999, "Load balancer port")
	flag.IntVar(&instances, "n", 3, "Number of backend instances")
	flag.IntVar(&basePort, "base", 4000, "Base port for backend servers")
	flag.Parse()

	// Create load balancer
	lb := loadbalance.NewLoadBalancer(basePort, instances, loadbalance.RoundRobin)

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: lb,
	}

	// Start health checker in background
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				lb.HealthCheck()
			case <-ctx.Done():
				return
			}
		}
	}()

	// Start server
	go func() {
		fmt.Printf("Load balancer running on http://localhost:%d\n", port)
		fmt.Printf("Balancing traffic across %d instances (ports %d-%d)\n",
			instances, basePort, basePort+instances-1)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Load balancer failed: %v\n", err)
		}
	}()

	// Wait for interrupt
	<-ctx.Done()
	fmt.Println("\nShutting down load balancer...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Load balancer forced to shutdown: %v\n", err)
	}

	fmt.Println("Load balancer stopped gracefully")
}

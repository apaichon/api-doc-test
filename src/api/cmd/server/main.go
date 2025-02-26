package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/spf13/viper"

	"api/config"
	"api/internal/monitoring"
	"api/internal/server"

	_ "github.com/mattn/go-sqlite3"
)

var cfg *config.Config

func init() {
	// Load configuration
	cfg = config.NewConfig()
}

// @title           API Service
// @version         1.0
// @description     A RESTful API service with authentication, payments, and monitoring
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host           localhost:4000
// @BasePath       /api
// @schemes        http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Parse command line arguments for number of instances
	var instances int
	flag.IntVar(&instances, "n", 1, "Number of server instances to run")
	flag.Parse()

	// If passed as argument without flag, check os.Args
	if len(os.Args) > 1 && instances == 1 {
		if n, err := strconv.Atoi(os.Args[1]); err == nil {
			instances = n
		}
	}

	shutdown, err := monitoring.InitTracer(viper.GetString("TRACE_EXPORTER_URL"))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// WaitGroup to wait for all servers to stop
	var wg sync.WaitGroup

	// Start multiple server instances
	for i := 0; i < instances; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			port := cfg.GraphQLPort + index
			server, err := server.NewServer(cfg)
			if err != nil {
				log.Fatalf("Failed to create server: %v", err)
			}
			server.Run(ctx, port)
		}(i)
	}

	// Wait for interrupt signal
	<-ctx.Done()
	log.Println("\nShutting down servers...")

	// Wait for all servers to stop
	wg.Wait()
	log.Println("All servers stopped gracefully")
}

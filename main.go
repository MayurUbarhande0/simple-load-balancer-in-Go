package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/MayurUbarhande0/reverseproxy/config"
	"github.com/MayurUbarhande0/reverseproxy/helper"
	"github.com/MayurUbarhande0/reverseproxy/middleware"
	"github.com/MayurUbarhande0/reverseproxy/models"
	"github.com/MayurUbarhande0/reverseproxy/routes"
)

func main() {
	// Parse command line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	log.Printf("Starting load balancer with %d backends\n", len(cfg.Backends))

	pool := &models.ServerPool{}

	// Initialize backends
	for _, s := range cfg.Backends {
		serverURL, err := url.Parse(s)
		if err != nil {
			log.Printf("Error parsing %s: %s\n", s, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(serverURL)

		// Customize error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
			log.Printf("Error proxying to %s: %v\n", serverURL, e)
			// Mark backend as down
			for _, backend := range pool.Backends {
				if backend.URL.String() == serverURL.String() {
					backend.SetAlive(false)
					break
				}
			}
			http.Error(w, "Service temporarily unavailable", http.StatusBadGateway)
		}

		backend := &models.Backend{
			URL:          serverURL,
			Alive:        true,
			ReverseProxy: proxy,
		}

		pool.Backends = append(pool.Backends, backend)
		log.Printf("Configured backend: %s\n", serverURL)
	}

	if len(pool.Backends) == 0 {
		log.Fatal("No backends configured")
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start health check routine
	go helper.HealthCheck(ctx, pool, cfg.HealthCheckInterval)

	// Setup HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.LBHandler(pool))
	mux.HandleFunc("/health", routes.HealthHandler(pool))
	mux.HandleFunc("/metrics", routes.MetricsHandler())

	// Wrap with middleware
	handler := middleware.LoggingMiddleware(middleware.MetricsMiddleware(mux))

	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Load Balancer starting on %s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cancel health check context
	cancel()

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*cfg.WriteTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited gracefully")
}


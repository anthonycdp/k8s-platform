package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/platform/api/internal/config"
	"github.com/example/platform/api/internal/handlers"
)

const (
	shutdownTimeout = 30 * time.Second
	serverTimeout   = 15 * time.Second
	idleTimeout     = 60 * time.Second
)

func main() {
	cfg := config.Load()

	handlers.SetConfig(cfg)
	handlers.SetAppInfo(config.AppVersion, cfg.Environment)

	router := buildRouter()
	handler := applyMiddleware(router, cfg)

	httpServer := createServer(cfg.HTTPPort, handler)
	metricsServer := createServer(cfg.MetricsPort, handlers.MetricsHandler())

	startServer(httpServer, cfg.AppName, cfg.HTTPPort, cfg.Environment)
	startServer(metricsServer, "Metrics", cfg.MetricsPort, "")

	gracefulShutdown(httpServer, metricsServer)
}

func buildRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/health", handlers.HealthHandler)
	router.HandleFunc("/ready", handlers.ReadyHandler)
	router.HandleFunc("/api/v1/status", handlers.StatusHandler)
	router.HandleFunc("/api/v1/info", handlers.InfoHandler)
	router.HandleFunc("/api/v1/echo", handlers.EchoHandler)
	return router
}

func applyMiddleware(router *http.ServeMux, cfg *config.Config) http.Handler {
	var handler http.Handler = router
	handler = handlers.MetricsMiddleware(handler)

	if cfg.EnableCORS {
		handler = corsMiddleware(handler)
	}

	return handler
}

func createServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
		IdleTimeout:  idleTimeout,
	}
}

func startServer(server *http.Server, name, port, env string) {
	go func() {
		if env != "" {
			log.Printf("%s starting on :%s (environment: %s)", name, port, env)
		} else {
			log.Printf("%s server listening on :%s", name, port)
		}

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s server error: %v", name, err)
		}
	}()
}

func gracefulShutdown(httpServer, metricsServer *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdownServer(ctx, httpServer, "HTTP")
	shutdownServer(ctx, metricsServer, "Metrics")

	log.Println("Servers stopped")
}

func shutdownServer(ctx context.Context, server *http.Server, name string) {
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("%s server shutdown error: %v", name, err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}


package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jibransaleem/Go_Api_Project/internal/config"
)

func main() {

	// Load configuration
	cfg := config.MustLoad()

	// Create router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the students api"))
	})

	// Create HTTP server
	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
	}

	// Channel to receive OS signals
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(done)

	// Start server in a separate goroutine
	go func() {

		slog.Info("server started", "address", server.Addr)

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("failed to start server: ", err)
		}

	}()

	// Wait for termination signal
	<-done

	// Create context with 5-second timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	slog.Info("shutting down the server")

	// Graceful shutdown
	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shut down server", "error", err)

		// Force close if graceful shutdown fails
		server.Close()
	}

	slog.Info("server stopped")
}

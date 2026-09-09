package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rashad-Muntar/forge/internal/platform/config"
	"github.com/Rashad-Muntar/forge/internal/platform/database"
	"github.com/Rashad-Muntar/forge/internal/platform/httpserver"
	"github.com/Rashad-Muntar/forge/internal/platform/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.AppEnv)

	ctx := context.Background()

	db, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)

			_, _ = w.Write([]byte(`{"status":"not_ready"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})

	server := httpserver.New(cfg.HTTPPort, router)

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(
			"HTTP server started",
			"port", cfg.HTTPPort,
			"environment", cfg.AppEnv,
		)

		if err := server.Start(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdownSignals := make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignals,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		log.Error("HTTP server failed", "error", err)
		os.Exit(1)

	case signal := <-shutdownSignals:
		log.Info("shutdown signal received", "signal", signal)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(); err != nil {
		log.Error("HTTP server shutdown failed", "error", err)
		os.Exit(1)
	}

	<-shutdownCtx.Done()

	log.Info("application stopped")
}

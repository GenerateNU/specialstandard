package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"specialstandard/internal/config"
	"specialstandard/internal/service"
	"syscall"
	"time"

	"github.com/sethvargo/go-envconfig"
)

func main() {
	// Load configuration from environment variables
	var cfg config.Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalln("Error processing .env file: ", err)
	}

	// Initialize application with config
	app := service.InitApp(cfg)

	// Close database connection when main exits
	defer func() {
		slog.Info("Closing database connection")
		if err := app.Repo.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	port := cfg.Application.Port

	// Listen for connections with a goroutine
	go func() {
		if err := app.Server.Listen(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal (SIGINT or SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("Shutting down server")

	// Shutdown server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.ShutdownWithContext(ctx); err != nil {
		slog.Error("failed to shutdown server gracefully", "error", err)
	}

	slog.Info("Server shutdown complete")
}

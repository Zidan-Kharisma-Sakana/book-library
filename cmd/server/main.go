package main

import (
	"context"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/api"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/db"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/config"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/logger"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Initialize(cfg.LogLevel)
	logger.Info("Starting book library API service")

	database, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer func() {
		err := db.Close(database)
		if err != nil {
			logger.Fatal("Failed to close database", "error", err)
		}
	}()

	validate := validator.New()

	// Initialize repositories
	bookRepo := repository.NewBookRepository(database)
	authorRepo := repository.NewAuthorRepository(database)
	userRepo := repository.NewUserRepository(database)

	// Initialize services
	authService := service.NewAuthService(cfg)
	bookService := service.NewBookService(validate, bookRepo, authorRepo)
	authorService := service.NewAuthorService(validate, authorRepo)
	userService := service.NewUserService(cfg, validate, userRepo, *authService)

	srv := api.NewServer(
		cfg.ServerAddress,
		authService,
		bookService,
		authorService,
		userService,
	)

	// Start the server in a goroutine
	go func() {
		logger.Info("Starting server", "address", cfg.ServerAddress)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	// Something clean up logiv

	logger.Info("Server exited properly")
}

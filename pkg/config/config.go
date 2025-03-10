package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	LogLevel      string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080" // Default port
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Default log level
	}

	return &Config{
		ServerAddress: serverAddress,
		DatabaseURL:   databaseURL,
		LogLevel:      logLevel,
	}, nil
}

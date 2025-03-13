package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress  string
	DatabaseURL    string
	LogLevel       string
	MigrationsPath string
	JWTSecret      string
	TokenExpiry    time.Duration
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

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations" // Default migrations path
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	tokenExpiry := os.Getenv("TOKEN_EXPIRY")
	if tokenExpiry == "" {
		return nil, errors.New("TOKEN_EXPIRY environment variable is required")
	}
	tokenExpiryDuration, err := time.ParseDuration(tokenExpiry)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddress:  serverAddress,
		DatabaseURL:    databaseURL,
		LogLevel:       logLevel,
		MigrationsPath: migrationsPath,
		JWTSecret:      jwtSecret,
		TokenExpiry:    tokenExpiryDuration,
	}, nil
}

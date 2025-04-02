package config

import (
	"errors"
	"golang.org/x/time/rate"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type RateLimitConfig struct {
	Enabled bool
	Limit   rate.Limit
	Burst   int
}

type Config struct {
	ServerAddress  string
	DatabaseURL    string
	LogLevel       string
	MigrationsPath string
	JWTSecret      string
	TokenExpiry    time.Duration
	AllowedOrigins []string
	RateLimit      *RateLimitConfig
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

	isRateLimitEnabledStr := os.Getenv("RATE_LIMIT_ENABLED")
	isRateLimitEnabled := false

	if isRateLimitEnabledStr == "true" || isRateLimitEnabledStr == "1" || isRateLimitEnabledStr == "yes" {
		isRateLimitEnabled = true
	} else if isRateLimitEnabledStr == "" {
		return nil, errors.New("RATE_LIMIT_ENABLED environment variable is required")
	}

	rateLimitConfig := &RateLimitConfig{
		Enabled: isRateLimitEnabled,
		Limit:   rate.Limit(float64(tokenExpiryDuration.Seconds())),
		Burst:   1,
	}
	
	if isRateLimitEnabled {
		requestPerSecond := os.Getenv("RATE_LIMIT_REQUEST_PER_SECOND")
		if requestPerSecond == "" {
			return nil, errors.New("RATE_LIMIT_REQUEST_PER_SECOND environment variable is required")
		}
		rps, _ := strconv.ParseFloat(requestPerSecond, 64)
		rateLimitConfig.Limit = rate.Limit(rps)

		requestBurst := os.Getenv("RATE_LIMIT_REQUEST_BURST")
		if requestBurst == "" {
			return nil, errors.New("RATE_LIMIT_REQUEST_BURST environment variable is required")
		}
		burst, _ := strconv.Atoi(requestBurst)
		rateLimitConfig.Burst = burst
	}

	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		return nil, errors.New("ALLOWED_ORIGINS environment variable is required")
	}
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	return &Config{
		ServerAddress:  serverAddress,
		DatabaseURL:    databaseURL,
		LogLevel:       logLevel,
		MigrationsPath: migrationsPath,
		JWTSecret:      jwtSecret,
		TokenExpiry:    tokenExpiryDuration,
		AllowedOrigins: allowedOrigins,
		RateLimit:      rateLimitConfig,
	}, nil
}

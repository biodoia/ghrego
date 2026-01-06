// Package config provides application configuration management
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port             string
	DatabaseURL      string
	APIKey           string
	LogLevel         string
	SkipBackendCheck bool

	// CORS
	AllowedOrigins []string

	// Timeouts
	ServerTimeout   time.Duration
	BackendTimeout  time.Duration
	ShutdownTimeout time.Duration

	// Redis (optional)
	RedisEnabled  bool
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RedisPrefix   string

	// Limits
	MaxRequestSize int64
	RateLimitRPS   int
	MaxBulkRepos   int
}

// Load loads configuration from environment and .env file
func Load() *Config {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("No .env file found, using environment variables only")
	}

	return &Config{
		// Server
		Port:             getEnvOrDefault("PORT", "8080"),
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", ""),
		APIKey:           os.Getenv("API_KEY"),
		LogLevel:         getEnvOrDefault("LOG_LEVEL", "info"),
		SkipBackendCheck: getEnvBool("SKIP_BACKEND_CHECK"),

		// CORS
		AllowedOrigins: getEnvSlice("ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://localhost:3000"}),

		// Timeouts
		ServerTimeout:   getEnvDuration("SERVER_TIMEOUT", 30*time.Second),
		BackendTimeout:  getEnvDuration("BACKEND_TIMEOUT", 30*time.Second),
		ShutdownTimeout: getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second),

		// Redis
		RedisEnabled:  os.Getenv("REDIS_ADDR") != "",
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		RedisPrefix:   getEnvOrDefault("REDIS_PREFIX", "gateway:"),

		// Limits
		MaxRequestSize: getEnvInt64("MAX_REQUEST_SIZE", 10*1024*1024),
		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 100),
		MaxBulkRepos:   getEnvInt("MAX_BULK_REPOS", 100),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}
	b, _ := strconv.ParseBool(value)
	return b
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, _ := strconv.Atoi(value)
	return i
}

func getEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, _ := strconv.ParseInt(value, 10, 64)
	return i
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Warn().Str("key", key).Str("value", value).Msg("Invalid duration format")
		return defaultValue
	}
	return d
}

func getEnvSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// Validate checks configuration validity
func (c *Config) Validate() error {
	if c.Port == "" {
		return ErrInvalidConfig("PORT cannot be empty")
	}
	if c.DatabaseURL == "" {
		return ErrInvalidConfig("DATABASE_URL cannot be empty")
	}
	if c.ServerTimeout <= 0 {
		return ErrInvalidConfig("SERVER_TIMEOUT must be positive")
	}
	if c.MaxRequestSize <= 0 {
		return ErrInvalidConfig("MAX_REQUEST_SIZE must be positive")
	}
	return nil
}

// ErrInvalidConfig is returned when configuration is invalid
type ErrInvalidConfig string

func (e ErrInvalidConfig) Error() string {
	return string(e)
}

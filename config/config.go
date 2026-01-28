package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Environment     string
	DatabaseURL     string
	RedisURL        string
	JWTSecret       string
	EncryptionKey   string
	OpenAIKey       string
	MaxUploadSize   int64
	UploadDir       string
	RateLimitReqs   int
	RateLimitWindow int
}

func Load() (*Config, error) {
	// Load .env file in development
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("Error loading .env file: %v", err)
		}
	}
	config := &Config{
		Port:          getEnv("PORT", "8080"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", " "),
		RedisURL:      getEnv("REDIS_URL", " "),
		JWTSecret:     getEnv("JWT_SECRET", " "),
		EncryptionKey: getEnv("ENCRYPTION_KEY", " "),
		OpenAIKey:     getEnv("OPENAI_KEY", " "),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
	}

	//Parse integers
	var err error
	config.MaxUploadSize, err = strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64) //10 MB
	if err != nil {
		return nil, fmt.Errorf("Invalid MAX_UPLOAD_SIZE: %v", err)
	}
	config.RateLimitReqs, err = strconv.Atoi(getEnv("RATE_LIMIT_REQS", "100"))
	if err != nil {
		return nil, fmt.Errorf("Invalid RATE_LIMIT_REQS: %v", err)
	}
	config.RateLimitWindow, err = strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "3600"))
	if err != nil {
		return nil, fmt.Errorf("Invalid RATE_LIMIT_WINDOW: %v", err)
	}

	//Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DatabaseURL is required")
	}
	if c.JWTSecret == " " {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}
	if c.EncryptionKey == " " {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	if len(c.EncryptionKey) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be 32 characters long")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

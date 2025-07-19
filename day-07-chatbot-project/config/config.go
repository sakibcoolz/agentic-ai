package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the chatbot
type Config struct {
	OpenAIAPIKey  string
	Model         string
	MaxTokens     int
	Temperature   float64
	MaxHistory    int
	RetryAttempts int
	RetryDelay    time.Duration
	SaveDirectory string
}

// Load creates a new configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	cfg := &Config{
		OpenAIAPIKey:  getEnvWithDefault("OPENAI_API_KEY", ""),
		Model:         getEnvWithDefault("OPENAI_MODEL", "gpt-3.5-turbo"),
		MaxTokens:     getEnvIntWithDefault("MAX_TOKENS", 150),
		Temperature:   getEnvFloatWithDefault("TEMPERATURE", 0.7),
		MaxHistory:    getEnvIntWithDefault("MAX_HISTORY", 10),
		RetryAttempts: getEnvIntWithDefault("RETRY_ATTEMPTS", 3),
		RetryDelay:    time.Duration(getEnvIntWithDefault("RETRY_DELAY_MS", 1000)) * time.Millisecond,
		SaveDirectory: getEnvWithDefault("SAVE_DIRECTORY", "./data/conversations"),
	}

	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloatWithDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

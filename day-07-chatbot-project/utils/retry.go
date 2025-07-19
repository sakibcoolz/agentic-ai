package utils

import (
	"fmt"
	"time"
)

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
	}
}

// RetryFunc represents a function that can be retried
type RetryFunc func() error

// Retry executes a function with exponential backoff
func Retry(fn RetryFunc, config RetryConfig) error {
	var lastErr error

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if err := fn(); err != nil {
			lastErr = err

			// Don't sleep after the last attempt
			if attempt < config.MaxAttempts-1 {
				delay := calculateDelay(attempt, config)
				time.Sleep(delay)
			}
		} else {
			return nil // Success
		}
	}

	return fmt.Errorf("operation failed after %d attempts, last error: %w", config.MaxAttempts, lastErr)
}

// calculateDelay calculates the delay for a given attempt using exponential backoff
func calculateDelay(attempt int, config RetryConfig) time.Duration {
	delay := time.Duration(float64(config.BaseDelay) * pow(config.Multiplier, float64(attempt)))

	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	return delay
}

// pow is a simple power function for float64
func pow(base, exp float64) float64 {
	if exp == 0 {
		return 1
	}

	result := base
	for i := 1; i < int(exp); i++ {
		result *= base
	}

	return result
}

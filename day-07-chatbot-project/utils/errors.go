package utils

import (
	"fmt"
	"strings"
)

// ChatbotError represents different types of errors that can occur
type ChatbotError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// ErrorType defines different categories of errors
type ErrorType int

const (
	// ErrorTypeAPI indicates an API-related error
	ErrorTypeAPI ErrorType = iota
	// ErrorTypeRateLimit indicates a rate limiting error
	ErrorTypeRateLimit
	// ErrorTypeNetwork indicates a network-related error
	ErrorTypeNetwork
	// ErrorTypeConfig indicates a configuration error
	ErrorTypeConfig
	// ErrorTypeValidation indicates a validation error
	ErrorTypeValidation
	// ErrorTypeInternal indicates an internal application error
	ErrorTypeInternal
)

// Error implements the error interface
func (e *ChatbotError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap returns the underlying cause
func (e *ChatbotError) Unwrap() error {
	return e.Cause
}

// String returns a string representation of the error type
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeAPI:
		return "API_ERROR"
	case ErrorTypeRateLimit:
		return "RATE_LIMIT_ERROR"
	case ErrorTypeNetwork:
		return "NETWORK_ERROR"
	case ErrorTypeConfig:
		return "CONFIG_ERROR"
	case ErrorTypeValidation:
		return "VALIDATION_ERROR"
	case ErrorTypeInternal:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

// NewChatbotError creates a new ChatbotError
func NewChatbotError(errorType ErrorType, message string, cause error) *ChatbotError {
	return &ChatbotError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// IsRetryable determines if an error should be retried
func IsRetryable(err error) bool {
	if chatbotErr, ok := err.(*ChatbotError); ok {
		switch chatbotErr.Type {
		case ErrorTypeRateLimit, ErrorTypeNetwork:
			return true
		case ErrorTypeAPI:
			// Some API errors are retryable (e.g., temporary server errors)
			return strings.Contains(strings.ToLower(chatbotErr.Message), "temporary") ||
				strings.Contains(strings.ToLower(chatbotErr.Message), "timeout") ||
				strings.Contains(strings.ToLower(chatbotErr.Message), "server error")
		default:
			return false
		}
	}

	// For unknown errors, check if they contain retryable keywords
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "network") ||
		strings.Contains(errMsg, "temporary")
}

// WrapError wraps an existing error with additional context
func WrapError(errorType ErrorType, message string, cause error) error {
	return NewChatbotError(errorType, message, cause)
}

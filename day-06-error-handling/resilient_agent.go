package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

// ResilientAgent represents an AI agent with comprehensive error handling
type ResilientAgent struct {
	client         *openai.Client
	config         *ReliabilityConfig
	retryManager   *RetryManager
	circuitBreaker *CircuitBreaker
	rateLimiter    *RateLimiter
	monitor        *Monitor
	faultInjector  *FaultInjector
	mu             sync.RWMutex
}

// ReliabilityConfig contains all reliability settings
type ReliabilityConfig struct {
	Retry          RetryConfig
	CircuitBreaker CircuitBreakerConfig
	RateLimit      RateLimitConfig
	Monitoring     MonitoringConfig
}

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts       int
	BaseDelay         time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
	JitterPercent     int
	RetriableErrors   []string
}

// CircuitBreakerConfig defines circuit breaker behavior
type CircuitBreakerConfig struct {
	FailureThreshold     int
	RecoveryTimeout      time.Duration
	TestRequestRate      float64
	ConsecutiveSuccesses int
}

// RateLimitConfig defines rate limiting behavior
type RateLimitConfig struct {
	RequestsPerMinute int
	BurstSize         int
	AdaptiveEnabled   bool
	QuotaPercentage   float64
}

// MonitoringConfig defines monitoring behavior
type MonitoringConfig struct {
	MetricsEnabled      bool
	HealthChecksEnabled bool
	AlertThreshold      float64
	MetricsRetention    time.Duration
}

// RetryManager handles retry logic with exponential backoff
type RetryManager struct {
	config RetryConfig
	random *rand.Rand
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config          CircuitBreakerConfig
	state           CircuitState
	failureCount    int
	lastFailureTime time.Time
	successCount    int
	mu              sync.RWMutex
}

// CircuitState represents the circuit breaker state
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "CLOSED"
	case CircuitOpen:
		return "OPEN"
	case CircuitHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	config       RateLimitConfig
	tokens       float64
	lastRefill   time.Time
	requestTimes []time.Time
	mu           sync.Mutex
}

// Monitor collects metrics and health information
type Monitor struct {
	config              MonitoringConfig
	totalRequests       int64
	successfulRequests  int64
	failedRequests      int64
	totalRetries        int64
	successfulRetries   int64
	failedRetries       int64
	circuitBreakerTrips int64
	rateLimitedRequests int64
	responseTimes       []time.Duration
	lastAPISuccess      time.Time
	lastAPIFailure      time.Time
	mu                  sync.RWMutex
}

// FaultInjector simulates various failure scenarios
type FaultInjector struct {
	activeFailures map[string]time.Time
	mu             sync.RWMutex
}

// Metrics represents system metrics
type Metrics struct {
	TotalRequests          int64
	SuccessfulRequests     int64
	FailedRequests         int64
	ErrorRate              float64
	TotalRetries           int64
	SuccessfulRetries      int64
	FailedRetries          int64
	RetrySuccessRate       float64
	CircuitBreakerTrips    int64
	CircuitBreakerState    string
	LastCircuitBreakerTrip time.Time
	RateLimitedRequests    int64
	RequestsPerMinute      float64
	QuotaUsage             float64
	AvgResponseTime        time.Duration
	P95ResponseTime        time.Duration
	FastestResponse        time.Duration
	SlowestResponse        time.Duration
}

// HealthStatus represents system health
type HealthStatus struct {
	Overall             bool
	APIConnection       bool
	CircuitBreakerOpen  bool
	RateLimitExceeded   bool
	LastAPISuccess      time.Time
	LastAPIFailure      time.Time
	ConsecutiveFailures int
	AvailableTokens     int
	MemoryUsage         float64
	GoroutineCount      int
}

// DefaultReliabilityConfig returns default reliability settings
func DefaultReliabilityConfig() *ReliabilityConfig {
	return &ReliabilityConfig{
		Retry: RetryConfig{
			MaxAttempts:       3,
			BaseDelay:         100 * time.Millisecond,
			MaxDelay:          30 * time.Second,
			BackoffMultiplier: 2.0,
			JitterPercent:     25,
			RetriableErrors:   []string{"rate_limit", "timeout", "server_error", "network"},
		},
		CircuitBreaker: CircuitBreakerConfig{
			FailureThreshold:     5,
			RecoveryTimeout:      60 * time.Second,
			TestRequestRate:      0.1,
			ConsecutiveSuccesses: 3,
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: 60,
			BurstSize:         10,
			AdaptiveEnabled:   true,
			QuotaPercentage:   80.0,
		},
		Monitoring: MonitoringConfig{
			MetricsEnabled:      true,
			HealthChecksEnabled: true,
			AlertThreshold:      0.05, // 5% error rate
			MetricsRetention:    24 * time.Hour,
		},
	}
}

// NewResilientAgent creates a new resilient AI agent
func NewResilientAgent(apiKey string, config *ReliabilityConfig) (*ResilientAgent, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if config == nil {
		config = DefaultReliabilityConfig()
	}

	client := openai.NewClient(apiKey)

	agent := &ResilientAgent{
		client:         client,
		config:         config,
		retryManager:   NewRetryManager(config.Retry),
		circuitBreaker: NewCircuitBreaker(config.CircuitBreaker),
		rateLimiter:    NewRateLimiter(config.RateLimit),
		monitor:        NewMonitor(config.Monitoring),
		faultInjector:  NewFaultInjector(),
	}

	return agent, nil
}

// NewRetryManager creates a new retry manager
func NewRetryManager(config RetryConfig) *RetryManager {
	return &RetryManager{
		config: config,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  CircuitClosed,
	}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		config:     config,
		tokens:     float64(config.BurstSize),
		lastRefill: time.Now(),
	}
}

// NewMonitor creates a new monitor
func NewMonitor(config MonitoringConfig) *Monitor {
	return &Monitor{
		config:        config,
		responseTimes: make([]time.Duration, 0, 1000),
	}
}

// NewFaultInjector creates a new fault injector
func NewFaultInjector() *FaultInjector {
	return &FaultInjector{
		activeFailures: make(map[string]time.Time),
	}
}

// Chat sends a message and returns a response with full error handling
func (ra *ResilientAgent) Chat(ctx context.Context, message string) (string, error) {
	startTime := time.Now()

	// Check rate limit
	if !ra.rateLimiter.Allow() {
		ra.monitor.RecordRateLimited()
		return "", fmt.Errorf("rate limit exceeded")
	}

	// Check circuit breaker
	if !ra.circuitBreaker.Allow() {
		ra.monitor.RecordFailure(time.Since(startTime))
		return "", fmt.Errorf("circuit breaker is open")
	}

	// Perform the request with retry logic
	response, err := ra.retryManager.Execute(ctx, func() (string, error) {
		return ra.performRequest(ctx, message)
	})

	duration := time.Since(startTime)

	if err != nil {
		ra.circuitBreaker.RecordFailure()
		ra.monitor.RecordFailure(duration)
		return "", err
	}

	ra.circuitBreaker.RecordSuccess()
	ra.monitor.RecordSuccess(duration)
	return response, nil
}

// performRequest makes the actual API request
func (ra *ResilientAgent) performRequest(ctx context.Context, message string) (string, error) {
	// Check for fault injection
	if err := ra.faultInjector.ShouldFail(); err != nil {
		return "", err
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		MaxTokens:   150,
		Temperature: 0.7,
	}

	resp, err := ra.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", ra.classifyError(err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices received")
	}

	return resp.Choices[0].Message.Content, nil
}

// Execute performs an operation with retry logic
func (rm *RetryManager) Execute(ctx context.Context, operation func() (string, error)) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= rm.config.MaxAttempts; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Don't retry if it's the last attempt or error is not retriable
		if attempt == rm.config.MaxAttempts || !rm.isRetriable(err) {
			break
		}

		// Calculate delay with exponential backoff and jitter
		delay := rm.calculateDelay(attempt)

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return "", lastErr
}

// calculateDelay calculates the delay for the next retry attempt
func (rm *RetryManager) calculateDelay(attempt int) time.Duration {
	exponentialDelay := float64(rm.config.BaseDelay) * math.Pow(rm.config.BackoffMultiplier, float64(attempt-1))

	// Apply maximum delay cap
	if exponentialDelay > float64(rm.config.MaxDelay) {
		exponentialDelay = float64(rm.config.MaxDelay)
	}

	// Add jitter to prevent thundering herd
	jitter := 1.0
	if rm.config.JitterPercent > 0 {
		jitterRange := float64(rm.config.JitterPercent) / 100.0
		jitter = 1.0 + (rm.random.Float64()*2-1)*jitterRange
	}

	finalDelay := time.Duration(exponentialDelay * jitter)
	return finalDelay
}

// isRetriable determines if an error should be retried
func (rm *RetryManager) isRetriable(err error) bool {
	errStr := err.Error()
	for _, retriableErr := range rm.config.RetriableErrors {
		if contains(errStr, retriableErr) {
			return true
		}
	}
	return false
}

// Allow checks if a request is allowed through the circuit breaker
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		return cb.shouldAttemptReset()
	case CircuitHalfOpen:
		return cb.shouldAllowTestRequest()
	default:
		return false
	}
}

// shouldAttemptReset checks if the circuit breaker should attempt to reset
func (cb *CircuitBreaker) shouldAttemptReset() bool {
	return time.Since(cb.lastFailureTime) >= cb.config.RecoveryTimeout
}

// shouldAllowTestRequest checks if a test request should be allowed in half-open state
func (cb *CircuitBreaker) shouldAllowTestRequest() bool {
	return rand.Float64() < cb.config.TestRequestRate
}

// RecordFailure records a failure in the circuit breaker
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()
	cb.successCount = 0

	if cb.state == CircuitClosed && cb.failureCount >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
	} else if cb.state == CircuitHalfOpen {
		cb.state = CircuitOpen
	}
}

// RecordSuccess records a success in the circuit breaker
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount = 0
	cb.successCount++

	if cb.state == CircuitOpen {
		cb.state = CircuitHalfOpen
	} else if cb.state == CircuitHalfOpen && cb.successCount >= cb.config.ConsecutiveSuccesses {
		cb.state = CircuitClosed
		cb.successCount = 0
	}
}

// Allow checks if a request is allowed by the rate limiter
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Refill tokens based on time elapsed
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := elapsed.Seconds() * float64(rl.config.RequestsPerMinute) / 60.0

	rl.tokens = math.Min(rl.tokens+tokensToAdd, float64(rl.config.BurstSize))
	rl.lastRefill = now

	// Check if we have tokens available
	if rl.tokens >= 1.0 {
		rl.tokens--

		// Record request time for rate calculation
		rl.requestTimes = append(rl.requestTimes, now)

		// Clean old request times (keep only last minute)
		cutoff := now.Add(-time.Minute)
		for i, reqTime := range rl.requestTimes {
			if reqTime.After(cutoff) {
				rl.requestTimes = rl.requestTimes[i:]
				break
			}
		}

		return true
	}

	return false
}

// classifyError classifies errors for retry and circuit breaker logic
func (ra *ResilientAgent) classifyError(err error) error {
	errStr := err.Error()

	switch {
	case contains(errStr, "rate limit"):
		return fmt.Errorf("rate_limit: %w", err)
	case contains(errStr, "timeout"):
		return fmt.Errorf("timeout: %w", err)
	case contains(errStr, "server error") || contains(errStr, "internal error"):
		return fmt.Errorf("server_error: %w", err)
	case contains(errStr, "network") || contains(errStr, "connection"):
		return fmt.Errorf("network: %w", err)
	default:
		return err
	}
}

// GetMetrics returns current system metrics
func (ra *ResilientAgent) GetMetrics() Metrics {
	return ra.monitor.GetMetrics(ra.circuitBreaker, ra.rateLimiter)
}

// GetHealthStatus returns current health status
func (ra *ResilientAgent) GetHealthStatus() HealthStatus {
	return ra.monitor.GetHealthStatus(ra.circuitBreaker, ra.rateLimiter)
}

// GetConfig returns the current configuration
func (ra *ResilientAgent) GetConfig() *ReliabilityConfig {
	return ra.config
}

// ResetCircuitBreakers resets all circuit breakers
func (ra *ResilientAgent) ResetCircuitBreakers() {
	ra.circuitBreaker.Reset()
}

// ResetMetrics resets all metrics
func (ra *ResilientAgent) ResetMetrics() {
	ra.monitor.Reset()
}

// InjectFault injects a fault for testing
func (ra *ResilientAgent) InjectFault(faultType string, duration time.Duration) {
	ra.faultInjector.InjectFault(faultType, duration)
}

// ClearFaults clears all injected faults
func (ra *ResilientAgent) ClearFaults() {
	ra.faultInjector.ClearFaults()
}

// Helper functions for the circuit breaker and monitor components
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitClosed
	cb.failureCount = 0
	cb.successCount = 0
}

func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (m *Monitor) RecordSuccess(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	m.successfulRequests++
	m.responseTimes = append(m.responseTimes, duration)
	m.lastAPISuccess = time.Now()

	// Keep only recent response times
	if len(m.responseTimes) > 1000 {
		m.responseTimes = m.responseTimes[len(m.responseTimes)-1000:]
	}
}

func (m *Monitor) RecordFailure(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	m.failedRequests++
	m.responseTimes = append(m.responseTimes, duration)
	m.lastAPIFailure = time.Now()
}

func (m *Monitor) RecordRateLimited() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests++
	m.rateLimitedRequests++
}

func (m *Monitor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests = 0
	m.successfulRequests = 0
	m.failedRequests = 0
	m.totalRetries = 0
	m.successfulRetries = 0
	m.failedRetries = 0
	m.circuitBreakerTrips = 0
	m.rateLimitedRequests = 0
	m.responseTimes = m.responseTimes[:0]
}

func (m *Monitor) GetMetrics(cb *CircuitBreaker, rl *RateLimiter) Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics := Metrics{
		TotalRequests:       m.totalRequests,
		SuccessfulRequests:  m.successfulRequests,
		FailedRequests:      m.failedRequests,
		TotalRetries:        m.totalRetries,
		SuccessfulRetries:   m.successfulRetries,
		FailedRetries:       m.failedRetries,
		CircuitBreakerTrips: m.circuitBreakerTrips,
		CircuitBreakerState: cb.GetState().String(),
		RateLimitedRequests: m.rateLimitedRequests,
	}

	if m.totalRequests > 0 {
		metrics.ErrorRate = float64(m.failedRequests) / float64(m.totalRequests)
	}

	if m.totalRetries > 0 {
		metrics.RetrySuccessRate = float64(m.successfulRetries) / float64(m.totalRetries)
	}

	// Calculate response time metrics
	if len(m.responseTimes) > 0 {
		total := time.Duration(0)
		fastest := m.responseTimes[0]
		slowest := m.responseTimes[0]

		for _, rt := range m.responseTimes {
			total += rt
			if rt < fastest {
				fastest = rt
			}
			if rt > slowest {
				slowest = rt
			}
		}

		metrics.AvgResponseTime = total / time.Duration(len(m.responseTimes))
		metrics.FastestResponse = fastest
		metrics.SlowestResponse = slowest

		// Calculate P95
		if len(m.responseTimes) >= 20 {
			sorted := make([]time.Duration, len(m.responseTimes))
			copy(sorted, m.responseTimes)
			// Simple sort for P95 calculation
			for i := 0; i < len(sorted)-1; i++ {
				for j := 0; j < len(sorted)-i-1; j++ {
					if sorted[j] > sorted[j+1] {
						sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
					}
				}
			}
			p95Index := int(float64(len(sorted)) * 0.95)
			metrics.P95ResponseTime = sorted[p95Index]
		}
	}

	// Calculate requests per minute
	rl.mu.Lock()
	metrics.RequestsPerMinute = float64(len(rl.requestTimes))
	rl.mu.Unlock()

	return metrics
}

func (m *Monitor) GetHealthStatus(cb *CircuitBreaker, rl *RateLimiter) HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	cb.mu.RLock()
	circuitOpen := cb.state == CircuitOpen
	consecutiveFailures := cb.failureCount
	cb.mu.RUnlock()

	rl.mu.Lock()
	availableTokens := int(rl.tokens)
	rateLimitExceeded := rl.tokens < 1.0
	rl.mu.Unlock()

	apiConnected := time.Since(m.lastAPISuccess) < time.Since(m.lastAPIFailure)

	overall := apiConnected && !circuitOpen && !rateLimitExceeded

	return HealthStatus{
		Overall:             overall,
		APIConnection:       apiConnected,
		CircuitBreakerOpen:  circuitOpen,
		RateLimitExceeded:   rateLimitExceeded,
		LastAPISuccess:      m.lastAPISuccess,
		LastAPIFailure:      m.lastAPIFailure,
		ConsecutiveFailures: consecutiveFailures,
		AvailableTokens:     availableTokens,
		MemoryUsage:         float64(memStats.Alloc),
		GoroutineCount:      runtime.NumGoroutine(),
	}
}

func (fi *FaultInjector) InjectFault(faultType string, duration time.Duration) {
	fi.mu.Lock()
	defer fi.mu.Unlock()

	fi.activeFailures[faultType] = time.Now().Add(duration)
}

func (fi *FaultInjector) ClearFaults() {
	fi.mu.Lock()
	defer fi.mu.Unlock()

	fi.activeFailures = make(map[string]time.Time)
}

func (fi *FaultInjector) ShouldFail() error {
	fi.mu.RLock()
	defer fi.mu.RUnlock()

	now := time.Now()
	for faultType, expiry := range fi.activeFailures {
		if now.Before(expiry) {
			switch faultType {
			case "timeout":
				return fmt.Errorf("timeout: simulated timeout error")
			case "ratelimit":
				return fmt.Errorf("rate_limit: simulated rate limit error")
			case "server_error":
				return fmt.Errorf("server_error: simulated server error")
			case "network":
				return fmt.Errorf("network: simulated network error")
			case "quota":
				return fmt.Errorf("quota: simulated quota exhaustion")
			}
		}
	}

	// Clean expired faults
	for faultType, expiry := range fi.activeFailures {
		if now.After(expiry) {
			delete(fi.activeFailures, faultType)
		}
	}

	return nil
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOfSubstring(s, substr) != -1)))
}

// Simple substring search
func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

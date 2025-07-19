package main

import (
	"fmt"
	"sync"
	"time"
)

// Lab Exercise 2: Circuit Breaker Pattern
// Learn to implement circuit breakers for fault tolerance

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// Exercise 2.1: Basic Circuit Breaker
type BasicCircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration

	state           CircuitState
	failureCount    int
	lastFailureTime time.Time
	mu              sync.RWMutex
}

func NewBasicCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *BasicCircuitBreaker {
	return &BasicCircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		state:            StateClosed,
	}
}

func (cb *BasicCircuitBreaker) Call(operation func() error) error {
	// TODO: Implement circuit breaker logic
	// 1. Check if circuit is open and if reset timeout has passed
	// 2. If closed, execute operation
	// 3. Handle success/failure and update state
	// 4. Return appropriate error if circuit is open

	fmt.Printf("🔘 Circuit State: %s, Failures: %d/%d\n",
		cb.state, cb.failureCount, cb.failureThreshold)

	return fmt.Errorf("not implemented")
}

func (cb *BasicCircuitBreaker) recordSuccess() {
	// TODO: Reset failure count and ensure circuit is closed
}

func (cb *BasicCircuitBreaker) recordFailure() {
	// TODO: Increment failure count and open circuit if threshold reached
}

func (cb *BasicCircuitBreaker) canAttemptReset() bool {
	// TODO: Check if enough time has passed since last failure
	return false
}

// Exercise 2.2: Advanced Circuit Breaker with Half-Open State
type AdvancedCircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration
	successThreshold int
	halfOpenMaxCalls int

	state           CircuitState
	failureCount    int
	successCount    int
	halfOpenCalls   int
	lastFailureTime time.Time
	mu              sync.RWMutex
}

func NewAdvancedCircuitBreaker(failureThreshold, successThreshold, halfOpenMaxCalls int, resetTimeout time.Duration) *AdvancedCircuitBreaker {
	return &AdvancedCircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		successThreshold: successThreshold,
		halfOpenMaxCalls: halfOpenMaxCalls,
		state:            StateClosed,
	}
}

func (cb *AdvancedCircuitBreaker) Call(operation func() error) error {
	// TODO: Implement advanced circuit breaker with half-open state
	// 1. Handle closed state (normal operation)
	// 2. Handle open state (reject calls, check for reset)
	// 3. Handle half-open state (limited test calls)
	// 4. Transition between states appropriately

	fmt.Printf("🔘 Advanced Circuit - State: %s, Failures: %d, Successes: %d, Half-Open Calls: %d\n",
		cb.state, cb.failureCount, cb.successCount, cb.halfOpenCalls)

	return fmt.Errorf("not implemented")
}

// Exercise 2.3: Circuit Breaker with Metrics
type MetricsCircuitBreaker struct {
	failureThreshold int
	resetTimeout     time.Duration

	state           CircuitState
	failureCount    int
	successCount    int
	totalCalls      int
	rejectedCalls   int
	lastFailureTime time.Time
	lastSuccessTime time.Time
	mu              sync.RWMutex

	// Metrics
	stateChanges []StateChange
	callHistory  []CallResult
}

type StateChange struct {
	From      CircuitState
	To        CircuitState
	Timestamp time.Time
	Reason    string
}

type CallResult struct {
	Success   bool
	Duration  time.Duration
	Timestamp time.Time
	Error     string
}

func NewMetricsCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *MetricsCircuitBreaker {
	return &MetricsCircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		state:            StateClosed,
		stateChanges:     make([]StateChange, 0),
		callHistory:      make([]CallResult, 0),
	}
}

func (cb *MetricsCircuitBreaker) Call(operation func() error) error {
	start := time.Now()

	// TODO: Implement circuit breaker with metrics collection
	// 1. Record call attempt
	// 2. Check circuit state and handle accordingly
	// 3. Execute operation if allowed
	// 4. Record call result with timing
	// 5. Update state and record state changes

	duration := time.Since(start)
	fmt.Printf("📊 Metrics Circuit - State: %s, Total Calls: %d, Rejected: %d, Duration: %v\n",
		cb.state, cb.totalCalls, cb.rejectedCalls, duration)

	return fmt.Errorf("not implemented")
}

func (cb *MetricsCircuitBreaker) GetMetrics() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	successRate := 0.0
	if cb.totalCalls > 0 {
		successRate = float64(cb.successCount) / float64(cb.totalCalls) * 100
	}

	return map[string]interface{}{
		"state":            cb.state.String(),
		"total_calls":      cb.totalCalls,
		"successful_calls": cb.successCount,
		"failed_calls":     cb.failureCount,
		"rejected_calls":   cb.rejectedCalls,
		"success_rate":     successRate,
		"state_changes":    len(cb.stateChanges),
		"last_failure":     cb.lastFailureTime,
		"last_success":     cb.lastSuccessTime,
	}
}

// Mock Operations for Testing

func AlwaysSucceedOperationCB() error {
	time.Sleep(50 * time.Millisecond) // Simulate work
	return nil
}

func AlwaysFailOperationCB() error {
	time.Sleep(50 * time.Millisecond) // Simulate work
	return fmt.Errorf("operation failed")
}

func IntermittentFailOperation(failRate float64) func() error {
	callCount := 0
	return func() error {
		callCount++
		time.Sleep(50 * time.Millisecond)

		// Fail for first few calls, then succeed
		if callCount <= 3 {
			return fmt.Errorf("operation failed (call %d)", callCount)
		}
		return nil
	}
}

func SlowlyRecoveringOperation() func() error {
	callCount := 0
	return func() error {
		callCount++
		time.Sleep(50 * time.Millisecond)

		// Gradually improve success rate
		failureRate := 1.0 - (float64(callCount) * 0.1)
		if failureRate <= 0 {
			return nil
		}

		if time.Now().UnixNano()%100 < int64(failureRate*100) {
			return fmt.Errorf("recovering operation failed (rate: %.1f%%)", failureRate*100)
		}
		return nil
	}
}

// Test Functions

func testBasicCircuitBreaker() {
	fmt.Println("\n🧪 Exercise 2.1: Basic Circuit Breaker")
	fmt.Println("=====================================")

	cb := NewBasicCircuitBreaker(3, 2*time.Second)

	fmt.Println("Testing with failing operation (should open circuit):")
	for i := 0; i < 5; i++ {
		err := cb.Call(AlwaysFailOperationCB)
		fmt.Printf("Call %d: %v\n", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nWaiting for reset timeout...")
	time.Sleep(2500 * time.Millisecond)

	fmt.Println("Testing after timeout (should allow calls):")
	for i := 0; i < 3; i++ {
		err := cb.Call(AlwaysSucceedOperationCB)
		fmt.Printf("Call %d: %v\n", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}
}

func testAdvancedCircuitBreaker() {
	fmt.Println("\n🧪 Exercise 2.2: Advanced Circuit Breaker")
	fmt.Println("=========================================")

	cb := NewAdvancedCircuitBreaker(3, 2, 5, 1*time.Second)

	fmt.Println("Phase 1: Testing with failing operation:")
	for i := 0; i < 5; i++ {
		err := cb.Call(AlwaysFailOperationCB)
		fmt.Printf("Call %d: %v\n", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nPhase 2: Waiting for half-open transition...")
	time.Sleep(1200 * time.Millisecond)

	fmt.Println("Phase 3: Testing with recovering operation:")
	operation := IntermittentFailOperation(0.5)
	for i := 0; i < 8; i++ {
		err := cb.Call(operation)
		fmt.Printf("Call %d: %v\n", i+1, err)
		time.Sleep(200 * time.Millisecond)
	}
}

func testMetricsCircuitBreaker() {
	fmt.Println("\n🧪 Exercise 2.3: Circuit Breaker with Metrics")
	fmt.Println("============================================")

	cb := NewMetricsCircuitBreaker(4, 1*time.Second)

	fmt.Println("Phase 1: Mixed success/failure pattern:")
	operations := []func() error{
		AlwaysSucceedOperationCB,
		AlwaysFailOperationCB,
		AlwaysSucceedOperationCB,
		AlwaysFailOperationCB,
		AlwaysFailOperationCB,
		AlwaysFailOperationCB,
		AlwaysFailOperationCB,    // This should open the circuit
		AlwaysSucceedOperationCB, // This should be rejected
	}

	for i, op := range operations {
		err := cb.Call(op)
		fmt.Printf("Call %d: %v\n", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nCurrent Metrics:")
	metrics := cb.GetMetrics()
	for key, value := range metrics {
		fmt.Printf("  %s: %v\n", key, value)
	}

	fmt.Println("\nPhase 2: Recovery testing:")
	time.Sleep(1200 * time.Millisecond)

	for i := 0; i < 5; i++ {
		err := cb.Call(AlwaysSucceedOperationCB)
		fmt.Printf("Recovery call %d: %v\n", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nFinal Metrics:")
	metrics = cb.GetMetrics()
	for key, value := range metrics {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

// TODO Instructions

func printTodoInstructions() {
	fmt.Println("\n📝 TODO: Complete the Circuit Breaker Implementation")
	fmt.Println("==================================================")
	fmt.Println()
	fmt.Println("Your tasks:")
	fmt.Println("1. ✅ Implement BasicCircuitBreaker.Call() method")
	fmt.Println("2. ✅ Add recordSuccess() and recordFailure() methods")
	fmt.Println("3. ✅ Implement canAttemptReset() logic")
	fmt.Println("4. ✅ Build AdvancedCircuitBreaker with half-open state")
	fmt.Println("5. ✅ Add comprehensive metrics collection")
	fmt.Println()
	fmt.Println("💡 Circuit Breaker States:")
	fmt.Println("• CLOSED: Normal operation, count failures")
	fmt.Println("• OPEN: Reject calls, wait for reset timeout")
	fmt.Println("• HALF_OPEN: Test with limited calls")
	fmt.Println()
	fmt.Println("🎯 Key Concepts:")
	fmt.Println("• Failure threshold triggers state changes")
	fmt.Println("• Reset timeout allows recovery attempts")
	fmt.Println("• Half-open state tests system recovery")
	fmt.Println("• Metrics help optimize thresholds")
	fmt.Println()
	fmt.Println("🛡️ Benefits:")
	fmt.Println("• Prevent cascading failures")
	fmt.Println("• Enable graceful degradation")
	fmt.Println("• Automatic failure detection")
	fmt.Println("• System recovery assistance")
}

func main() {
	fmt.Println("⚡ Lab 2: Circuit Breaker Patterns")
	fmt.Println("=================================")

	printTodoInstructions()

	// Uncomment these as you complete each exercise
	// testBasicCircuitBreaker()
	// testAdvancedCircuitBreaker()
	// testMetricsCircuitBreaker()

	fmt.Println("\n🎉 Lab 2 Complete!")
	fmt.Println("==================")
	fmt.Println("✅ You've implemented circuit breaker patterns")
	fmt.Println("✅ You understand state transitions")
	fmt.Println("✅ You've learned about half-open testing")
	fmt.Println("✅ You can collect reliability metrics")
	fmt.Println()
	fmt.Println("Next: Lab 3 - Rate Limiting Systems")
}

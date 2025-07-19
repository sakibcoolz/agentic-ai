package main

import (
	"fmt"
	"sync"
	"time"
)

// Lab Exercise 3: Rate Limiting Systems
// Learn to implement various rate limiting strategies

// Exercise 3.1: Token Bucket Rate Limiter
type TokenBucket struct {
	capacity   int
	tokens     float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// TODO: Implement token bucket algorithm
	// 1. Calculate time elapsed since last refill
	// 2. Add tokens based on refill rate
	// 3. Cap tokens at capacity
	// 4. Check if tokens available
	// 5. Consume token if available

	fmt.Printf("🪣 Token Bucket - Tokens: %.2f/%d\n", tb.tokens, tb.capacity)
	return false // TODO: Replace with actual implementation
}

func (tb *TokenBucket) refill() {
	// TODO: Implement token refill logic
}

// Exercise 3.2: Sliding Window Rate Limiter
type SlidingWindow struct {
	windowSize  time.Duration
	maxRequests int
	requests    []time.Time
	mu          sync.Mutex
}

func NewSlidingWindow(windowSize time.Duration, maxRequests int) *SlidingWindow {
	return &SlidingWindow{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		requests:    make([]time.Time, 0),
	}
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	// TODO: Implement sliding window algorithm
	// 1. Remove requests outside the window
	// 2. Check if under the request limit
	// 3. Add current request if allowed

	fmt.Printf("📊 Sliding Window - Requests: %d/%d in %v\n",
		len(sw.requests), sw.maxRequests, sw.windowSize)
	return false // TODO: Replace with actual implementation
}

func (sw *SlidingWindow) cleanupOldRequests(now time.Time) {
	// TODO: Remove requests older than window size
}

// Exercise 3.3: Adaptive Rate Limiter
type AdaptiveRateLimiter struct {
	baseRate         float64
	currentRate      float64
	maxRate          float64
	minRate          float64
	successCount     int
	failureCount     int
	adjustmentPeriod time.Duration
	lastAdjustment   time.Time
	tokenBucket      *TokenBucket
	mu               sync.RWMutex
}

func NewAdaptiveRateLimiter(baseRate, minRate, maxRate float64, adjustmentPeriod time.Duration) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		baseRate:         baseRate,
		currentRate:      baseRate,
		maxRate:          maxRate,
		minRate:          minRate,
		adjustmentPeriod: adjustmentPeriod,
		lastAdjustment:   time.Now(),
		tokenBucket:      NewTokenBucket(int(baseRate), baseRate),
	}
}

func (arl *AdaptiveRateLimiter) Allow() bool {
	arl.adjustRateIfNeeded()

	arl.mu.RLock()
	allowed := arl.tokenBucket.Allow()
	arl.mu.RUnlock()

	fmt.Printf("🎛️  Adaptive Rate - Current: %.1f/s, Success: %d, Failures: %d\n",
		arl.currentRate, arl.successCount, arl.failureCount)

	return allowed
}

func (arl *AdaptiveRateLimiter) RecordSuccess() {
	arl.mu.Lock()
	defer arl.mu.Unlock()
	arl.successCount++
}

func (arl *AdaptiveRateLimiter) RecordFailure() {
	arl.mu.Lock()
	defer arl.mu.Unlock()
	arl.failureCount++
}

func (arl *AdaptiveRateLimiter) adjustRateIfNeeded() {
	arl.mu.Lock()
	defer arl.mu.Unlock()

	if time.Since(arl.lastAdjustment) < arl.adjustmentPeriod {
		return
	}

	// TODO: Implement adaptive rate adjustment
	// 1. Calculate success rate
	// 2. Increase rate if high success rate
	// 3. Decrease rate if high failure rate
	// 4. Update token bucket with new rate
	// 5. Reset counters

	arl.lastAdjustment = time.Now()
}

// Exercise 3.4: Distributed Rate Limiter (Simulation)
type DistributedRateLimiter struct {
	instanceID    string
	globalLimit   int
	instanceLimit int
	windowSize    time.Duration

	localRequests  []time.Time
	globalRequests []RequestRecord // Simulated shared state
	mu             sync.RWMutex
}

type RequestRecord struct {
	InstanceID string
	Timestamp  time.Time
}

func NewDistributedRateLimiter(instanceID string, globalLimit, instanceLimit int, windowSize time.Duration) *DistributedRateLimiter {
	return &DistributedRateLimiter{
		instanceID:     instanceID,
		globalLimit:    globalLimit,
		instanceLimit:  instanceLimit,
		windowSize:     windowSize,
		localRequests:  make([]time.Time, 0),
		globalRequests: make([]RequestRecord, 0),
	}
}

func (drl *DistributedRateLimiter) Allow() bool {
	drl.mu.Lock()
	defer drl.mu.Unlock()

	// TODO: Implement distributed rate limiting
	// 1. Check local instance limit
	// 2. Check global limit across all instances
	// 3. Add request to both local and global records
	// 4. Clean up old records

	fmt.Printf("🌐 Distributed Rate - Instance: %d/%d, Global: %d/%d\n",
		len(drl.localRequests), drl.instanceLimit,
		len(drl.globalRequests), drl.globalLimit)

	return false // TODO: Replace with actual implementation
}

// Exercise 3.5: Rate Limiter with Priority Queues
type PriorityRateLimiter struct {
	highPriorityLimiter   *TokenBucket
	normalPriorityLimiter *TokenBucket
	lowPriorityLimiter    *TokenBucket

	highPriorityRatio   float64
	normalPriorityRatio float64
	lowPriorityRatio    float64
}

type Priority int

const (
	LowPriority Priority = iota
	NormalPriority
	HighPriority
)

func NewPriorityRateLimiter(totalRate float64, highRatio, normalRatio, lowRatio float64) *PriorityRateLimiter {
	return &PriorityRateLimiter{
		highPriorityLimiter:   NewTokenBucket(int(totalRate*highRatio), totalRate*highRatio),
		normalPriorityLimiter: NewTokenBucket(int(totalRate*normalRatio), totalRate*normalRatio),
		lowPriorityLimiter:    NewTokenBucket(int(totalRate*lowRatio), totalRate*lowRatio),
		highPriorityRatio:     highRatio,
		normalPriorityRatio:   normalRatio,
		lowPriorityRatio:      lowRatio,
	}
}

func (prl *PriorityRateLimiter) Allow(priority Priority) bool {
	// TODO: Implement priority-based rate limiting
	// 1. Route request to appropriate limiter based on priority
	// 2. Allow high priority requests to borrow from lower priorities
	// 3. Implement spillover logic for unused capacity

	switch priority {
	case HighPriority:
		fmt.Printf("🔴 High Priority Request\n")
		return false // TODO: Replace with actual implementation
	case NormalPriority:
		fmt.Printf("🟡 Normal Priority Request\n")
		return false // TODO: Replace with actual implementation
	case LowPriority:
		fmt.Printf("🟢 Low Priority Request\n")
		return false // TODO: Replace with actual implementation
	default:
		return false
	}
}

// Test Functions

func testTokenBucket() {
	fmt.Println("\n🧪 Exercise 3.1: Token Bucket Rate Limiter")
	fmt.Println("==========================================")

	// Create bucket: 5 tokens, refill 2 tokens/second
	bucket := NewTokenBucket(5, 2.0)

	fmt.Println("Testing burst capacity:")
	for i := 0; i < 8; i++ {
		allowed := bucket.Allow()
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\nWaiting for refill...")
	time.Sleep(2 * time.Second)

	fmt.Println("Testing after refill:")
	for i := 0; i < 4; i++ {
		allowed := bucket.Allow()
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(100 * time.Millisecond)
	}
}

func testSlidingWindow() {
	fmt.Println("\n🧪 Exercise 3.2: Sliding Window Rate Limiter")
	fmt.Println("============================================")

	// Create window: 3 requests per 2 seconds
	window := NewSlidingWindow(2*time.Second, 3)

	fmt.Println("Testing within window limit:")
	for i := 0; i < 6; i++ {
		allowed := window.Allow()
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nWaiting for window to slide...")
	time.Sleep(2500 * time.Millisecond)

	fmt.Println("Testing with fresh window:")
	for i := 0; i < 4; i++ {
		allowed := window.Allow()
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(200 * time.Millisecond)
	}
}

func testAdaptiveRateLimiter() {
	fmt.Println("\n🧪 Exercise 3.3: Adaptive Rate Limiter")
	fmt.Println("======================================")

	// Create adaptive limiter: base 2/s, min 1/s, max 5/s, adjust every 3s
	limiter := NewAdaptiveRateLimiter(2.0, 1.0, 5.0, 3*time.Second)

	fmt.Println("Phase 1: High success rate (should increase rate):")
	for i := 0; i < 10; i++ {
		allowed := limiter.Allow()
		if allowed {
			limiter.RecordSuccess() // High success rate
		}
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(400 * time.Millisecond)
	}

	fmt.Println("\nPhase 2: High failure rate (should decrease rate):")
	for i := 0; i < 10; i++ {
		allowed := limiter.Allow()
		if allowed {
			limiter.RecordFailure() // High failure rate
		}
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(400 * time.Millisecond)
	}
}

func testDistributedRateLimiter() {
	fmt.Println("\n🧪 Exercise 3.4: Distributed Rate Limiter")
	fmt.Println("=========================================")

	// Simulate multiple instances
	instance1 := NewDistributedRateLimiter("instance-1", 10, 4, 5*time.Second)
	instance2 := NewDistributedRateLimiter("instance-2", 10, 4, 5*time.Second)

	fmt.Println("Testing multiple instances sharing global limit:")

	for i := 0; i < 6; i++ {
		// Alternate between instances
		var allowed bool
		var instanceName string

		if i%2 == 0 {
			allowed = instance1.Allow()
			instanceName = "Instance-1"
		} else {
			allowed = instance2.Allow()
			instanceName = "Instance-2"
		}

		fmt.Printf("%s Request %d: %v\n", instanceName, i+1, allowed)
		time.Sleep(200 * time.Millisecond)
	}
}

func testPriorityRateLimiter() {
	fmt.Println("\n🧪 Exercise 3.5: Priority Rate Limiter")
	fmt.Println("======================================")

	// Create priority limiter: 5/s total (50% high, 30% normal, 20% low)
	limiter := NewPriorityRateLimiter(5.0, 0.5, 0.3, 0.2)

	requests := []Priority{
		HighPriority, NormalPriority, LowPriority,
		HighPriority, HighPriority, NormalPriority,
		LowPriority, HighPriority, NormalPriority,
		LowPriority, HighPriority, LowPriority,
	}

	fmt.Println("Testing priority-based rate limiting:")
	for i, priority := range requests {
		allowed := limiter.Allow(priority)
		fmt.Printf("Request %d: %v\n", i+1, allowed)
		time.Sleep(150 * time.Millisecond)
	}
}

// TODO Instructions

func printTodoInstructions() {
	fmt.Println("\n📝 TODO: Complete the Rate Limiting Implementation")
	fmt.Println("=================================================")
	fmt.Println()
	fmt.Println("Your tasks:")
	fmt.Println("1. ✅ Implement TokenBucket.Allow() method")
	fmt.Println("2. ✅ Add SlidingWindow request tracking")
	fmt.Println("3. ✅ Build adaptive rate adjustment logic")
	fmt.Println("4. ✅ Create distributed rate limiting")
	fmt.Println("5. ✅ Implement priority-based routing")
	fmt.Println()
	fmt.Println("💡 Rate Limiting Concepts:")
	fmt.Println("• Token Bucket: Allows bursts up to capacity")
	fmt.Println("• Sliding Window: Fixed number of requests per time window")
	fmt.Println("• Adaptive: Adjusts rate based on success/failure")
	fmt.Println("• Distributed: Coordinates across multiple instances")
	fmt.Println("• Priority: Different limits for different request types")
	fmt.Println()
	fmt.Println("🎯 Key Algorithms:")
	fmt.Println("• Token refill: tokens += (elapsed * rate)")
	fmt.Println("• Window cleanup: remove requests older than window")
	fmt.Println("• Rate adjustment: increase on success, decrease on failure")
	fmt.Println("• Priority spillover: high priority can use lower priority tokens")
	fmt.Println()
	fmt.Println("🚦 Benefits:")
	fmt.Println("• API quota management")
	fmt.Println("• System overload protection")
	fmt.Println("• Fair resource allocation")
	fmt.Println("• Cost optimization")
}

func main() {
	fmt.Println("🚦 Lab 3: Rate Limiting Systems")
	fmt.Println("===============================")

	printTodoInstructions()

	// Uncomment these as you complete each exercise
	// testTokenBucket()
	// testSlidingWindow()
	// testAdaptiveRateLimiter()
	// testDistributedRateLimiter()
	// testPriorityRateLimiter()

	fmt.Println("\n🎉 Lab 3 Complete!")
	fmt.Println("==================")
	fmt.Println("✅ You've implemented rate limiting algorithms")
	fmt.Println("✅ You understand burst vs sustained rates")
	fmt.Println("✅ You've learned adaptive rate control")
	fmt.Println("✅ You can handle distributed coordination")
	fmt.Println("✅ You've built priority-based systems")
	fmt.Println()
	fmt.Println("Next: Lab 4 - End-to-End Integration")
}

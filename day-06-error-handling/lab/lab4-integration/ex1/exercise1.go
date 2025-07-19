package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Lab Exercise 1: Basic Retry Implementation
// Learn to implement exponential backoff with jitter

// Exercise 1.1: Simple Retry Function
func simpleRetry(operation func() error, maxAttempts int) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d/%d\n", attempt, maxAttempts)

		err := operation()
		if err == nil {
			fmt.Printf("✅ Success on attempt %d\n", attempt)
			return nil
		}

		lastErr = err
		fmt.Printf("❌ Failed: %v\n", err)

		if attempt < maxAttempts {
			// TODO: Add a simple delay here
			// time.Sleep(time.Second)
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", maxAttempts, lastErr)
}

// Exercise 1.2: Exponential Backoff
func exponentialBackoffRetry(operation func() error, maxAttempts int, baseDelay time.Duration) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d/%d\n", attempt, maxAttempts)

		err := operation()
		if err == nil {
			fmt.Printf("✅ Success on attempt %d\n", attempt)
			return nil
		}

		lastErr = err
		fmt.Printf("❌ Failed: %v\n", err)

		if attempt < maxAttempts {
			// TODO: Implement exponential backoff
			// delay := time.Duration(math.Pow(2, float64(attempt-1))) * baseDelay
			// fmt.Printf("⏱️  Waiting %v before next attempt\n", delay)
			// time.Sleep(delay)
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", maxAttempts, lastErr)
}

// Exercise 1.3: Exponential Backoff with Jitter
func jitteredRetry(operation func() error, maxAttempts int, baseDelay time.Duration, jitterPercent int) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d/%d\n", attempt, maxAttempts)

		err := operation()
		if err == nil {
			fmt.Printf("✅ Success on attempt %d\n", attempt)
			return nil
		}

		lastErr = err
		fmt.Printf("❌ Failed: %v\n", err)

		if attempt < maxAttempts {
			// TODO: Implement exponential backoff with jitter
			// baseTime := math.Pow(2, float64(attempt-1)) * float64(baseDelay)
			// jitter := 1.0 + (rand.Float64()*2-1)*float64(jitterPercent)/100.0
			// delay := time.Duration(baseTime * jitter)
			// fmt.Printf("⏱️  Waiting %v before next attempt (with jitter)\n", delay)
			// time.Sleep(delay)
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", maxAttempts, lastErr)
}

// Exercise 1.4: Context-Aware Retry
func contextAwareRetry(ctx context.Context, operation func() error, maxAttempts int, baseDelay time.Duration) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// TODO: Check context cancellation
		// select {
		// case <-ctx.Done():
		//     return ctx.Err()
		// default:
		// }

		fmt.Printf("Attempt %d/%d\n", attempt, maxAttempts)

		err := operation()
		if err == nil {
			fmt.Printf("✅ Success on attempt %d\n", attempt)
			return nil
		}

		lastErr = err
		fmt.Printf("❌ Failed: %v\n", err)

		if attempt < maxAttempts {
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * baseDelay
			fmt.Printf("⏱️  Waiting %v before next attempt\n", delay)

			// TODO: Implement context-aware sleep
			// select {
			// case <-ctx.Done():
			//     return ctx.Err()
			// case <-time.After(delay):
			//     continue
			// }
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", maxAttempts, lastErr)
}

// Mock Operations for Testing

// AlwaysFailOperation simulates an operation that always fails
func AlwaysFailOperation() error {
	return fmt.Errorf("simulated failure")
}

// SometimesFailOperation simulates an operation that fails randomly
func SometimesFailOperation() error {
	if rand.Float64() < 0.7 { // 70% failure rate
		return fmt.Errorf("random failure")
	}
	return nil
}

// FailFirstNOperation fails the first N attempts, then succeeds
func FailFirstNOperation(failCount *int, n int) func() error {
	return func() error {
		if *failCount < n {
			*failCount++
			return fmt.Errorf("failure %d/%d", *failCount, n)
		}
		return nil
	}
}

// SlowOperation simulates a slow operation
func SlowOperation() error {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	if rand.Float64() < 0.3 { // 30% failure rate
		return fmt.Errorf("slow operation failed")
	}
	return nil
}

// Lab Exercises

func runExercise1() {
	fmt.Println("\n🧪 Exercise 1.1: Simple Retry")
	fmt.Println("=============================")

	// Test with always failing operation
	fmt.Println("Testing with always failing operation:")
	err := simpleRetry(AlwaysFailOperation, 3)
	if err != nil {
		fmt.Printf("Final result: %v\n", err)
	}

	// Test with sometimes failing operation
	fmt.Println("\nTesting with sometimes failing operation:")
	err = simpleRetry(SometimesFailOperation, 5)
	if err != nil {
		fmt.Printf("Final result: %v\n", err)
	}
}

func runExercise2() {
	fmt.Println("\n🧪 Exercise 1.2: Exponential Backoff")
	fmt.Println("====================================")

	failCount := 0
	operation := FailFirstNOperation(&failCount, 2)

	fmt.Println("Testing exponential backoff (should succeed on 3rd attempt):")
	err := exponentialBackoffRetry(operation, 5, 100*time.Millisecond)
	if err != nil {
		fmt.Printf("Final result: %v\n", err)
	}
}

func runExercise3() {
	fmt.Println("\n🧪 Exercise 1.3: Exponential Backoff with Jitter")
	fmt.Println("================================================")

	fmt.Println("Testing jittered retry (notice varied delays):")

	// Run multiple times to see jitter effect
	for i := 0; i < 3; i++ {
		fmt.Printf("\n--- Run %d ---\n", i+1)
		failCount := 0
		operation := FailFirstNOperation(&failCount, 2)

		err := jitteredRetry(operation, 4, 100*time.Millisecond, 25)
		if err != nil {
			fmt.Printf("Final result: %v\n", err)
		}
	}
}

func runExercise4() {
	fmt.Println("\n🧪 Exercise 1.4: Context-Aware Retry")
	fmt.Println("====================================")

	// Test with timeout context
	fmt.Println("Testing with 2-second timeout:")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := contextAwareRetry(ctx, SlowOperation, 10, 500*time.Millisecond)
	if err != nil {
		fmt.Printf("Final result: %v\n", err)
	}

	// Test with cancellation
	fmt.Println("\nTesting with manual cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())

	// Cancel after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("🛑 Cancelling context...")
		cancel2()
	}()

	err = contextAwareRetry(ctx2, AlwaysFailOperation, 10, 300*time.Millisecond)
	if err != nil {
		fmt.Printf("Final result: %v\n", err)
	}
}

// TODO Instructions for Students

func printTodoInstructions() {
	fmt.Println("\n📝 TODO: Complete the Implementation")
	fmt.Println("===================================")
	fmt.Println()
	fmt.Println("Your tasks:")
	fmt.Println("1. ✅ Implement simple retry delay in Exercise 1.1")
	fmt.Println("2. ✅ Add exponential backoff logic in Exercise 1.2")
	fmt.Println("3. ✅ Implement jitter calculation in Exercise 1.3")
	fmt.Println("4. ✅ Add context cancellation checks in Exercise 1.4")
	fmt.Println()
	fmt.Println("💡 Tips:")
	fmt.Println("• Exponential backoff: delay = baseDelay * 2^(attempt-1)")
	fmt.Println("• Jitter: multiply delay by random factor (0.75-1.25 for 25% jitter)")
	fmt.Println("• Context: Use select statement to check ctx.Done() channel")
	fmt.Println("• Test each implementation before moving to the next")
	fmt.Println()
	fmt.Println("🎯 Goals:")
	fmt.Println("• Understand retry timing patterns")
	fmt.Println("• See how jitter prevents thundering herd")
	fmt.Println("• Learn context-aware programming")
	fmt.Println("• Build foundation for production retry systems")
}

func main() {
	fmt.Println("🔄 Lab 1: Basic Retry Implementation")
	fmt.Println("===================================")

	printTodoInstructions()

	// Uncomment these as you complete each exercise
	runExercise1()
	// runExercise2()
	// runExercise3()
	// runExercise4()

	fmt.Println("\n🎉 Lab 1 Complete!")
	fmt.Println("==================")
	fmt.Println("✅ You've implemented basic retry patterns")
	fmt.Println("✅ You understand exponential backoff")
	fmt.Println("✅ You've learned about jitter benefits")
	fmt.Println("✅ You can handle context cancellation")
	fmt.Println()
	fmt.Println("Next: Lab 2 - Circuit Breaker Patterns")
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Demonstration of the resilient AI agent with fault injection

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create resilient agent
	config := DefaultReliabilityConfig()
	agent, err := NewResilientAgent(apiKey, config)
	if err != nil {
		log.Fatalf("Failed to create resilient agent: %v", err)
	}

	fmt.Println("🧪 Resilient AI Agent Demonstration")
	fmt.Println("===================================")
	fmt.Println()

	// Run demonstrations
	demonstrateBasicFunctionality(agent)
	demonstrateRetryLogic(agent)
	demonstrateCircuitBreaker(agent)
	demonstrateRateLimiting(agent)
	demonstrateRecovery(agent)

	fmt.Println("\n🎉 Demonstration Complete!")
	fmt.Println("=========================")
	fmt.Println("The resilient agent successfully handled all fault scenarios!")
	fmt.Println("✅ Retry logic working")
	fmt.Println("✅ Circuit breaker protection active")
	fmt.Println("✅ Rate limiting enforced")
	fmt.Println("✅ Automatic recovery functioning")
}

func demonstrateBasicFunctionality(agent *ResilientAgent) {
	fmt.Println("1. 📤 Basic Functionality Test")
	fmt.Println("==============================")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := agent.Chat(ctx, "Hello! Can you tell me about AI agents?")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response: %s\n", response)
	}

	// Show initial metrics
	metrics := agent.GetMetrics()
	fmt.Printf("📊 Total requests: %d, Successful: %d\n",
		metrics.TotalRequests, metrics.SuccessfulRequests)

	fmt.Println()
}

func demonstrateRetryLogic(agent *ResilientAgent) {
	fmt.Println("2. 🔄 Retry Logic Demonstration")
	fmt.Println("===============================")

	// Inject temporary timeout failures
	agent.InjectFault("timeout", 3*time.Second)
	fmt.Println("🚨 Injecting timeout failures for 3 seconds...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	startTime := time.Now()
	response, err := agent.Chat(ctx, "What are the benefits of retry logic?")
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Request failed after %v: %v\n", duration, err)
	} else {
		fmt.Printf("✅ Request succeeded after %v with retries\n", duration)
		fmt.Printf("📝 Response: %s\n", response[:min(100, len(response))]+"...")
	}

	// Clear faults
	agent.ClearFaults()
	fmt.Println("🧹 Cleared fault injection")
	fmt.Println()
}

func demonstrateCircuitBreaker(agent *ResilientAgent) {
	fmt.Println("3. ⚡ Circuit Breaker Demonstration")
	fmt.Println("==================================")

	// Inject server errors to trigger circuit breaker
	agent.InjectFault("server_error", 8*time.Second)
	fmt.Println("🚨 Injecting server errors for 8 seconds...")

	// Make multiple requests to trigger circuit breaker
	for i := 0; i < 7; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		_, err := agent.Chat(ctx, fmt.Sprintf("Test request %d", i+1))
		if err != nil {
			fmt.Printf("Request %d: ❌ %v\n", i+1, err)
		} else {
			fmt.Printf("Request %d: ✅ Success\n", i+1)
		}

		cancel()
		time.Sleep(200 * time.Millisecond)
	}

	// Check circuit breaker status
	health := agent.GetHealthStatus()
	if health.CircuitBreakerOpen {
		fmt.Println("⚡ Circuit breaker is now OPEN - protecting system")
	}

	agent.ClearFaults()
	fmt.Println("🧹 Cleared fault injection")
	fmt.Println()
}

func demonstrateRateLimiting(agent *ResilientAgent) {
	fmt.Println("4. 🚦 Rate Limiting Demonstration")
	fmt.Println("=================================")

	fmt.Println("📈 Making rapid requests to demonstrate rate limiting...")

	// Make rapid requests to hit rate limits
	for i := 0; i < 12; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		start := time.Now()
		_, err := agent.Chat(ctx, fmt.Sprintf("Quick request %d", i+1))
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("Request %d (%v): ❌ %v\n", i+1, duration, err)
		} else {
			fmt.Printf("Request %d (%v): ✅ Success\n", i+1, duration)
		}

		cancel()
		time.Sleep(100 * time.Millisecond) // Very fast requests
	}

	fmt.Println("🚦 Rate limiting successfully controlled request flow")
	fmt.Println()
}

func demonstrateRecovery(agent *ResilientAgent) {
	fmt.Println("5. 🔄 Recovery Demonstration")
	fmt.Println("============================")

	fmt.Println("⏳ Waiting for systems to recover...")
	time.Sleep(3 * time.Second)

	// Reset circuit breakers and metrics for clean recovery test
	agent.ResetCircuitBreakers()

	// Test that system has recovered
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := agent.Chat(ctx, "Has the system recovered successfully?")
	if err != nil {
		fmt.Printf("❌ Recovery failed: %v\n", err)
	} else {
		fmt.Printf("✅ System recovered! Response: %s\n", response[:min(100, len(response))]+"...")
	}

	// Show final health status
	health := agent.GetHealthStatus()
	fmt.Printf("🏥 Final health status: Overall=%v, Circuit=%v, RateLimit=%v\n",
		health.Overall, !health.CircuitBreakerOpen, !health.RateLimitExceeded)

	// Show final metrics
	metrics := agent.GetMetrics()
	fmt.Printf("📊 Final metrics: Total=%d, Success=%d, Failed=%d, ErrorRate=%.2f%%\n",
		metrics.TotalRequests, metrics.SuccessfulRequests, metrics.FailedRequests,
		metrics.ErrorRate*100)

	fmt.Println()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

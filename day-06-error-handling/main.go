package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create resilient agent with comprehensive error handling
	config := DefaultReliabilityConfig()
	agent, err := NewResilientAgent(apiKey, config)
	if err != nil {
		log.Fatalf("Failed to create resilient agent: %v", err)
	}

	fmt.Println("🛡️ Production-Ready AI Agent with Error Handling")
	fmt.Println("==============================================")
	fmt.Println()
	fmt.Println("Features:")
	fmt.Println("• Intelligent retry strategies with exponential backoff")
	fmt.Println("• Circuit breakers for fault tolerance")
	fmt.Println("• Rate limiting and quota management")
	fmt.Println("• Real-time monitoring and health checks")
	fmt.Println("• Graceful error recovery")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("• 'stats' - View system health and metrics")
	fmt.Println("• 'health' - Check component health status")
	fmt.Println("• 'config' - Show current reliability configuration")
	fmt.Println("• 'test [scenario]' - Run fault injection tests")
	fmt.Println("• 'demo' - Run comprehensive reliability demonstration")
	fmt.Println("• 'reset' - Reset all circuit breakers and metrics")
	fmt.Println("• 'quit' - Exit the program")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("💬 You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle special commands
		switch {
		case input == "quit":
			fmt.Println("👋 Goodbye!")
			return

		case input == "stats":
			displaySystemStats(agent)
			continue

		case input == "health":
			displayHealthStatus(agent)
			continue

		case input == "config":
			displayConfiguration(agent)
			continue

		case strings.HasPrefix(input, "test "):
			scenario := strings.TrimPrefix(input, "test ")
			runFaultInjectionTest(agent, scenario)
			continue

		case input == "demo":
			fmt.Println("🚀 Starting comprehensive reliability demonstration...")
			runDemo(agent)
			continue

		case input == "reset":
			agent.ResetCircuitBreakers()
			agent.ResetMetrics()
			fmt.Println("✅ System reset completed")
			continue
		}

		// Process regular chat message with full error handling
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		startTime := time.Now()
		response, err := agent.Chat(ctx, input)
		duration := time.Since(startTime)

		cancel()

		if err != nil {
			handleChatError(err, duration)
		} else {
			fmt.Printf("🤖 AI: %s\n", response)
			fmt.Printf("⏱️  Response time: %v\n", duration.Round(time.Millisecond))
		}

		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

func displaySystemStats(agent *ResilientAgent) {
	fmt.Println("\n📊 System Statistics")
	fmt.Println("==================")

	metrics := agent.GetMetrics()

	fmt.Printf("🔄 Requests:\n")
	fmt.Printf("  Total: %d\n", metrics.TotalRequests)
	fmt.Printf("  Successful: %d\n", metrics.SuccessfulRequests)
	fmt.Printf("  Failed: %d\n", metrics.FailedRequests)
	fmt.Printf("  Error Rate: %.2f%%\n", metrics.ErrorRate*100)

	fmt.Printf("\n⏱️  Performance:\n")
	fmt.Printf("  Avg Response Time: %v\n", metrics.AvgResponseTime.Round(time.Millisecond))
	fmt.Printf("  P95 Response Time: %v\n", metrics.P95ResponseTime.Round(time.Millisecond))
	fmt.Printf("  Fastest Response: %v\n", metrics.FastestResponse.Round(time.Millisecond))
	fmt.Printf("  Slowest Response: %v\n", metrics.SlowestResponse.Round(time.Millisecond))

	fmt.Printf("\n🔄 Retries:\n")
	fmt.Printf("  Total Retries: %d\n", metrics.TotalRetries)
	fmt.Printf("  Successful Retries: %d\n", metrics.SuccessfulRetries)
	fmt.Printf("  Failed Retries: %d\n", metrics.FailedRetries)
	fmt.Printf("  Retry Success Rate: %.2f%%\n", metrics.RetrySuccessRate*100)

	fmt.Printf("\n⚡ Circuit Breakers:\n")
	fmt.Printf("  Trips: %d\n", metrics.CircuitBreakerTrips)
	fmt.Printf("  Current State: %s\n", metrics.CircuitBreakerState)
	fmt.Printf("  Time Since Last Trip: %v\n", time.Since(metrics.LastCircuitBreakerTrip).Round(time.Second))

	fmt.Printf("\n🚦 Rate Limiting:\n")
	fmt.Printf("  Requests/Min: %.1f\n", metrics.RequestsPerMinute)
	fmt.Printf("  Rate Limited: %d\n", metrics.RateLimitedRequests)
	fmt.Printf("  Current Quota Usage: %.1f%%\n", metrics.QuotaUsage*100)
}

func displayHealthStatus(agent *ResilientAgent) {
	fmt.Println("\n🏥 Health Status")
	fmt.Println("===============")

	health := agent.GetHealthStatus()

	overallStatus := "🟢 HEALTHY"
	if !health.Overall {
		overallStatus = "🔴 UNHEALTHY"
	}
	fmt.Printf("Overall: %s\n", overallStatus)

	fmt.Printf("\n📡 API Connection:\n")
	if health.APIConnection {
		fmt.Printf("  Status: 🟢 Connected\n")
		fmt.Printf("  Last Success: %v ago\n", time.Since(health.LastAPISuccess).Round(time.Second))
	} else {
		fmt.Printf("  Status: 🔴 Disconnected\n")
		fmt.Printf("  Last Failure: %v ago\n", time.Since(health.LastAPIFailure).Round(time.Second))
	}

	fmt.Printf("\n⚡ Circuit Breaker:\n")
	circuitStatus := "🟢 CLOSED"
	if health.CircuitBreakerOpen {
		circuitStatus = "🔴 OPEN"
	}
	fmt.Printf("  State: %s\n", circuitStatus)
	fmt.Printf("  Failure Count: %d\n", health.ConsecutiveFailures)

	fmt.Printf("\n🚦 Rate Limiter:\n")
	rateLimitStatus := "🟢 AVAILABLE"
	if health.RateLimitExceeded {
		rateLimitStatus = "🟡 LIMITED"
	}
	fmt.Printf("  Status: %s\n", rateLimitStatus)
	fmt.Printf("  Tokens Available: %d\n", health.AvailableTokens)

	fmt.Printf("\n💾 Memory Usage:\n")
	fmt.Printf("  Heap Size: %.2f MB\n", health.MemoryUsage/1024/1024)
	fmt.Printf("  Goroutines: %d\n", health.GoroutineCount)
}

func displayConfiguration(agent *ResilientAgent) {
	fmt.Println("\n⚙️ Reliability Configuration")
	fmt.Println("============================")

	config := agent.GetConfig()

	fmt.Printf("🔄 Retry Settings:\n")
	fmt.Printf("  Max Attempts: %d\n", config.Retry.MaxAttempts)
	fmt.Printf("  Base Delay: %v\n", config.Retry.BaseDelay)
	fmt.Printf("  Max Delay: %v\n", config.Retry.MaxDelay)
	fmt.Printf("  Backoff Multiplier: %.1fx\n", config.Retry.BackoffMultiplier)
	fmt.Printf("  Jitter: %d%%\n", config.Retry.JitterPercent)

	fmt.Printf("\n⚡ Circuit Breaker:\n")
	fmt.Printf("  Failure Threshold: %d\n", config.CircuitBreaker.FailureThreshold)
	fmt.Printf("  Recovery Timeout: %v\n", config.CircuitBreaker.RecoveryTimeout)
	fmt.Printf("  Test Request Rate: %.1f%%\n", config.CircuitBreaker.TestRequestRate*100)

	fmt.Printf("\n🚦 Rate Limiting:\n")
	fmt.Printf("  Requests/Min: %d\n", config.RateLimit.RequestsPerMinute)
	fmt.Printf("  Burst Size: %d\n", config.RateLimit.BurstSize)
	fmt.Printf("  Adaptive: %t\n", config.RateLimit.AdaptiveEnabled)

	fmt.Printf("\n📊 Monitoring:\n")
	fmt.Printf("  Metrics Enabled: %t\n", config.Monitoring.MetricsEnabled)
	fmt.Printf("  Health Checks: %t\n", config.Monitoring.HealthChecksEnabled)
	fmt.Printf("  Alert Threshold: %.1f%%\n", config.Monitoring.AlertThreshold*100)
}

func runFaultInjectionTest(agent *ResilientAgent, scenario string) {
	fmt.Printf("\n🧪 Running Fault Injection Test: %s\n", scenario)
	fmt.Println("=========================================")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	switch scenario {
	case "timeout":
		fmt.Println("Simulating API timeout...")
		agent.InjectFault("timeout", 5*time.Second)

	case "ratelimit":
		fmt.Println("Simulating rate limit...")
		agent.InjectFault("ratelimit", 10*time.Second)

	case "server_error":
		fmt.Println("Simulating server errors...")
		agent.InjectFault("server_error", 8*time.Second)

	case "network":
		fmt.Println("Simulating network issues...")
		agent.InjectFault("network", 15*time.Second)

	case "quota":
		fmt.Println("Simulating quota exhaustion...")
		agent.InjectFault("quota", 20*time.Second)

	default:
		fmt.Printf("Unknown scenario: %s\n", scenario)
		fmt.Println("Available scenarios: timeout, ratelimit, server_error, network, quota")
		return
	}

	// Test the agent under fault conditions
	testMessages := []string{
		"Hello, how are you?",
		"What's the weather like?",
		"Tell me a joke",
		"How do circuit breakers work?",
		"What's your favorite color?",
	}

	successCount := 0
	for i, msg := range testMessages {
		fmt.Printf("\nTest %d/5: %s\n", i+1, msg)

		response, err := agent.Chat(ctx, msg)
		if err != nil {
			fmt.Printf("❌ Failed: %v\n", err)
		} else {
			fmt.Printf("✅ Success: %s\n", response[:min(100, len(response))]+"...")
			successCount++
		}

		time.Sleep(1 * time.Second) // Brief pause between tests
	}

	agent.ClearFaults()

	fmt.Printf("\n📊 Test Results:\n")
	fmt.Printf("  Success Rate: %d/5 (%.0f%%)\n", successCount, float64(successCount)/5*100)
	fmt.Printf("  Fault injection cleared\n")
}

func handleChatError(err error, duration time.Duration) {
	fmt.Printf("❌ Error: %v\n", err)
	fmt.Printf("⏱️  Failed after: %v\n", duration.Round(time.Millisecond))

	// Provide helpful guidance based on error type
	switch {
	case strings.Contains(err.Error(), "rate limit"):
		fmt.Println("💡 Tip: Rate limiting is active. The system will automatically retry with backoff.")

	case strings.Contains(err.Error(), "timeout"):
		fmt.Println("💡 Tip: Request timed out. Check your network connection or try again.")

	case strings.Contains(err.Error(), "circuit breaker"):
		fmt.Println("💡 Tip: Circuit breaker is open. The system is protecting against cascading failures.")

	case strings.Contains(err.Error(), "quota"):
		fmt.Println("💡 Tip: API quota may be exhausted. Check your OpenAI usage limits.")

	default:
		fmt.Println("💡 Tip: Use 'health' command to check system status or 'stats' for detailed metrics.")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

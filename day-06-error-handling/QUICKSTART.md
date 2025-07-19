# Quick Start: Day 6 Error Handling & Retries

## 🚀 Get Started in 5 Minutes

### 1. **Environment Setup**
```bash
# Navigate to Day 6 directory
cd day-06-error-handling

# Copy environment template
cp .env.example .env

# Edit with your OpenAI API key
echo "OPENAI_API_KEY=your_api_key_here" > .env
```

### 2. **Install Dependencies**
```bash
# Initialize Go module
go mod init day-06-error-handling

# Install required packages
go get github.com/joho/godotenv@latest
go get github.com/sashabaranov/go-openai@latest

# Download dependencies
go mod tidy
```

### 3. **Run the Resilient Agent**
```bash
# Start the production-ready AI agent
go run *.go
```

## 💬 Try These Commands

### **Basic Chat**
```
You: Hello, how are you?
```

### **System Health Check**
```
You: health
```

### **View Performance Metrics**
```
You: stats
```

### **Test Fault Tolerance**
```
You: test timeout
You: test ratelimit
You: test server_error
```

### **View Configuration**
```
You: config
```

### **Reset System State**
```
You: reset
```

## 🧪 Fault Injection Testing

Test how the system handles various failure scenarios:

### **API Timeout Simulation**
```bash
You: test timeout
# Watch how the system handles timeouts with retries
```

### **Rate Limit Testing**
```bash
You: test ratelimit
# See exponential backoff in action
```

### **Server Error Handling**
```bash
You: test server_error
# Observe circuit breaker protection
```

### **Network Issues**
```bash
You: test network
# Test network failure recovery
```

### **Quota Exhaustion**
```bash
You: test quota
# Simulate API quota limits
```

## 📊 Understanding the Output

### **Success Response**
```
🤖 AI: Hello! I'm doing well, thank you for asking.
⏱️  Response time: 234ms
```

### **Error with Recovery Information**
```
❌ Error: rate_limit: rate limit exceeded
⏱️  Failed after: 1.2s
💡 Tip: Rate limiting is active. The system will automatically retry with backoff.
```

### **Health Status Display**
```
🏥 Health Status
===============
Overall: 🟢 HEALTHY

📡 API Connection:
  Status: 🟢 Connected
  Last Success: 2s ago

⚡ Circuit Breaker:
  State: 🟢 CLOSED
  Failure Count: 0

🚦 Rate Limiter:
  Status: 🟢 AVAILABLE
  Tokens Available: 8
```

### **Performance Metrics**
```
📊 System Statistics
==================
🔄 Requests:
  Total: 25
  Successful: 23
  Failed: 2
  Error Rate: 8.00%

⏱️  Performance:
  Avg Response Time: 456ms
  P95 Response Time: 1.2s
  Fastest Response: 123ms
  Slowest Response: 2.1s
```

## 🎯 Key Features in Action

### **1. Exponential Backoff**
- Watch delays increase: 100ms → 200ms → 400ms
- Jitter prevents thundering herd problems
- Respects maximum delay limits

### **2. Circuit Breaker Protection**
- Opens after 5 consecutive failures
- Prevents cascading failures
- Automatic recovery testing

### **3. Intelligent Rate Limiting**
- Token bucket with burst capacity
- Adaptive adjustment based on conditions
- Quota usage monitoring

### **4. Real-time Monitoring**
- Request success/failure rates
- Response time percentiles
- System health indicators

## 🔧 Configuration Customization

### **Modify Retry Behavior**
Edit `resilient_agent.go` and adjust:
```go
RetryConfig{
    MaxAttempts: 5,           // More retry attempts
    BaseDelay: 50 * time.Millisecond,  // Faster initial retry
    BackoffMultiplier: 1.5,   // Gentler backoff
}
```

### **Tune Circuit Breaker**
```go
CircuitBreakerConfig{
    FailureThreshold: 10,     // More tolerant
    RecoveryTimeout: 30 * time.Second,  // Faster recovery
}
```

### **Adjust Rate Limits**
```go
RateLimitConfig{
    RequestsPerMinute: 120,   // Higher rate limit
    BurstSize: 20,           // Larger burst capacity
}
```

## 🚨 Troubleshooting

### **"API key not found" Error**
```bash
# Check environment file
cat .env

# Verify API key format
echo $OPENAI_API_KEY
```

### **Import Errors**
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod tidy
```

### **Circuit Breaker Stuck Open**
```bash
# Reset the system
You: reset

# Check health status
You: health
```

### **Rate Limit Issues**
```bash
# Check current limits
You: config

# Wait for token refill
# Or increase burst size in config
```

## 📈 Performance Tips

### **Optimize for Your Use Case**
- **High Traffic**: Increase rate limits and burst size
- **Reliability Focus**: Lower failure thresholds
- **Cost Optimization**: Reduce retry attempts
- **Low Latency**: Decrease base retry delays

### **Production Checklist**
- [ ] Set appropriate rate limits for your API tier
- [ ] Configure circuit breaker thresholds for your SLA
- [ ] Set up monitoring alerts for error rates
- [ ] Test fault injection scenarios regularly
- [ ] Document operational procedures

## 🎓 Next Steps

### **Explore Advanced Features**
1. **Custom Error Handling**: Add domain-specific error types
2. **Distributed Coordination**: Share circuit breaker state
3. **Advanced Monitoring**: Add Prometheus metrics
4. **Cost Optimization**: Implement adaptive rate limiting

### **Integration Opportunities**
- **Day 5 Memory**: Add memory persistence during failures
- **Day 4 Prompts**: Fallback to simpler prompts
- **Day 7 Project**: Use as foundation for robust chatbot

## 🏆 Success Indicators

You'll know the system is working when:
- ✅ **Error Recovery**: Failed requests automatically retry
- ✅ **Circuit Protection**: System prevents cascading failures
- ✅ **Rate Compliance**: Stays within API limits
- ✅ **Monitoring**: Real-time visibility into system health
- ✅ **Graceful Degradation**: Maintains service during issues

---

**Ready to build bulletproof AI systems? Start chatting and test the fault tolerance! 🛡️**

# Day 6 Demo Instructions

## 🎮 How to Run the Demo

### 1. **Setup Environment**
```bash
cd day-06-error-handling
cp .env.example .env
# Edit .env with your OpenAI API key
```

### 2. **Run the Main Application**
```bash
go run *.go
```

### 3. **Run Comprehensive Demo**

#### **Automatic Demonstration**
```
You: demo
```
This will run a complete fault tolerance demonstration showing:
- Basic functionality test
- Retry logic with exponential backoff
- Circuit breaker protection
- Rate limiting enforcement
- System recovery validation

### 4. **Interactive Testing**

#### **Basic Chat Test**
```
You: Hello, how are you doing?
```

#### **View System Health**
```
You: health
```

#### **Check Performance Metrics**
```
You: stats
```

#### **Test Individual Fault Scenarios**
```
You: test timeout
You: test ratelimit
You: test server_error
You: test network
```

#### **Reset System State**
```
You: reset
```

## 🧪 What You'll See

### **Successful Response**
```
🤖 AI: Hello! I'm doing well, thank you for asking. How can I assist you today?
⏱️  Response time: 234ms
```

### **Error with Recovery**
```
❌ Error: rate_limit: rate limit exceeded
⏱️  Failed after: 1.2s
💡 Tip: Rate limiting is active. The system will automatically retry with backoff.
```

### **System Health Status**
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
```

## 🎯 Key Demonstrations

1. **Retry Logic**: Watch exponential backoff in action
2. **Circuit Breaker**: See system protection during failures
3. **Rate Limiting**: Observe intelligent request throttling
4. **Recovery**: Experience automatic system healing
5. **Monitoring**: Real-time visibility into system health

## 🚀 Advanced Testing

For comprehensive testing, check out the lab exercises:

```bash
# Test retry patterns
cd lab/lab1-retry && go run *.go

# Test circuit breakers  
cd lab/lab2-circuit && go run *.go

# Test rate limiting
cd lab/lab3-ratelimit && go run *.go

# Test full integration
cd lab/lab4-integration && go run *.go
```

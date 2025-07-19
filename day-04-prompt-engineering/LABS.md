# Day 4 Labs: Prompt Engineering Exercises

## Lab 1: Basic Template System

**Objective**: Build your own simple prompt template system

### Exercise 1.1: String Template
Create a function that takes a template string and replaces placeholders:

```go
func SimpleTemplate(template string, vars map[string]string) string {
    // Replace {{variable}} with actual values
    // Example: "Hello {{name}}" with vars["name"] = "World" -> "Hello World"
}
```

### Exercise 1.2: Template Validation
Add validation to ensure all required variables are provided:

```go
func ValidateTemplate(template string, vars map[string]string) []string {
    // Return list of missing variables
}
```

## Lab 2: Prompting Techniques

### Exercise 2.1: Zero-Shot vs Few-Shot
Compare these two approaches for the same task:

**Zero-shot prompt:**
```
"Classify this email as spam or not spam: [email content]"
```

**Few-shot prompt:**
```
"Classify emails as spam or not spam:

Example 1: 'Win money now!' -> spam
Example 2: 'Meeting at 3pm' -> not spam
Example 3: 'Free pills!!!' -> spam

Now classify: [email content]"
```

### Exercise 2.2: Chain-of-Thought
Create a prompt that guides the model through step-by-step reasoning:

```
"Solve this math problem step by step:
1. Identify what we're looking for
2. List the given information  
3. Choose the appropriate formula
4. Substitute values
5. Calculate the result
6. Verify the answer makes sense

Problem: [math problem]"
```

## Lab 3: Domain-Specific Prompts

### Exercise 3.1: Code Review Prompt
Design a prompt for code review that covers:
- Code quality
- Best practices
- Potential bugs
- Suggestions for improvement

### Exercise 3.2: Data Analysis Prompt
Create a template for analyzing datasets:
- Summary statistics
- Key insights
- Trends and patterns
- Recommendations

## Lab 4: Prompt Optimization

### Exercise 4.1: Response Quality Metrics
Implement functions to score prompt responses:

```go
// Measures response completeness
func CompletenessScore(response string, expectedElements []string) float64

// Measures response relevance
func RelevanceScore(response string, query string) float64

// Measures response clarity
func ClarityScore(response string) float64
```

### Exercise 4.2: Iterative Improvement
Design an experiment to improve a prompt through iterations:

1. Start with a baseline prompt
2. Identify weaknesses in responses
3. Modify the prompt to address issues
4. Test and measure improvement
5. Repeat until satisfied

## Lab 5: Advanced Techniques

### Exercise 5.1: Context-Aware Prompting
Build a system that adapts prompts based on:
- User expertise level
- Previous conversation history
- Task complexity
- Available context length

### Exercise 5.2: Multi-Turn Conversations
Design prompts that work well in conversational contexts:
- Reference previous messages
- Maintain context coherence
- Handle topic transitions
- Manage conversation flow

## Evaluation Criteria

For each lab, evaluate your solutions on:

1. **Effectiveness**: Does it achieve the intended goal?
2. **Clarity**: Are the prompts clear and unambiguous?
3. **Consistency**: Do similar inputs produce similar outputs?
4. **Efficiency**: Token usage vs. output quality
5. **Robustness**: How well does it handle edge cases?

## Bonus Challenges

### Challenge 1: Prompt Marketplace
Design a system for sharing and rating prompt templates:
- Template categorization
- User ratings and reviews
- Usage statistics
- Best practice guidelines

### Challenge 2: Adaptive Prompting
Create a system that automatically improves prompts based on user feedback:
- Collect response ratings
- Identify common failure patterns
- Suggest prompt modifications
- A/B test improvements

### Challenge 3: Multi-Language Prompts
Extend your templates to support multiple languages:
- Language detection
- Template translation
- Cultural adaptation
- Localized examples

## Success Metrics

By the end of these labs, you should be able to:
- ✅ Create effective prompt templates
- ✅ Apply various prompting techniques appropriately
- ✅ Measure and improve prompt effectiveness
- ✅ Design domain-specific prompting strategies
- ✅ Handle complex, multi-turn conversations

## Reflection Questions

1. Which prompting techniques work best for different types of tasks?
2. How do you balance prompt specificity with flexibility?
3. What role does context play in prompt effectiveness?
4. How can you systematically improve prompts over time?
5. What are the trade-offs between prompt complexity and results?

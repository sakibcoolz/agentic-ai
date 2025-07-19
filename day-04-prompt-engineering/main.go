package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// PromptTemplate represents a reusable prompt template
type PromptTemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Template    string                 `json:"template"`
	Variables   []string               `json:"variables"`
	Category    string                 `json:"category"`
	Examples    []PromptExample        `json:"examples"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PromptExample shows how to use a template
type PromptExample struct {
	Input       map[string]string `json:"input"`
	Description string            `json:"description"`
}

// PromptEngine manages prompt templates and generation
type PromptEngine struct {
	templates map[string]PromptTemplate
	client    *openai.Client
	history   []PromptExecution
}

// PromptExecution tracks prompt usage and results
type PromptExecution struct {
	Template        string                 `json:"template"`
	Variables       map[string]string      `json:"variables"`
	GeneratedPrompt string                 `json:"generated_prompt"`
	Response        string                 `json:"response"`
	Timestamp       time.Time              `json:"timestamp"`
	TokensUsed      int                    `json:"tokens_used"`
	Quality         float64                `json:"quality"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NewPromptEngine creates a new prompt engineering system
func NewPromptEngine(apiKey string) *PromptEngine {
	engine := &PromptEngine{
		templates: make(map[string]PromptTemplate),
		client:    openai.NewClient(apiKey),
		history:   make([]PromptExecution, 0),
	}

	// Load built-in templates
	engine.loadBuiltinTemplates()

	return engine
}

// loadBuiltinTemplates adds pre-defined prompt templates
func (pe *PromptEngine) loadBuiltinTemplates() {
	// Code generation template
	pe.AddTemplate(PromptTemplate{
		Name:        "code_generation",
		Description: "Generate Go code based on requirements",
		Category:    "programming",
		Template: `You are an expert Go programmer. Generate clean, efficient, and well-documented Go code.

Task: {{.task}}
Requirements:
{{range .requirements}}
- {{.}}
{{end}}

Additional Context: {{.context}}

Please provide:
1. Complete, runnable Go code
2. Inline comments explaining key logic
3. Error handling where appropriate
4. Follow Go best practices and conventions

Code:`,
		Variables: []string{"task", "requirements", "context"},
		Examples: []PromptExample{
			{
				Input: map[string]string{
					"task":         "Create a function to calculate Fibonacci numbers",
					"requirements": "Efficient algorithm,Handle edge cases,Include tests",
					"context":      "Part of a math utility package",
				},
				Description: "Generate Fibonacci function with tests",
			},
		},
	})

	// Data analysis template
	pe.AddTemplate(PromptTemplate{
		Name:        "data_analysis",
		Description: "Analyze data and provide insights",
		Category:    "analysis",
		Template: `You are a senior data analyst with expertise in {{.domain}}. Analyze the following data and provide comprehensive insights.

Data: {{.data}}
Analysis Type: {{.analysis_type}}
Business Context: {{.context}}

Please provide:
1. **Key Findings**: What are the main insights?
2. **Trends & Patterns**: What trends do you observe?
3. **Anomalies**: Any unusual patterns or outliers?
4. **Recommendations**: Actionable next steps
5. **Confidence Level**: How confident are you in these insights?

Format your response with clear sections and bullet points.`,
		Variables: []string{"domain", "data", "analysis_type", "context"},
		Examples: []PromptExample{
			{
				Input: map[string]string{
					"domain":        "e-commerce",
					"data":          "Monthly sales data showing 20% increase",
					"analysis_type": "trend analysis",
					"context":       "Q4 holiday season performance",
				},
				Description: "Analyze e-commerce sales trends",
			},
		},
	})

	// Chain-of-thought problem solving
	pe.AddTemplate(PromptTemplate{
		Name:        "chain_of_thought",
		Description: "Step-by-step problem solving with reasoning",
		Category:    "reasoning",
		Template: `Let's solve this problem step by step using clear reasoning.

Problem: {{.problem}}
Context: {{.context}}
Constraints: {{.constraints}}

I'll work through this systematically:

Step 1: Understand the Problem
- What exactly are we trying to solve?
- What information do we have?
- What are we looking for?

Step 2: Break Down the Problem
- What are the key components?
- Are there sub-problems to solve first?
- What's the logical sequence?

Step 3: Apply Solution Strategy
- What approach should we use?
- Why is this the best method?
- How do we implement it?

Step 4: Verify the Solution
- Does this make sense?
- Have we addressed all requirements?
- Are there any edge cases?

Let me work through each step:`,
		Variables: []string{"problem", "context", "constraints"},
		Examples: []PromptExample{
			{
				Input: map[string]string{
					"problem":     "Optimize database query performance",
					"context":     "E-commerce application with slow product searches",
					"constraints": "Cannot change database schema",
				},
				Description: "Database optimization problem",
			},
		},
	})

	// Few-shot learning template
	pe.AddTemplate(PromptTemplate{
		Name:        "few_shot_learning",
		Description: "Learn from examples to perform similar tasks",
		Category:    "learning",
		Template: `I'll show you some examples of {{.task_type}}, then ask you to do a similar task.

{{range .examples}}
Example {{.number}}:
Input: {{.input}}
Output: {{.output}}
Explanation: {{.explanation}}

{{end}}

Now, please apply the same pattern to this new case:
Input: {{.new_input}}
Output:`,
		Variables: []string{"task_type", "examples", "new_input"},
		Examples: []PromptExample{
			{
				Input: map[string]string{
					"task_type": "function naming in Go",
					"new_input": "Function that converts string to uppercase",
				},
				Description: "Learn Go naming conventions from examples",
			},
		},
	})

	// Creative writing template
	pe.AddTemplate(PromptTemplate{
		Name:        "creative_writing",
		Description: "Generate creative content with specific style and requirements",
		Category:    "creative",
		Template: `You are a talented {{.writer_type}} with expertise in {{.domain}}.

Writing Task: {{.task}}
Style: {{.style}}
Tone: {{.tone}}
Target Audience: {{.audience}}
Length: {{.length}}

Key Requirements:
{{range .requirements}}
- {{.}}
{{end}}

Theme/Message: {{.theme}}

Please create engaging content that:
1. Captures the reader's attention immediately
2. Maintains the specified tone throughout
3. Delivers the core message effectively
4. Is appropriate for the target audience
5. Follows the style guidelines

Content:`,
		Variables: []string{"writer_type", "domain", "task", "style", "tone", "audience", "length", "requirements", "theme"},
		Examples: []PromptExample{
			{
				Input: map[string]string{
					"writer_type": "technical blogger",
					"domain":      "software development",
					"task":        "Explain microservices architecture",
					"style":       "conversational yet informative",
					"tone":        "friendly and approachable",
					"audience":    "junior developers",
					"length":      "800-1000 words",
					"theme":       "Making complex concepts accessible",
				},
				Description: "Technical blog post about microservices",
			},
		},
	})
}

// AddTemplate adds a new template to the engine
func (pe *PromptEngine) AddTemplate(template PromptTemplate) {
	pe.templates[template.Name] = template
}

// GetTemplate retrieves a template by name
func (pe *PromptEngine) GetTemplate(name string) (PromptTemplate, error) {
	template, exists := pe.templates[name]
	if !exists {
		return PromptTemplate{}, fmt.Errorf("template '%s' not found", name)
	}
	return template, nil
}

// ListTemplates returns all available templates
func (pe *PromptEngine) ListTemplates() map[string]PromptTemplate {
	return pe.templates
}

// GeneratePrompt creates a prompt from a template with variables
func (pe *PromptEngine) GeneratePrompt(templateName string, variables map[string]interface{}) (string, error) {
	templateObj, err := pe.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	// Create Go template
	tmpl, err := template.New(templateName).Parse(templateObj.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with variables
	var result strings.Builder
	err = tmpl.Execute(&result, variables)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// ExecutePrompt generates and executes a prompt using the LLM
func (pe *PromptEngine) ExecutePrompt(ctx context.Context, templateName string, variables map[string]interface{}) (*PromptExecution, error) {
	// Generate the prompt
	prompt, err := pe.GeneratePrompt(templateName, variables)
	if err != nil {
		return nil, err
	}

	// Convert variables to string map for storage
	stringVars := make(map[string]string)
	for k, v := range variables {
		stringVars[k] = fmt.Sprintf("%v", v)
	}

	// Execute with LLM
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	resp, err := pe.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM execution failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from LLM")
	}

	// Create execution record
	execution := &PromptExecution{
		Template:        templateName,
		Variables:       stringVars,
		GeneratedPrompt: prompt,
		Response:        resp.Choices[0].Message.Content,
		Timestamp:       time.Now(),
		TokensUsed:      resp.Usage.TotalTokens,
		Quality:         0, // To be set by evaluation
		Metadata:        make(map[string]interface{}),
	}

	// Store in history
	pe.history = append(pe.history, *execution)

	return execution, nil
}

// AnalyzePromptEffectiveness provides metrics on prompt usage
func (pe *PromptEngine) AnalyzePromptEffectiveness() map[string]interface{} {
	if len(pe.history) == 0 {
		return map[string]interface{}{
			"total_executions": 0,
			"message":          "No prompt executions recorded yet",
		}
	}

	// Calculate metrics
	totalExecutions := len(pe.history)
	totalTokens := 0
	templateUsage := make(map[string]int)
	avgTokensByTemplate := make(map[string]float64)

	for _, execution := range pe.history {
		totalTokens += execution.TokensUsed
		templateUsage[execution.Template]++
	}

	// Calculate average tokens by template
	for template, count := range templateUsage {
		totalForTemplate := 0
		for _, execution := range pe.history {
			if execution.Template == template {
				totalForTemplate += execution.TokensUsed
			}
		}
		avgTokensByTemplate[template] = float64(totalForTemplate) / float64(count)
	}

	return map[string]interface{}{
		"total_executions":       totalExecutions,
		"total_tokens_used":      totalTokens,
		"average_tokens":         float64(totalTokens) / float64(totalExecutions),
		"template_usage":         templateUsage,
		"avg_tokens_by_template": avgTokensByTemplate,
		"most_used_template":     findMostUsedTemplate(templateUsage),
	}
}

// findMostUsedTemplate finds the template with highest usage
func findMostUsedTemplate(usage map[string]int) string {
	maxUsage := 0
	mostUsed := ""

	for template, count := range usage {
		if count > maxUsage {
			maxUsage = count
			mostUsed = template
		}
	}

	return mostUsed
}

// GetPromptHistory returns execution history
func (pe *PromptEngine) GetPromptHistory() []PromptExecution {
	return pe.history
}

// ValidateTemplate checks if a template has all required components
func (pe *PromptEngine) ValidateTemplate(template PromptTemplate) []string {
	var issues []string

	// Check required fields
	if template.Name == "" {
		issues = append(issues, "Template name is required")
	}
	if template.Template == "" {
		issues = append(issues, "Template content is required")
	}

	// Check for undefined variables in template
	re := regexp.MustCompile(`\{\{\.(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(template.Template, -1)

	templateVars := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			templateVars[match[1]] = true
		}
	}

	// Check if all template variables are declared
	declaredVars := make(map[string]bool)
	for _, v := range template.Variables {
		declaredVars[v] = true
	}

	for templateVar := range templateVars {
		if !declaredVars[templateVar] {
			issues = append(issues, fmt.Sprintf("Variable '%s' used in template but not declared", templateVar))
		}
	}

	return issues
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create prompt engine
	engine := NewPromptEngine(apiKey)
	ctx := context.Background()

	fmt.Println("🎯 Prompt Engineering System")
	fmt.Println("=============================")
	fmt.Printf("Available templates: %d\n\n", len(engine.ListTemplates()))

	// Show available templates
	fmt.Println("📋 Available Templates:")
	for name, template := range engine.ListTemplates() {
		fmt.Printf("- %s (%s): %s\n", name, template.Category, template.Description)
	}

	fmt.Println("\nCommands:")
	fmt.Println("- 'list' - Show all templates")
	fmt.Println("- 'demo <template>' - Run a demo of a template")
	fmt.Println("- 'stats' - Show prompt usage statistics")
	fmt.Println("- 'custom' - Create a custom prompt")
	fmt.Println("- 'quit' - Exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Prompt> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "quit" {
			fmt.Println("👋 Goodbye!")
			break
		}

		parts := strings.Fields(input)
		command := strings.ToLower(parts[0])

		switch command {
		case "list":
			fmt.Println("\n📋 Template Details:")
			for name, template := range engine.ListTemplates() {
				fmt.Printf("\n%s (%s):\n", name, template.Category)
				fmt.Printf("  Description: %s\n", template.Description)
				fmt.Printf("  Variables: %v\n", template.Variables)
				if len(template.Examples) > 0 {
					fmt.Printf("  Example: %s\n", template.Examples[0].Description)
				}
			}

		case "demo":
			if len(parts) < 2 {
				fmt.Println("Usage: demo <template_name>")
				continue
			}

			templateName := parts[1]
			template, err := engine.GetTemplate(templateName)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			if len(template.Examples) == 0 {
				fmt.Printf("No examples available for template '%s'\n", templateName)
				continue
			}

			// Use the first example
			example := template.Examples[0]
			variables := make(map[string]interface{})
			for k, v := range example.Input {
				variables[k] = v
			}

			fmt.Printf("\n🔍 Demo: %s\n", example.Description)
			fmt.Printf("Template: %s\n\n", templateName)

			execution, err := engine.ExecutePrompt(ctx, templateName, variables)
			if err != nil {
				fmt.Printf("Error executing prompt: %v\n", err)
				continue
			}

			fmt.Printf("Generated Prompt:\n%s\n\n", execution.GeneratedPrompt)
			fmt.Printf("Response:\n%s\n\n", execution.Response)
			fmt.Printf("Tokens used: %d\n\n", execution.TokensUsed)

		case "stats":
			stats := engine.AnalyzePromptEffectiveness()
			fmt.Println("\n📊 Prompt Usage Statistics:")
			for key, value := range stats {
				fmt.Printf("  %s: %v\n", key, value)
			}
			fmt.Println()

		case "custom":
			fmt.Println("\n✏️ Custom Prompt Creator")
			fmt.Print("Enter your prompt: ")
			if !scanner.Scan() {
				continue
			}

			customPrompt := scanner.Text()
			if customPrompt == "" {
				fmt.Println("Empty prompt, skipping.")
				continue
			}

			// Execute custom prompt directly
			req := openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: customPrompt,
					},
				},
				Temperature: 0.7,
				MaxTokens:   1000,
			}

			resp, err := engine.client.CreateChatCompletion(ctx, req)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			if len(resp.Choices) > 0 {
				fmt.Printf("\nResponse:\n%s\n\n", resp.Choices[0].Message.Content)
				fmt.Printf("Tokens used: %d\n\n", resp.Usage.TotalTokens)
			}

		default:
			fmt.Println("Unknown command. Try 'list', 'demo <template>', 'stats', 'custom', or 'quit'")
		}
	}

	// Final statistics
	if len(engine.GetPromptHistory()) > 0 {
		fmt.Println("\n📊 Session Summary:")
		stats := engine.AnalyzePromptEffectiveness()
		for key, value := range stats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}

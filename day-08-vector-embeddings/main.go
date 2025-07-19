package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// Embedding represents a text embedding with metadata
type Embedding struct {
	ID       string                 `json:"id"`
	Text     string                 `json:"text"`
	Vector   []float64              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata"`
}

// VectorStore provides in-memory vector storage and search
type VectorStore struct {
	embeddings []Embedding
	client     *openai.Client
}

// SearchResult represents a search result with similarity score
type SearchResult struct {
	Embedding  Embedding `json:"embedding"`
	Similarity float64   `json:"similarity"`
}

// NewVectorStore creates a new vector store
func NewVectorStore(apiKey string) *VectorStore {
	return &VectorStore{
		embeddings: make([]Embedding, 0),
		client:     openai.NewClient(apiKey),
	}
}

// GenerateEmbedding creates an embedding for the given text
func (vs *VectorStore) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	}

	resp, err := vs.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding data returned")
	}

	// Convert float32 to float64
	embedding := resp.Data[0].Embedding
	result := make([]float64, len(embedding))
	for i, v := range embedding {
		result[i] = float64(v)
	}

	return result, nil
}

// AddDocument adds a document to the vector store
func (vs *VectorStore) AddDocument(ctx context.Context, id, text string, metadata map[string]interface{}) error {
	vector, err := vs.GenerateEmbedding(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	embedding := Embedding{
		ID:       id,
		Text:     text,
		Vector:   vector,
		Metadata: metadata,
	}

	vs.embeddings = append(vs.embeddings, embedding)
	return nil
}

// CosineSimilarity calculates cosine similarity between two vectors
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Search performs semantic search in the vector store
func (vs *VectorStore) Search(ctx context.Context, query string, topK int) ([]SearchResult, error) {
	queryVector, err := vs.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	results := make([]SearchResult, 0, len(vs.embeddings))

	for _, embedding := range vs.embeddings {
		similarity := CosineSimilarity(queryVector, embedding.Vector)
		results = append(results, SearchResult{
			Embedding:  embedding,
			Similarity: similarity,
		})
	}

	// Sort by similarity (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Return top K results
	if topK > len(results) {
		topK = len(results)
	}

	return results[:topK], nil
}

// GetDocumentCount returns the number of documents in the store
func (vs *VectorStore) GetDocumentCount() int {
	return len(vs.embeddings)
}

// GetDocument retrieves a document by ID
func (vs *VectorStore) GetDocument(id string) (*Embedding, error) {
	for _, embedding := range vs.embeddings {
		if embedding.ID == id {
			return &embedding, nil
		}
	}
	return nil, fmt.Errorf("document with ID %s not found", id)
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

	// Create vector store
	vectorStore := NewVectorStore(apiKey)
	ctx := context.Background()

	fmt.Println("🔍 Vector Database & Embeddings Demo")
	fmt.Println("=====================================")

	// Sample documents to add to the vector store
	documents := []struct {
		id       string
		text     string
		metadata map[string]interface{}
	}{
		{
			id:   "doc1",
			text: "Artificial intelligence is the simulation of human intelligence in machines that are programmed to think and learn like humans.",
			metadata: map[string]interface{}{
				"category": "AI",
				"source":   "encyclopedia",
			},
		},
		{
			id:   "doc2",
			text: "Machine learning is a subset of artificial intelligence that focuses on the development of algorithms that allow computers to learn from data.",
			metadata: map[string]interface{}{
				"category": "ML",
				"source":   "textbook",
			},
		},
		{
			id:   "doc3",
			text: "Natural language processing enables computers to understand, interpret, and generate human language in a valuable way.",
			metadata: map[string]interface{}{
				"category": "NLP",
				"source":   "research",
			},
		},
		{
			id:   "doc4",
			text: "Deep learning uses neural networks with multiple layers to model and understand complex patterns in data.",
			metadata: map[string]interface{}{
				"category": "DL",
				"source":   "article",
			},
		},
		{
			id:   "doc5",
			text: "Computer vision allows machines to interpret and understand visual information from the world around them.",
			metadata: map[string]interface{}{
				"category": "CV",
				"source":   "journal",
			},
		},
		{
			id:   "doc6",
			text: "Go is a programming language developed by Google that emphasizes simplicity, efficiency, and strong support for concurrent programming.",
			metadata: map[string]interface{}{
				"category": "Programming",
				"source":   "documentation",
			},
		},
	}

	// Add documents to vector store
	fmt.Println("📥 Adding documents to vector store...")
	for _, doc := range documents {
		fmt.Printf("Adding: %s... ", doc.id)
		err := vectorStore.AddDocument(ctx, doc.id, doc.text, doc.metadata)
		if err != nil {
			log.Printf("Error adding document %s: %v", doc.id, err)
			continue
		}
		fmt.Println("✅")
	}

	fmt.Printf("\n📊 Vector store contains %d documents\n\n", vectorStore.GetDocumentCount())

	// Demonstrate semantic search
	queries := []string{
		"What is machine learning?",
		"How do computers understand images?",
		"Programming languages for AI",
		"Neural networks and deep learning",
		"Language understanding by machines",
	}

	for _, query := range queries {
		fmt.Printf("🔍 Query: %s\n", query)
		fmt.Println(strings.Repeat("-", 50))

		results, err := vectorStore.Search(ctx, query, 3)
		if err != nil {
			log.Printf("Search error: %v", err)
			continue
		}

		for i, result := range results {
			fmt.Printf("%d. [%.3f] %s\n",
				i+1,
				result.Similarity,
				result.Embedding.Text)
			fmt.Printf("   ID: %s, Category: %v\n",
				result.Embedding.ID,
				result.Embedding.Metadata["category"])
		}
		fmt.Println()
	}

	// Demonstrate similarity between documents
	fmt.Println("📈 Document Similarity Analysis")
	fmt.Println(strings.Repeat("=", 40))

	doc1, _ := vectorStore.GetDocument("doc1") // AI
	doc2, _ := vectorStore.GetDocument("doc2") // ML
	doc6, _ := vectorStore.GetDocument("doc6") // Go programming

	sim12 := CosineSimilarity(doc1.Vector, doc2.Vector)
	sim16 := CosineSimilarity(doc1.Vector, doc6.Vector)
	sim26 := CosineSimilarity(doc2.Vector, doc6.Vector)

	fmt.Printf("AI ↔ ML similarity: %.3f\n", sim12)
	fmt.Printf("AI ↔ Go similarity: %.3f\n", sim16)
	fmt.Printf("ML ↔ Go similarity: %.3f\n", sim26)

	fmt.Println("\n✨ Vector search demo complete!")
	fmt.Println("Notice how semantically similar documents have higher similarity scores!")
}

package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/DanielPopoola/cencori-go"
)

func main() {
	client, err := cencori.NewClient(
		cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 1: Single text embedding
	fmt.Println("Single Text Embedding")
	fmt.Println(repeat("=", 60))

	resp, err := client.Chat.Embeddings(context.Background(), cencori.EmbeddingParams{
		Input: "Hello, world!",
		Model: "text-embedding-3-small",
	})

	if err != nil {
		log.Fatalf("Embedding failed: %v", err)
	}

	fmt.Printf("Text: \"Hello, world!\"\n")
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Embedding dimension: %d\n", len(resp.Data[0].Embedding))
	fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
	fmt.Printf("First 5 values: %.4f, %.4f, %.4f, %.4f, %.4f\n",
		resp.Data[0].Embedding[0],
		resp.Data[0].Embedding[1],
		resp.Data[0].Embedding[2],
		resp.Data[0].Embedding[3],
		resp.Data[0].Embedding[4],
	)

	// Example 2: Batch embeddings
	fmt.Println("\n Batch Embeddings")
	fmt.Println(repeat("=", 60))

	texts := []string{
		"The cat sits on the mat",
		"A feline rests on a rug",
		"Dogs bark loudly at night",
	}

	batchResp, err := client.Chat.Embeddings(context.Background(), cencori.EmbeddingParams{
		Input: texts,
		Model: "text-embedding-3-small",
	})

	if err != nil {
		log.Fatalf("Batch embedding failed: %v", err)
	}

	fmt.Printf("Embedded %d texts\n", len(batchResp.Data))
	fmt.Printf("Total tokens: %d\n\n", batchResp.Usage.TotalTokens)

	for i, text := range texts {
		fmt.Printf("%d. \"%s\"\n", i+1, text)
		fmt.Printf("   Dimension: %d\n", len(batchResp.Data[i].Embedding))
	}

	// Example 3: Similarity comparison
	fmt.Println("\n Semantic Similarity")
	fmt.Println(repeat("=", 60))

	// Calculate cosine similarity between first two texts (similar meaning)
	sim1 := cosineSimilarity(
		batchResp.Data[0].Embedding,
		batchResp.Data[1].Embedding,
	)

	// Calculate similarity between first and third (different meaning)
	sim2 := cosineSimilarity(
		batchResp.Data[0].Embedding,
		batchResp.Data[2].Embedding,
	)

	fmt.Printf("\nSimilarity: \"%s\" vs \"%s\"\n", texts[0], texts[1])
	fmt.Printf("  Score: %.4f (similar meaning)\n", sim1)

	fmt.Printf("\nSimilarity: \"%s\" vs \"%s\"\n", texts[0], texts[2])
	fmt.Printf("  Score: %.4f (different meaning)\n", sim2)

	if sim1 > sim2 {
		fmt.Println("\n Similar sentences have higher similarity score!")
	}

	// Example 4: Search/retrieval use case
	fmt.Println("\n Semantic Search Demo")
	fmt.Println(repeat("=", 60))

	documents := []string{
		"Python is a high-level programming language",
		"Go is known for its concurrency features",
		"JavaScript runs in web browsers",
		"Coffee is a popular morning beverage",
	}

	query := "Tell me about programming languages"

	// Embed all documents + query
	allTexts := append([]string{query}, documents...)
	searchResp, err := client.Chat.Embeddings(context.Background(), cencori.EmbeddingParams{
		Input: allTexts,
		Model: "text-embedding-3-small",
	})

	if err != nil {
		log.Fatalf("Search embedding failed: %v", err)
	}

	queryEmbed := searchResp.Data[0].Embedding
	fmt.Printf("Query: \"%s\"\n\n", query)

	// Calculate similarities
	type result struct {
		text  string
		score float64
	}

	results := make([]result, 0, len(documents))
	for i, doc := range documents {
		score := cosineSimilarity(queryEmbed, searchResp.Data[i+1].Embedding)
		results = append(results, result{text: doc, score: score})
	}

	// Simple sort by score (descending)
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	fmt.Println("Ranked results:")
	for i, r := range results {
		fmt.Printf("%d. [%.4f] %s\n", i+1, r.score, r.text)
	}
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
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

func repeat(s string, count int) string {
	var result strings.Builder
	for range count {
		result.WriteString(s)
	}
	return result.String()
}

package main

import (
	"context"
	"fmt"
	"log"
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

	// Simple text completion (no conversation history)
	fmt.Println("Text Completion Example")
	fmt.Println("=" + repeat("=", 50))

	temp := 0.7
	maxTokens := 100

	resp, err := client.Chat.Completions(context.Background(), cencori.CompletionParams{
		Prompt:      "Write a haiku about coding in Go",
		Model:       "gpt-4o",
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	})

	if err != nil {
		log.Fatalf("Completion failed: %v", err)
	}

	fmt.Printf("\nPrompt: Write a haiku about coding in Go\n")
	fmt.Printf("\nResponse:\n%s\n", resp.Choices[0].Message.Content)
	fmt.Printf("\nTokens used: %d\n", resp.Usage.TotalTokens)

	// Compare with different models
	fmt.Println("\n\n Comparing Models")
	fmt.Println("=" + repeat("=", 50))

	prompt := "Explain recursion in one sentence."
	models := []string{"gpt-4o-mini", "gemini-2.5-flash", "claude-3-haiku"}

	for _, model := range models {
		resp, err := client.Chat.Completions(context.Background(), cencori.CompletionParams{
			Prompt: prompt,
			Model:  model,
		})

		if err != nil {
			fmt.Printf("\n %s: %v\n", model, err)
			continue
		}

		fmt.Printf("\n%s:\n", model)
		fmt.Printf("  %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("  Tokens: %d\n", resp.Usage.TotalTokens)
	}
}

func repeat(s string, count int) string {
	var result strings.Builder
	for range count {
		result.WriteString(s)
	}
	return result.String()
}

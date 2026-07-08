// Documents API example.
//
// Run with:
//
//	CENCORI_API_KEY=csk_... go run examples/13-documents/main.go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/cencori/cencori-go"
)

func main() {
	client, err := cencori.NewClient(
		cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
	)
	if err != nil {
		log.Fatalf("client: %v", err)
	}

	ctx := context.Background()
	pdfURL := "https://www.orimi.com/pdf-test.pdf"

	// 1. Extract — native PDF text (free, no LLM)
	fmt.Println("== extract ==")
	extracted, err := client.Documents.Extract(ctx, &cencori.DocumentParams{
		DocumentURL: pdfURL,
	})
	if err != nil {
		log.Fatalf("extract: %v", err)
	}
	fmt.Printf("method: %s\n", extracted.Method)
	fmt.Printf("pages:  %d\n", extracted.PageCount)
	fmt.Printf("chars:  %d\n", len(extracted.Text))
	if extracted.Cost != nil {
		fmt.Printf("cost:   $%.6f\n", extracted.Cost.CencoriChargeUsd)
	} else {
		fmt.Println("cost:   $0 (native extraction, no LLM call)")
	}

	// 2. Summarize
	fmt.Println("\n== summarize ==")
	summary, err := client.Documents.Summarize(ctx, &cencori.DocumentParams{DocumentURL: pdfURL})
	if err != nil {
		log.Fatalf("summarize: %v", err)
	}
	fmt.Println(summary.Summary)

	// 3. Query
	fmt.Println("\n== query ==")
	answer, err := client.Documents.Query(ctx, &cencori.DocumentParams{
		DocumentURL: pdfURL,
		Question:    "What is the main topic of this document?",
	})
	if err != nil {
		log.Fatalf("query: %v", err)
	}
	fmt.Println(answer.Answer)

	// 4. From a local file
	if data, err := os.ReadFile("sample.pdf"); err == nil {
		local, err := client.Documents.Extract(ctx, &cencori.DocumentParams{
			DocumentBase64: base64.StdEncoding.EncodeToString(data),
			MimeType:       "application/pdf",
		})
		if err != nil {
			log.Fatalf("local extract: %v", err)
		}
		fmt.Printf("\n== local file ==\nextracted %d chars via %s\n", len(local.Text), local.Method)
	}
}

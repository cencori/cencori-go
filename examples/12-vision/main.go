// Vision API example.
//
// Run with:
//
//	CENCORI_API_KEY=csk_... go run examples/12-vision/main.go
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
	imageURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3a/Cat03.jpg/1200px-Cat03.jpg"

	// 1. Analyze from a URL with a custom prompt
	analysis, err := client.Vision.Analyze(ctx, &cencori.VisionParams{
		ImageURL: imageURL,
		Prompt:   "Describe this animal and estimate its age.",
	})
	if err != nil {
		log.Fatalf("analyze: %v", err)
	}
	fmt.Println("== analyze ==")
	fmt.Println(analysis.Analysis)
	fmt.Printf("cost: $%.6f\n", analysis.Cost.CencoriChargeUsd)

	// 2. Describe (preset prompt)
	described, err := client.Vision.Describe(ctx, &cencori.VisionParams{ImageURL: imageURL})
	if err != nil {
		log.Fatalf("describe: %v", err)
	}
	fmt.Println("\n== describe ==")
	fmt.Println(described.Description)

	// 3. OCR from a local file (if present)
	if data, err := os.ReadFile("sample_receipt.png"); err == nil {
		ocr, err := client.Vision.OCR(ctx, &cencori.VisionParams{
			ImageBase64: base64.StdEncoding.EncodeToString(data),
			MimeType:    "image/png",
		})
		if err != nil {
			log.Fatalf("ocr: %v", err)
		}
		fmt.Println("\n== ocr ==")
		fmt.Println(ocr.Text)
	}

	// 4. Classify — structured JSON output
	classified, err := client.Vision.Classify(ctx, &cencori.VisionParams{ImageURL: imageURL})
	if err != nil {
		log.Fatalf("classify: %v", err)
	}
	fmt.Println("\n== classify ==")
	fmt.Printf("%#v\n", classified.Classification)
}

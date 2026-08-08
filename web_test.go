package cencori

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/web/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var request WebSearchParams
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Query != "agent evidence" {
			t.Fatalf("unexpected query: %s", request.Query)
		}
		_ = json.NewEncoder(w).Encode(WebSearchResponse{
			Query:        request.Query,
			Count:        1,
			SearchEngine: "cencori-web-lexical-v1",
			Results: []WebSearchHit{{
				ID:           "doc_1",
				Title:        "Evidence",
				CanonicalURL: "https://example.com/evidence",
				ContentHash:  "abc",
			}},
		})
	}))
	defer server.Close()

	client, err := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	response, err := client.Web.Search(context.Background(), &WebSearchParams{Query: "agent evidence"})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if response.Count != 1 || response.Results[0].Title != "Evidence" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestWebValidatesRequiredInputs(t *testing.T) {
	client, err := NewClient(WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}
	if _, err := client.Web.Fetch(context.Background(), &WebFetchParams{}); err == nil {
		t.Fatal("expected Fetch validation error")
	}
	if _, err := client.Web.Search(context.Background(), &WebSearchParams{}); err == nil {
		t.Fatal("expected Search validation error")
	}
	if _, err := client.Web.Crawl(context.Background(), &WebCrawlParams{}); err == nil {
		t.Fatal("expected Crawl validation error")
	}
}

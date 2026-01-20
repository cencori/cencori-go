# Cencori Go SDK Examples

Complete, production-ready examples for the Cencori Go SDK.

## Prerequisites

```bash
export CENCORI_API_KEY="your-key-here"
go get github.com/DanielPopoola/cencori-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/DanielPopoola/cencori-go"
)

func main() {
    client, _ := cencori.NewClient(
        cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
    )
    
    resp, err := client.Chat.Create(context.Background(), &cencori.ChatParams{
        Model: "gpt-4o",
        Messages: []cencori.Message{
            {Role: "user", Content: "Hello!"},
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## Examples by Category

### Chat API
- **[01-basic-chat](01-basic-chat/main.go)** - Simple chat completion
- **[02-streaming](02-streaming/main.go)** - Real-time streaming
- **[03-error-handling](03-error-handling/main.go)** - Robust error handling
- **[04-advanced-params](04-advanced-params/main.go)** - All parameters
- **[10-completions](10-completions/main.go)** - Text completion
- **[11-embeddings](11-embeddings/main.go)** - Vector embeddings & semantic search


### Projects Management
- **[05-projects](05-projects/main.go)** - Project management

### API Keys
- **[06-key-rotation](06-key-rotation/main.go)** - API key rotation

### Metrics & Analytics
- **[07-metrics](07-metrics/main.go)** - Analytics dashboard

### Advanced Patterns
- **[08-context-timeout](08-context-timeout/main.go)** - Context patterns
- **[09-multi-provider](09-multi-provider/main.go)** - Provider comparison

## Running Examples

```bash
# Run specific example
go run examples/01-basic-chat/main.go

# Run with race detector
go run -race examples/02-streaming/main.go

# Run all examples
for f in examples/**/*.go; do go run "$f"; done
```

## Common Patterns

### Error Handling

```go
import "errors"

_, err := client.Chat.Create(ctx, params)
if err != nil {
    switch {
    case errors.Is(err, cencori.ErrInvalidAPIKey):
        // Handle auth error
    case errors.Is(err, cencori.ErrRateLimited):
        // Implement backoff
    case errors.Is(err, context.DeadlineExceeded):
        // Handle timeout
    default:
        // Unknown error
    }
}
```

### Streaming

```go
stream, err := client.Chat.Stream(ctx, params)
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Err != nil {
        log.Fatal(chunk.Err)
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.Chat.Create(ctx, params)
// Request auto-cancelled after 10s
```

## Testing

Each example includes inline comments explaining:
- What it demonstrates
- Expected output
- Common gotchas
- Production considerations

## Need Help?

- 📚 [Full Documentation](https://cencori.com/docs)

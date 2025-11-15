## <!-- Based on: https://github.com/github/awesome-copilot/blob/main/instructions/go-mcp-server.instructions.md -->

description: 'Best practices and patterns for building Model Context Protocol (MCP) servers in Go using the official github.com/modelcontextprotocol/go-sdk package.'
applyTo: "**/\*.go, **/go.mod, \*\*/go.sum"

---

# Go MCP Server Development Guidelines

When building MCP servers in Go, follow these best practices and patterns using the official Go SDK.

## Server Setup

Create an MCP server using `mcp.NewServer`:

```go
import "github.com/modelcontextprotocol/go-sdk/mcp"

server := mcp.NewServer(
    &mcp.Implementation{
        Name:    "my-server",
        Version: "v1.0.0",
    },
    nil, // or provide mcp.Options
)
```

## Adding Tools

Use `mcp.AddTool` with struct-based input and output for type safety:

- Define input/output structs with JSON schema tags
- Use `jsonschema` tags to document fields
- Implement handler functions with proper signature
- Return errors for invalid inputs
- Check context cancellation

## Adding Resources

Use `mcp.AddResource` for providing accessible data:

- Define resource URIs and MIME types
- Implement resource read handlers
- Return appropriate content types
- Handle resource not found errors

## Adding Prompts

Use `mcp.AddPrompt` for reusable prompt templates:

- Define prompt arguments
- Return formatted prompt messages
- Use context for cancellation

## Transport Configuration

### Stdio Transport

For communication over stdin/stdout (most common for desktop integrations):

```go
if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
    log.Fatal(err)
}
```

### Streamable HTTP Transport

For HTTP-based communication using the current MCP specification:

```go
handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
    return server
}, &mcp.StreamableHTTPOptions{
    JSONResponse: true,
    Logger:       slog.Default(),
})

httpServer := &http.Server{
    Addr:              ":8080",
    Handler:           handler,
    ReadHeaderTimeout: 5 * time.Second,
}

if err := httpServer.ListenAndServe(); err != nil {
    log.Fatal(err)
}
```

## Error Handling

- Always return proper errors from handlers
- Use context for cancellation checking
- Validate inputs before processing
- Wrap errors with context using `fmt.Errorf("%w", err)`

## JSON Schema Tags

Use `jsonschema` tags to document your structs:

- `required` - Mark required fields
- `description` - Add field descriptions
- `minimum`, `maximum` - Set numeric bounds
- `format` - Specify formats (email, uri, etc.)
- `uniqueItems` - For arrays with unique items
- `default` - Specify default values

## Context Usage

Always respect context cancellation and deadlines:

- Check `ctx.Err()` in long-running operations
- Pass context to downstream functions
- Use `context.WithTimeout` for operations with deadlines

## Server Options

Configure server behavior with options:

- Set capabilities (tools, resources, prompts)
- Enable resource subscriptions if needed
- Configure server metadata

## Testing

Test your MCP tools using standard Go testing patterns:

- Write unit tests for tool handlers
- Use table-driven tests
- Mock external dependencies
- Test error conditions
- Verify context cancellation

## Module Setup

Initialize your Go module properly:

```bash
go mod init github.com/yourusername/yourserver
go get github.com/modelcontextprotocol/go-sdk@latest
```

## Common Patterns

### Logging

Use structured logging with `log/slog`

### Configuration

Use environment variables or config files for settings

### Graceful Shutdown

Handle shutdown signals properly using `signal.Notify`

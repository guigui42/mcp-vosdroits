// Package main provides the entry point for the VosDroits MCP server.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guigui42/mcp-vosdroits/internal/config"
	"github.com/guigui42/mcp-vosdroits/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Load configuration
	cfg := config.Load()

	// Override version if set at build time
	if version != "dev" {
		cfg.ServerVersion = version
	}

	// Set up logging
	setupLogging(cfg.LogLevel)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		slog.Info("Shutting down gracefully...")
		cancel()
	}()

	// Create MCP server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    cfg.ServerName,
			Version: cfg.ServerVersion,
		},
		nil,
	)

	// Register tools
	if err := tools.RegisterTools(server, cfg); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}

	slog.Info("Starting MCP server",
		"name", cfg.ServerName,
		"version", cfg.ServerVersion,
	)

	// Run server with appropriate transport
	if err := runWithTransport(ctx, server, cfg); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// runWithTransport runs the server with the appropriate transport based on configuration.
// If HTTPPort is set, it starts an HTTP server with Streamable HTTP transport.
// Otherwise, it uses stdio transport for stdio-based communication.
func runWithTransport(ctx context.Context, server *mcp.Server, cfg *config.Config) error {
	if cfg.HTTPPort != "" {
		slog.Info("Using Streamable HTTP transport", "port", cfg.HTTPPort)
		
		// Create Streamable HTTP handler
		handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
			return server
		}, &mcp.StreamableHTTPOptions{
			JSONResponse: true,
			Logger:       slog.Default(),
		})
		
		// Create HTTP server
		httpServer := &http.Server{
			Addr:              ":" + cfg.HTTPPort,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		}
		
		// Start server in goroutine
		errCh := make(chan error, 1)
		go func() {
			slog.Info("HTTP server listening", "addr", httpServer.Addr)
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
		}()
		
		// Wait for context cancellation or error
		select {
		case <-ctx.Done():
			slog.Info("Shutting down HTTP server...")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return httpServer.Shutdown(shutdownCtx)
		case err := <-errCh:
			return fmt.Errorf("http server error: %w", err)
		}
	}
	
	slog.Info("Using stdio transport")
	return server.Run(ctx, &mcp.StdioTransport{})
}

func setupLogging(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
}

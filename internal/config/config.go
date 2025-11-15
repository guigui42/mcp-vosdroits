// Package config provides configuration management for the MCP server.
package config

import (
	"os"
	"time"
)

// Config holds the server configuration.
type Config struct {
	ServerName        string
	ServerVersion     string
	LogLevel          string
	HTTPTimeout       time.Duration
	HTTPPort          string
	AllowClientOrigin string
}

// Load returns a new Config loaded from environment variables.
func Load() *Config {
	return &Config{
		ServerName:        getEnv("SERVER_NAME", "vosdroits"),
		ServerVersion:     getEnv("SERVER_VERSION", "v1.0.0"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		HTTPTimeout:       getEnvDuration("HTTP_TIMEOUT", 30*time.Second),
		HTTPPort:          getEnv("HTTP_PORT", ""),
		AllowClientOrigin: getEnv("ALLOW_CLIENT_ORIGIN", "*"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

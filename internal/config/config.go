package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL    string
	RedisURL       string
	Port           string
	SlackToken     string
	SlackChannelID string
}

// NewConfig creates a new configuration instance
func NewConfig() *Config {
	return LoadConfig()
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		Port:           getEnv("PORT", "8080"),
		SlackToken:     getEnv("SLACK_TOKEN", ""),
		SlackChannelID: getEnv("SLACK_CHANNEL_ID", ""),
	}
}

// GetSlackParams returns Slack configuration parameters
func (c *Config) GetSlackParams() (string, string) {
	return c.SlackToken, c.SlackChannelID
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

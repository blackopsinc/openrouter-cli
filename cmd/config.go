package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppConfig represents the application configuration
type AppConfig struct {
	DefaultModel string            `json:"default_model"`
	Timeout      int               `json:"timeout_seconds"`
	MaxFileSize  int64             `json:"max_file_size_bytes"`
	Models       map[string]string `json:"models"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		DefaultModel: "anthropic/claude-3.7-sonnet:thinking",
		Timeout:      60,
		MaxFileSize:  10 * 1024 * 1024, // 10MB
		Models: map[string]string{
			"claude-3-opus":     "anthropic/claude-3-opus",
			"claude-3-sonnet":   "anthropic/claude-3-sonnet",
			"claude-3-haiku":    "anthropic/claude-3-haiku",
			"gpt-4":             "openai/gpt-4",
			"gpt-4-turbo":       "openai/gpt-4-turbo",
			"gpt-3.5-turbo":     "openai/gpt-3.5-turbo",
			"claude-thinking":   "anthropic/claude-3.7-sonnet:thinking",
		},
	}
}

// LoadConfig loads configuration from file or returns default
func LoadConfig() (*AppConfig, error) {
	configPath := getConfigPath()
	
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return &config, nil
}

// SaveConfig saves configuration to file
func (c *AppConfig) SaveConfig() error {
	configPath := getConfigPath()
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// ResolveModel resolves a model alias to its full name
func (c *AppConfig) ResolveModel(model string) string {
	if fullModel, exists := c.Models[model]; exists {
		return fullModel
	}
	return model
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".openrouter-cli.json"
	}
	return filepath.Join(homeDir, ".config", "openrouter-cli", "config.json")
}

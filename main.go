package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"openrouter-cli/cmd"
)

const (
	envAPIKey = "OPENROUTER_API_KEY"
)

// Config holds the application configuration
type Config struct {
	APIKey       string
	FilePath     string
	Model        string
	DirectPrompt string
	Verbose      bool
	Timeout      time.Duration
	ListModels   bool
	SaveConfig   bool
}

func main() {
	// Load application config
	appConfig, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	config, err := parseFlags(appConfig)
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Handle special commands
	if config.ListModels {
		listModels(appConfig)
		return
	}

	if config.SaveConfig {
		if err := appConfig.SaveConfig(); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
		fmt.Println("Configuration saved successfully")
		return
	}

	// Resolve model alias
	resolvedModel := appConfig.ResolveModel(config.Model)
	if config.Verbose {
		if resolvedModel != config.Model {
			log.Printf("Resolved model alias '%s' to '%s'", config.Model, resolvedModel)
		}
		log.Printf("Using model: %s", resolvedModel)
		if config.FilePath != "" {
			log.Printf("Reading from file: %s", config.FilePath)
		}
		log.Printf("Request timeout: %v", config.Timeout)
	}

	// Create client with timeout
	client := cmd.NewOpenRouterClient(config.APIKey)
	client.SetTimeout(config.Timeout)

	response, err := client.ProcessInput(config.FilePath, config.DirectPrompt, resolvedModel)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Println(response)
}

// parseFlags parses command line flags and returns configuration
func parseFlags(appConfig *cmd.AppConfig) (*Config, error) {
	config := &Config{}

	// Define flags
	flag.StringVar(&config.FilePath, "file", "", "Path to a file whose content will be used as the prompt")
	flag.StringVar(&config.APIKey, "key", "", "OpenRouter API key (can also be set via OPENROUTER_API_KEY env var)")
	flag.StringVar(&config.Model, "model", appConfig.DefaultModel, "Model to use for the request (or model alias)")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.Verbose, "v", false, "Enable verbose output (shorthand)")
	flag.DurationVar(&config.Timeout, "timeout", time.Duration(appConfig.Timeout)*time.Second, "Request timeout")
	flag.BoolVar(&config.ListModels, "list-models", false, "List available model aliases")
	flag.BoolVar(&config.SaveConfig, "save-config", false, "Save current configuration to file")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [PROMPT]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A CLI tool for sending prompts to OpenRouter's API.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nModel Aliases:\n")
		for alias, fullName := range appConfig.Models {
			fmt.Fprintf(os.Stderr, "  %-15s -> %s\n", alias, fullName)
		}
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -key YOUR_KEY \"What is Go?\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -file prompt.txt -model gpt-4\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --list-models\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  export OPENROUTER_API_KEY=your_key && %s \"Explain AI\"\n", os.Args[0])
	}

	flag.Parse()

	// Get API key from environment if not provided via flag
	if config.APIKey == "" {
		config.APIKey = os.Getenv(envAPIKey)
		if config.APIKey == "" && !config.ListModels && !config.SaveConfig {
			return nil, fmt.Errorf("API key is required. Provide it via -key flag or %s environment variable", envAPIKey)
		}
	}

	// Get direct prompt from remaining arguments
	args := flag.Args()
	if len(args) > 0 {
		config.DirectPrompt = strings.Join(args, " ")
	}

	// Validate configuration
	if config.FilePath == "" && config.DirectPrompt == "" && !config.ListModels && !config.SaveConfig {
		if config.Verbose {
			log.Println("No file or direct prompt provided, will read from stdin")
		}
	}

	return config, nil
}

// listModels displays available model aliases
func listModels(appConfig *cmd.AppConfig) {
	fmt.Println("Available model aliases:")
	fmt.Println()
	
	maxAliasLen := 0
	for alias := range appConfig.Models {
		if len(alias) > maxAliasLen {
			maxAliasLen = len(alias)
		}
	}
	
	for alias, fullName := range appConfig.Models {
		fmt.Printf("  %-*s -> %s\n", maxAliasLen, alias, fullName)
	}
	
	fmt.Printf("\nDefault model: %s\n", appConfig.DefaultModel)
}

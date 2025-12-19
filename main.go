package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"openrouter-cli/cmd"
)

const (
	envAPIKey    = "OPENROUTER_API_KEY"
	envModel     = "OPENROUTER_MODEL"
	envPrePrompt = "OPENROUTER_PRE_PROMPT"
)

func main() {
	// Get API key from environment (required)
	apiKey := os.Getenv(envAPIKey)
	if apiKey == "" {
		log.Fatalf("API key is required. Set %s environment variable", envAPIKey)
	}

	// Get model from environment or use default
	model := getModel()

	// Read and prepare input from stdin
	input, err := prepareInput()
	if err != nil {
		log.Fatalf("Failed to prepare input: %v", err)
	}

	// Send request to OpenRouter
	client := cmd.NewOpenRouterClient(apiKey)
	response, err := client.ProcessInput("", input, model)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Println(response)
}

// getModel returns the model to use, from environment or config default
func getModel() string {
	appConfig, err := cmd.LoadConfig()
	if err != nil {
		appConfig = cmd.DefaultConfig()
	}

	model := os.Getenv(envModel)
	if model == "" {
		return appConfig.DefaultModel
	}

	// Resolve model alias if it's an alias
	return appConfig.ResolveModel(model)
}

// prepareInput reads from stdin and optionally prepends a pre-prompt
func prepareInput() (string, error) {
	// Read input from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	input := strings.TrimSpace(string(data))
	if input == "" {
		return "", fmt.Errorf("input is empty")
	}

	// Prepend pre-prompt if set
	if prePrompt := os.Getenv(envPrePrompt); prePrompt != "" {
		input = prePrompt + "\n\n" + input
	}

	return input, nil
}

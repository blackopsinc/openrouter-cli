package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"openrouter-cli/cmd"
)

func main() {
	// Define flags
	filePath := flag.String("file", "", "Path to a file whose content will be used as the prompt")
	apiKey := flag.String("key", "", "OpenRouter API key (required)")
	modelName := flag.String("model", "anthropic/claude-3.7-sonnet:thinking", "Model to use for the request")
	
	// Parse flags
	flag.Parse()

	// Check if API key is provided via flag or environment variable
	if *apiKey == "" {
		// Try to get API key from environment variable
		envApiKey := os.Getenv("OPENROUTER_API_KEY")
		if envApiKey != "" {
			*apiKey = envApiKey
		} else {
			fmt.Println("Error: OpenRouter API key is required")
			fmt.Println("You can provide it with the -key flag or by setting the OPENROUTER_API_KEY environment variable")
			flag.Usage()
			os.Exit(1)
		}
	}

	// Get direct prompt from arguments if provided
	args := flag.Args()
	var directPrompt string
	if len(args) > 0 {
		directPrompt = strings.Join(args, " ")
	}

	// Create a new client
	client := cmd.NewOpenRouterClient(*apiKey)

	// Process input and make request
	response, err := client.ProcessInput(*filePath, directPrompt, *modelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print response
	fmt.Println(response)
}

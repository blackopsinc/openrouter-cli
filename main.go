package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	envAPIKey    = "OPENROUTER_API_KEY"
	envModel     = "OPENROUTER_MODEL"
	envPrePrompt = "OPENROUTER_PRE_PROMPT"

	openRouterURL  = "https://openrouter.ai/api/v1/chat/completions"
	defaultTimeout = 60 * time.Second
	defaultModel  = "openai/gpt-oss-20b:free"
	userAgent      = "OpenRouter-CLI/1.0"
)

// OpenRouterRequest represents the request body for the OpenRouter API
type OpenRouterRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from the OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func main() {
	// Get API key from environment (required)
	apiKey := os.Getenv(envAPIKey)
	if apiKey == "" {
		log.Fatalf("API key is required. Set %s environment variable", envAPIKey)
	}

	// Get model from environment or use default
	model := os.Getenv(envModel)
	if model == "" {
		model = defaultModel
	}

	// Read and prepare input from stdin
	input, err := prepareInput()
	if err != nil {
		log.Fatalf("Failed to prepare input: %v", err)
	}

	// Send request to OpenRouter
	response, err := sendRequest(apiKey, input, model)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	fmt.Println(response)
}

// prepareInput reads from stdin and optionally prepends a pre-prompt
func prepareInput() (string, error) {
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

// sendRequest sends a request to the OpenRouter API
func sendRequest(apiKey, input, modelName string) (string, error) {
	reqBody := OpenRouterRequest{
		Model: modelName,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: input,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("HTTP-Referer", "https://github.com/blackopsinc/openrouter-cli")
	req.Header.Set("X-Title", "OpenRouter CLI")

	client := &http.Client{Timeout: defaultTimeout}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timed out after %v", defaultTimeout)
		}
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openRouterResp.Error != nil {
		return "", fmt.Errorf("API error (%s): %s",
			openRouterResp.Error.Type, openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no response received from the API")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}

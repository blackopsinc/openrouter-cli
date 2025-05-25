package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	openRouterURL = "https://openrouter.ai/api/v1/chat/completions"
)

// OpenRouterClient handles interactions with the OpenRouter API
type OpenRouterClient struct {
	apiKey string
}

// NewOpenRouterClient creates a new OpenRouter API client
func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{apiKey: apiKey}
}

// ProcessInput processes input from various sources and sends it to OpenRouter
func (c *OpenRouterClient) ProcessInput(filePath string, directPrompt string, modelName string) (string, error) {
	var input string
	var err error

	// Priority for input sources
	if filePath != "" {
		// Prioritize file content - upload file as prompt
		fileContent, err := readFromFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read from file: %w", err)
		}
		input = fileContent
	} else if directPrompt != "" {
		// Use direct prompt if no file is provided
		input = directPrompt
	} else {
		// Fall back to stdin if neither file nor direct prompt is provided
		input, err = readFromStdin()
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
	}

	// Trim input
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("input is empty")
	}

	// Send request to OpenRouter
	return c.sendRequest(input, modelName)
}

// readFromFile reads text from a file
func readFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// readFromStdin reads text from standard input
func readFromStdin() (string, error) {
	fmt.Println("Enter your message (press Ctrl+D on Unix/Linux or Ctrl+Z followed by Enter on Windows when done):")
	
	scanner := bufio.NewScanner(os.Stdin)
	var inputBuilder strings.Builder
	
	for scanner.Scan() {
		inputBuilder.WriteString(scanner.Text() + "\n")
	}
	
	if err := scanner.Err(); err != nil {
		return "", err
	}
	
	return inputBuilder.String(), nil
}

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
	} `json:"error,omitempty"`
}

// sendRequest sends a request to the OpenRouter API
func (c *OpenRouterClient) sendRequest(input string, modelName string) (string, error) {
	// Create request body
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
	
	// Create HTTP request
	req, err := http.NewRequest("POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/openrouter-cli")
	req.Header.Set("X-Title", "OpenRouter CLI")
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	// Parse response
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Check for errors
	if openRouterResp.Error != nil && openRouterResp.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", openRouterResp.Error.Message)
	}
	
	// Check if we have any choices
	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no response received from the API")
	}
	
	// Return the content of the first choice
	return openRouterResp.Choices[0].Message.Content, nil
}

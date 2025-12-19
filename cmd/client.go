package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	openRouterURL    = "https://openrouter.ai/api/v1/chat/completions"
	defaultTimeout   = 60 * time.Second
	maxFileSize      = 10 * 1024 * 1024 // 10MB
	userAgent        = "OpenRouter-CLI/1.0"
)

// OpenRouterClient handles interactions with the OpenRouter API
type OpenRouterClient struct {
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
}

// NewOpenRouterClient creates a new OpenRouter API client
func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:  apiKey,
		timeout: defaultTimeout,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// SetTimeout sets the request timeout
func (c *OpenRouterClient) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	c.httpClient.Timeout = timeout
}

// ProcessInput processes input from various sources and sends it to OpenRouter
func (c *OpenRouterClient) ProcessInput(filePath string, directPrompt string, modelName string) (string, error) {
	var input string
	var err error

	// Priority for input sources
	switch {
	case filePath != "":
		// Prioritize file content - upload file as prompt
		input, err = c.readFromFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read from file '%s': %w", filePath, err)
		}
	case directPrompt != "":
		// Use direct prompt if no file is provided
		input = directPrompt
	default:
		// Fall back to stdin if neither file nor direct prompt is provided
		input, err = c.readFromStdin()
		if err != nil {
			return "", fmt.Errorf("failed to read from stdin: %w", err)
		}
	}

	// Validate and clean input
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("input is empty")
	}

	// Send request to OpenRouter
	return c.sendRequest(input, modelName)
}

// readFromFile reads text from a file with size validation
func (c *OpenRouterClient) readFromFile(filePath string) (string, error) {
	// Check file size first
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot access file: %w", err)
	}

	if fileInfo.Size() > maxFileSize {
		return "", fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", 
			fileInfo.Size(), maxFileSize)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot read file: %w", err)
	}

	return string(content), nil
}

// readFromStdin reads text from standard input (for piped data)
func (c *OpenRouterClient) readFromStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var inputBuilder strings.Builder
	
	for scanner.Scan() {
		inputBuilder.WriteString(scanner.Text())
		inputBuilder.WriteString("\n")
	}
	
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}
	
	return inputBuilder.String(), nil
}

// OpenRouterRequest represents the request body for the OpenRouter API
type OpenRouterRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from the OpenRouter API
type OpenRouterResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// sendRequest sends a request to the OpenRouter API with improved error handling
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
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", openRouterURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("HTTP-Referer", "https://github.com/openrouter-cli")
	req.Header.Set("X-Title", "OpenRouter CLI")
	
	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timed out after %v", c.timeout)
		}
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	// Handle HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	
	// Check for API errors
	if openRouterResp.Error != nil {
		return "", fmt.Errorf("API error (%s): %s", 
			openRouterResp.Error.Type, openRouterResp.Error.Message)
	}
	
	// Check if we have any choices
	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no response received from the API")
	}
	
	// Return the content of the first choice
	return openRouterResp.Choices[0].Message.Content, nil
}

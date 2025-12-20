package main

import (
	"bufio"
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
	envProvider  = "LLM_PROVIDER"
	envModel     = "LLM_MODEL"
	envPrePrompt = "LLM_PRE_PROMPT"
	envStream    = "LLM_STREAM"
	envVerbose   = "LLM_VERBOSE"
	envOllamaURL = "OLLAMA_URL"
	envLMStudioURL = "LM_STUDIO_URL"

	defaultTimeout = 60 * time.Second
	defaultModel  = "openai/gpt-oss-20b:free"
	userAgent      = "OpenRouter-CLI/1.0"

	// Provider URLs
	openRouterURL  = "https://openrouter.ai/api/v1/chat/completions"
	defaultOllamaURL = "http://localhost:11434/api/chat"
	defaultLMStudioURL = "http://10.8.0.22:1234/v1/chat/completions"
)

// Provider represents the LLM provider type
type Provider string

const (
	ProviderOpenRouter Provider = "openrouter"
	ProviderOllama     Provider = "ollama"
	ProviderLMStudio   Provider = "lmstudio"
)

// OpenRouterRequest represents the request body for the OpenRouter API
type OpenRouterRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from the OpenRouter API
type OpenRouterResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// OllamaResponse represents the response from Ollama's native API
type OllamaResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done   bool `json:"done"`
	Error  string `json:"error,omitempty"`
}

func main() {
	// Get provider from environment or default to openrouter
	providerStr := strings.ToLower(strings.TrimSpace(os.Getenv(envProvider)))
	if providerStr == "" {
		providerStr = string(ProviderOpenRouter)
	}

	provider := Provider(providerStr)
	if provider != ProviderOpenRouter && provider != ProviderOllama && provider != ProviderLMStudio {
		log.Fatalf("Invalid provider: %s. Must be one of: openrouter, ollama, lmstudio", providerStr)
	}

	// Get API key from environment (required for OpenRouter, optional for others)
	apiKey := strings.TrimSpace(os.Getenv(envAPIKey))
	if provider == ProviderOpenRouter && apiKey == "" {
		log.Fatalf("API key is required for OpenRouter. Set %s environment variable", envAPIKey)
	}

	// Get model from environment or use default
	model := strings.TrimSpace(os.Getenv(envModel))
	if model == "" {
		if provider == ProviderOllama {
			model = "llama2" // Default Ollama model
		} else if provider == ProviderLMStudio {
			model = "local-model" // Default LM Studio model
		} else {
			model = defaultModel
		}
	}

	// Get provider-specific URLs
	ollamaURL := strings.TrimSpace(os.Getenv(envOllamaURL))
	if ollamaURL == "" {
		ollamaURL = defaultOllamaURL
	}

	lmStudioURL := strings.TrimSpace(os.Getenv(envLMStudioURL))
	if lmStudioURL == "" {
		lmStudioURL = defaultLMStudioURL
	}

	// Check if streaming is enabled
	stream := isEnvSet(envStream)

	// Check if verbose mode is enabled
	verbose := isEnvSet(envVerbose)

	if verbose {
		log.Printf("[DEBUG] Starting LLM CLI")
		log.Printf("[DEBUG] Provider: %s", provider)
		log.Printf("[DEBUG] Model: %s", model)
		log.Printf("[DEBUG] Streaming: %v", stream)
		if provider == ProviderOllama {
			log.Printf("[DEBUG] Ollama URL: %s", ollamaURL)
		}
		if provider == ProviderLMStudio {
			log.Printf("[DEBUG] LM Studio URL: %s", lmStudioURL)
		}
	}

	// Read and prepare input from stdin
	input, err := prepareInput(verbose)
	if err != nil {
		log.Fatalf("Failed to prepare input: %v", err)
	}

	if verbose {
		log.Printf("[DEBUG] Input length: %d characters", len(input))
	}

	// Send request based on provider
	if stream {
		err = sendStreamingRequest(provider, apiKey, input, model, ollamaURL, lmStudioURL, verbose)
	} else {
		response, err := sendRequest(provider, apiKey, input, model, ollamaURL, lmStudioURL, verbose)
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}
		fmt.Println(response)
	}

	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
}

// isEnvSet checks if an environment variable is set to a truthy value
func isEnvSet(key string) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return val == "1" || val == "true" || val == "yes" || val == "on"
}

// prepareInput reads from stdin and optionally prepends a pre-prompt
func prepareInput(verbose bool) (string, error) {
	if verbose {
		log.Printf("[DEBUG] Reading input from stdin...")
	}

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
		if verbose {
			log.Printf("[DEBUG] Prepending pre-prompt (length: %d)", len(prePrompt))
		}
		input = prePrompt + "\n\n" + input
	}

	return input, nil
}

// getAPIURL returns the appropriate API URL based on provider
func getAPIURL(provider Provider, ollamaURL, lmStudioURL string) string {
	switch provider {
	case ProviderOllama:
		return ollamaURL
	case ProviderLMStudio:
		return lmStudioURL
	default:
		return openRouterURL
	}
}

// sendRequest sends a request to the LLM API (non-streaming)
func sendRequest(provider Provider, apiKey, input, modelName, ollamaURL, lmStudioURL string, verbose bool) (string, error) {
	apiURL := getAPIURL(provider, ollamaURL, lmStudioURL)

	// Ollama uses a slightly different request format
	var reqBody interface{}
	if provider == ProviderOllama {
		reqBody = map[string]interface{}{
			"model": modelName,
			"messages": []ChatMessage{
				{
					Role:    "user",
					Content: input,
				},
			},
			"stream": false,
		}
	} else {
		reqBody = OpenRouterRequest{
			Model: modelName,
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: input,
				},
			},
			Stream: false,
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	if verbose {
		log.Printf("[DEBUG] Request URL: %s", apiURL)
		log.Printf("[DEBUG] Request body size: %d bytes", len(jsonData))
		log.Printf("[DEBUG] Request model: %s", modelName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	// Only set Authorization header for OpenRouter
	if provider == ProviderOpenRouter && apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Referer", "https://github.com/blackopsinc/openrouter-cli")
		req.Header.Set("X-Title", "OpenRouter CLI")
	}

	if verbose {
		log.Printf("[DEBUG] Sending HTTP POST request...")
	}

	client := &http.Client{Timeout: defaultTimeout}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("request timed out after %v", defaultTimeout)
		}
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if verbose {
		log.Printf("[DEBUG] Response status: %d %s", resp.StatusCode, resp.Status)
		log.Printf("[DEBUG] Response headers: %v", resp.Header)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if verbose {
		log.Printf("[DEBUG] Response body size: %d bytes", len(body))
	}

	// Try to parse error response for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		if provider == ProviderOllama {
			var ollamaResp OllamaResponse
			if err := json.Unmarshal(body, &ollamaResp); err == nil && ollamaResp.Error != "" {
				return "", fmt.Errorf("HTTP %d - Ollama error: %s", resp.StatusCode, ollamaResp.Error)
			}
		} else {
			var openRouterResp OpenRouterResponse
			if err := json.Unmarshal(body, &openRouterResp); err == nil && openRouterResp.Error != nil {
				return "", fmt.Errorf("HTTP %d - API error (%s): %s",
					resp.StatusCode, openRouterResp.Error.Type, openRouterResp.Error.Message)
			}
		}
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Handle Ollama's native response format
	if provider == ProviderOllama {
		var ollamaResp OllamaResponse
		if err := json.Unmarshal(body, &ollamaResp); err != nil {
			if verbose {
				log.Printf("[DEBUG] Failed to parse Ollama JSON response: %v", err)
				log.Printf("[DEBUG] Response body (first 500 chars): %s", string(body[:min(500, len(body))]))
			}
			return "", fmt.Errorf("failed to parse Ollama response: %w", err)
		}

		if ollamaResp.Error != "" {
			return "", fmt.Errorf("Ollama error: %s", ollamaResp.Error)
		}

		if verbose {
			log.Printf("[DEBUG] Successfully received Ollama response")
		}

		return ollamaResp.Message.Content, nil
	}

	// Handle OpenAI-compatible response format (OpenRouter, LM Studio)
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		if verbose {
			log.Printf("[DEBUG] Failed to parse JSON response: %v", err)
			log.Printf("[DEBUG] Response body (first 500 chars): %s", string(body[:min(500, len(body))]))
		}
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openRouterResp.Error != nil {
		return "", fmt.Errorf("API error (%s): %s",
			openRouterResp.Error.Type, openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no response received from the API")
	}

	if verbose {
		log.Printf("[DEBUG] Successfully received response with %d choice(s)", len(openRouterResp.Choices))
	}

	return openRouterResp.Choices[0].Message.Content, nil
}

// sendStreamingRequest sends a streaming request to the LLM API (SSE)
func sendStreamingRequest(provider Provider, apiKey, input, modelName, ollamaURL, lmStudioURL string, verbose bool) error {
	apiURL := getAPIURL(provider, ollamaURL, lmStudioURL)

	// Ollama uses a slightly different request format
	var reqBody interface{}
	if provider == ProviderOllama {
		reqBody = map[string]interface{}{
			"model": modelName,
			"messages": []ChatMessage{
				{
					Role:    "user",
					Content: input,
				},
			},
			"stream": true,
		}
	} else {
		reqBody = OpenRouterRequest{
			Model: modelName,
			Messages: []ChatMessage{
				{
					Role:    "user",
					Content: input,
				},
			},
			Stream: true,
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if verbose {
		log.Printf("[DEBUG] Streaming request URL: %s", apiURL)
		log.Printf("[DEBUG] Request body size: %d bytes", len(jsonData))
		log.Printf("[DEBUG] Request model: %s", modelName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	// Only set Authorization header for OpenRouter
	if provider == ProviderOpenRouter && apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Referer", "https://github.com/blackopsinc/openrouter-cli")
		req.Header.Set("X-Title", "OpenRouter CLI")
	}

	if verbose {
		log.Printf("[DEBUG] Sending streaming HTTP POST request...")
	}

	client := &http.Client{Timeout: defaultTimeout}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("request timed out after %v", defaultTimeout)
		}
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if verbose {
		log.Printf("[DEBUG] Response status: %d %s", resp.StatusCode, resp.Status)
		log.Printf("[DEBUG] Content-Type: %s", resp.Header.Get("Content-Type"))
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if provider == ProviderOllama {
			var ollamaResp OllamaResponse
			if err := json.Unmarshal(body, &ollamaResp); err == nil && ollamaResp.Error != "" {
				return fmt.Errorf("HTTP %d - Ollama error: %s", resp.StatusCode, ollamaResp.Error)
			}
		} else {
			var openRouterResp OpenRouterResponse
			if err := json.Unmarshal(body, &openRouterResp); err == nil && openRouterResp.Error != nil {
				return fmt.Errorf("HTTP %d - API error (%s): %s",
					resp.StatusCode, openRouterResp.Error.Type, openRouterResp.Error.Message)
			}
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Handle Ollama's native streaming format (newline-delimited JSON, not SSE)
	if provider == ProviderOllama {
		scanner := bufio.NewScanner(resp.Body)
		var fullContent strings.Builder
		chunkCount := 0

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			if verbose {
				log.Printf("[DEBUG] Ollama stream line: %s", line)
			}

			var chunk OllamaResponse
			if err := json.Unmarshal([]byte(line), &chunk); err != nil {
				if verbose {
					log.Printf("[DEBUG] Failed to parse Ollama chunk: %v", err)
					log.Printf("[DEBUG] Chunk data: %s", line)
				}
				continue
			}

			if chunk.Error != "" {
				return fmt.Errorf("Ollama error in stream: %s", chunk.Error)
			}

			if chunk.Message.Content != "" {
				fmt.Print(chunk.Message.Content)
				fullContent.WriteString(chunk.Message.Content)
				chunkCount++
			}

			if chunk.Done {
				if verbose {
					log.Printf("[DEBUG] Ollama stream finished")
				}
				break
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to read Ollama stream: %w", err)
		}

		if verbose {
			log.Printf("[DEBUG] Ollama stream complete. Received %d chunks, total length: %d characters",
				chunkCount, fullContent.Len())
		}

		fmt.Println()
		return nil
	}

	// Parse SSE stream for OpenAI-compatible providers (OpenRouter, LM Studio)
	scanner := bufio.NewScanner(resp.Body)
	var fullContent strings.Builder
	chunkCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		if verbose {
			log.Printf("[DEBUG] SSE line: %s", line)
		}

		// Skip empty lines and non-data lines
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Extract JSON data
		data := strings.TrimPrefix(line, "data: ")

		// Check for [DONE] marker
		if data == "[DONE]" {
			if verbose {
				log.Printf("[DEBUG] Received [DONE] marker, stream complete")
			}
			break
		}

		// Parse JSON chunk
		var chunk OpenRouterResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			if verbose {
				log.Printf("[DEBUG] Failed to parse SSE chunk: %v", err)
				log.Printf("[DEBUG] Chunk data: %s", data)
			}
			continue
		}

		// Check for errors in chunk
		if chunk.Error != nil {
			return fmt.Errorf("API error in stream (%s): %s",
				chunk.Error.Type, chunk.Error.Message)
		}

		// Extract content from delta (streaming) or message (final)
		if len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			var content string

			// Streaming responses use delta, final responses use message
			if choice.Delta.Content != "" {
				content = choice.Delta.Content
			} else if choice.Message.Content != "" {
				content = choice.Message.Content
			}

			if content != "" {
				fmt.Print(content)
				fullContent.WriteString(content)
				chunkCount++
			}

			// Check for finish reason
			if choice.FinishReason != "" {
				if verbose {
					log.Printf("[DEBUG] Stream finished with reason: %s", choice.FinishReason)
				}
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read stream: %w", err)
	}

	if verbose {
		log.Printf("[DEBUG] Stream complete. Received %d chunks, total length: %d characters",
			chunkCount, fullContent.Len())
	}

	// Print newline after stream completes
	fmt.Println()

	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


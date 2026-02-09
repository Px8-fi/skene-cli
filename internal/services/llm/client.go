package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Provider endpoints
var providerEndpoints = map[string]string{
	"openai":    "https://api.openai.com/v1/chat/completions",
	"anthropic": "https://api.anthropic.com/v1/messages",
	"gemini":    "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions",
	"ollama":    "http://localhost:11434/v1/chat/completions",
	"lmstudio":  "http://localhost:1234/v1/chat/completions",
}

// Client handles LLM API calls
type Client struct {
	provider string
	model    string
	apiKey   string
	baseURL  string
	client   *http.Client
}

// NewClient creates a new LLM client
func NewClient(provider, model, apiKey, baseURL string) *Client {
	return &Client{
		provider: provider,
		model:    model,
		apiKey:   apiKey,
		baseURL:  baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	// Anthropic-specific fields
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content,omitempty"`
}

// AnthropicRequest represents an Anthropic API request
type AnthropicRequest struct {
	Model     string        `json:"model"`
	System    string        `json:"system,omitempty"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens"`
}

// Chat sends a chat completion request and returns the response text
func (c *Client) Chat(ctx context.Context, messages []ChatMessage, maxTokens int) (string, error) {
	if c.provider == "anthropic" {
		return c.chatAnthropic(ctx, messages, maxTokens)
	}
	return c.chatOpenAICompatible(ctx, messages, maxTokens)
}

// chatOpenAICompatible handles OpenAI-compatible APIs (OpenAI, Gemini, Ollama, LM Studio, generic)
func (c *Client) chatOpenAICompatible(ctx context.Context, messages []ChatMessage, maxTokens int) (string, error) {
	endpoint := c.getEndpoint()

	req := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Set auth header based on provider
	switch c.provider {
	case "gemini":
		// Gemini uses API key as query parameter
		q := httpReq.URL.Query()
		q.Add("key", c.apiKey)
		httpReq.URL.RawQuery = q.Encode()
	case "ollama", "lmstudio":
		// Local providers typically don't need auth
	default:
		if c.apiKey != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		}
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// chatAnthropic handles Anthropic's API format
func (c *Client) chatAnthropic(ctx context.Context, messages []ChatMessage, maxTokens int) (string, error) {
	endpoint := c.getEndpoint()

	if maxTokens == 0 {
		maxTokens = 4096
	}

	// Anthropic requires system prompt as a top-level field, not in messages
	var systemPrompt string
	var userMessages []ChatMessage
	for _, msg := range messages {
		if msg.Role == "system" {
			systemPrompt = msg.Content
		} else {
			userMessages = append(userMessages, msg)
		}
	}

	req := AnthropicRequest{
		Model:     c.model,
		System:    systemPrompt,
		Messages:  userMessages,
		MaxTokens: maxTokens,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Anthropic returns content array
	if len(chatResp.Content) > 0 {
		return chatResp.Content[0].Text, nil
	}

	// Fallback to OpenAI format
	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content in response")
}

// getEndpoint returns the API endpoint for the provider
func (c *Client) getEndpoint() string {
	if c.baseURL != "" {
		// Custom base URL - append chat completions path
		if c.provider == "anthropic" {
			return c.baseURL + "/v1/messages"
		}
		return c.baseURL + "/chat/completions"
	}

	if endpoint, ok := providerEndpoints[c.provider]; ok {
		return endpoint
	}

	// Default to OpenAI
	return providerEndpoints["openai"]
}

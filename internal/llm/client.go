package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient    *http.Client
	endpoint      string
	deployment    string
	apiVersion    string
	tokenProvider func(ctx context.Context, scope string) (string, error)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewClient(endpoint, deployment, apiVersion string, tokenProvider func(ctx context.Context, scope string) (string, error)) (*Client, error) {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		endpoint:      endpoint,
		deployment:    deployment,
		apiVersion:    apiVersion,
		tokenProvider: tokenProvider,
	}, nil
}

func (c *Client) GetCommand(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	scope := strings.TrimSuffix(c.endpoint, "/") + "/.default"
	token, err := c.tokenProvider(ctx, scope)
	if err != nil {
		return "", fmt.Errorf("getting token: %w", err)
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		c.endpoint, c.deployment, c.apiVersion)

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	body := ChatRequest{
		Messages:    messages,
		MaxTokens:   512,
		Temperature: 0.3,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return "", fmt.Errorf("API returned status %d (failed to read body: %v)", resp.StatusCode, readErr)
		}
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	command := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	command = strings.Trim(command, "`")
	command = strings.Trim(command, "```")
	command = strings.TrimSpace(command)

	return command, nil
}

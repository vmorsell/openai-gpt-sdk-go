package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultEndpoint = "https://api.openai.com/v1"
	defaultModel    = "gpt-3.5-turbo"
	jsonMIME        = "application/json"
)

// Config provides configuration to a client instance.
type Config struct {
	// The API key to use.
	APIKey string

	// The API endpoint to use for a client.
	Endpoint string
}

// NewConfig returns a pointer to a new initialized config.
func NewConfig() *Config {
	return &Config{
		Endpoint: defaultEndpoint,
	}
}

// WithAPIKey sets API key for a config.
func (c *Config) WithAPIKey(apiKey string) *Config {
	c.APIKey = apiKey
	return c
}

// WithEndpoint sets endpoint for a config.
func (c *Config) WithEndpoint(endpoint string) *Config {
	c.Endpoint = endpoint
	return c
}

// Client implements the request and response handling.
type Client struct {
	Config *Config
}

// NewClient returns a pointer to a new initialized client.
func NewClient(config *Config) *Client {
	return &Client{
		config,
	}
}

// makeCall makes a call to the OpenAI API.
func (c *Client) makeCall(httpPath string, payload interface{}, out interface{}) error {
	if payload == nil {
		return fmt.Errorf("empty payload")
	}

	if out == nil {
		return fmt.Errorf("missing return type")
	}

	url := strings.Join([]string{c.Config.Endpoint, httpPath}, "")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("encode payload: %w", err)
	}

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Add("Content-Type", jsonMIME)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.APIKey))

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return fmt.Errorf("http error %d: %s", res.StatusCode, body)
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

// ChatCompletion implements the chat completion API method.
func (c *Client) ChatCompletion(in ChatCompletionInput) (*ChatCompletionOutput, error) {
	if in.Model == "" {
		in.Model = defaultModel
	}

	path := "/chat/completions"
	out := ChatCompletionOutput{}

	if err := c.makeCall(path, in, &out); err != nil {
		return nil, fmt.Errorf("make call: %w", err)
	}

	return &out, nil
}

// ChatCompletionInput is the input to a ChatCompletion call.
type ChatCompletionInput struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatCompletionOutput is the output to a ChatCompletion call.
type ChatCompletionOutput struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice is a single returned suggestion of what the next chat
// message could be.
type Choice struct {
	Index        int `json:"index"`
	Message      Message
	FinishReason string `json:"finish_reason"`
}

// Message represents a message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage holds token usage reporting from ChatGPT.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
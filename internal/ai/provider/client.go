package provider

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Creates and owns the shared Gemini client.

type Client struct {
	client     *genai.Client
	model      string
	maxRetries int
}

func NewClient(
	ctx context.Context,
	apiKey string,
	model string,
	maxRetries int,
) (*Client, error) {

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	if model == "" {
		model = "gemini-2.5-flash"
	}

	if maxRetries < 0 {
		maxRetries = 0
	}

	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize gemini client: %w", err)
	}

	return &Client{
		client:     geminiClient,
		model:      model,
		maxRetries: maxRetries,
	}, nil
}
package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client  *Client
	prompts *PromptLoader
}

func NewGeminiProvider(
	client *Client,
	prompts *PromptLoader,
) LLMProvider {
	return &GeminiProvider{
		client:  client,
		prompts: prompts,
	}
}

// ----------------------------------------------------------------------
// Generate Structured JSON
// ----------------------------------------------------------------------

func (g *GeminiProvider) GenerateStructuredOutput(
	ctx context.Context,
	promptPath string,
	input map[string]any,
	output any,
) error {

	prompt, err := g.prompts.Load(promptPath)
	if err != nil {
		return err
	}

	requestBody, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal prompt input: %w", err)
	}

	finalPrompt := fmt.Sprintf(
		"%s\n\nInput:\n%s",
		prompt,
		string(requestBody),
	)

	err = retry(ctx, g.client.maxRetries, func(ctx context.Context) error {

		response, err := g.client.client.Models.GenerateContent(
			ctx,
			g.client.model,
			genai.Text(finalPrompt),
			&genai.GenerateContentConfig{
				ResponseMIMEType: "application/json",
			},
		)
		if err != nil {
			return fmt.Errorf("gemini request failed: %w", err)
		}

		if len(response.Candidates) == 0 ||
			response.Candidates[0].Content == nil ||
			len(response.Candidates[0].Content.Parts) == 0 {
			return fmt.Errorf("empty response from gemini")
		}

		responseText := response.Text()

		if err := json.Unmarshal([]byte(responseText), output); err != nil {
			return fmt.Errorf("invalid structured json: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------------
// Generate Free-form Text
// ----------------------------------------------------------------------

func (g *GeminiProvider) GenerateText(
	ctx context.Context,
	promptPath string,
	input map[string]any,
) (string, error) {

	prompt, err := g.prompts.Load(promptPath)
	if err != nil {
		return "", err
	}

	requestBody, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal prompt input: %w", err)
	}

	finalPrompt := fmt.Sprintf(
		"%s\n\nInput:\n%s",
		prompt,
		string(requestBody),
	)

	var responseText string

	err = retry(ctx, g.client.maxRetries, func(ctx context.Context) error {

		response, err := g.client.client.Models.GenerateContent(
			ctx,
			g.client.model,
			genai.Text(finalPrompt),
			nil,
		)
		if err != nil {
			return fmt.Errorf("gemini request failed: %w", err)
		}

		if len(response.Candidates) == 0 ||
			response.Candidates[0].Content == nil ||
			len(response.Candidates[0].Content.Parts) == 0 {
			return fmt.Errorf("empty response from gemini")
		}

		responseText = response.Text()
		return nil
	})

	if err != nil {
		return "", err
	}

	return responseText, nil
}
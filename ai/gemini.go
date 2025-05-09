package ai

import (
	"context"

	"google.golang.org/genai"
)

type GeminiAIClient struct {
	client *genai.Client
}

func NewGeminiAIClient(apiKey string) (*GeminiAIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiAIClient{ client: client }, nil
}

// Compile-time check to ensure GeminiAIClient implements AIClient interface
var _ AIClient = (*GeminiAIClient)(nil)

func (c *GeminiAIClient) GenerateChangelogEntry(params GenerateChangelogEntryParams) (GenerateChangelogEntryResponse, error) {
	// TODO: Implement the changelog generation logic using the Gemini API - https://linear.app/ammar-ahmed/issue/CHL-7/prompt-with-schema-response-with-gemini
	return GenerateChangelogEntryResponse{}, nil
}

package ai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) (*OpenAIClient, error) {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &OpenAIClient{client: &client}, nil
}

// Compile-time check to ensure OpenAIClient implements AIClient interface
var _ AIClient = (*OpenAIClient)(nil)

func (c *OpenAIClient) GenerateChangelog(from, to string) (string, error) {
	return "TODO", nil
}


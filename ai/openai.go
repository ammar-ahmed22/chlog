package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ammar-ahmed22/chlog/git"
	"github.com/ammar-ahmed22/chlog/models"
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

func (c *OpenAIClient) GenerateChangelogEntry(params GenerateChangelogEntryParams) (GenerateChangelogEntryResponse, error) {
	ctx := context.Background()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "changelog_entry",
		Description: openai.String("The change log entry for the commit range"),
		Schema:      models.ChangelogEntrySchema,
		Strict:      openai.Bool(true),
	}

	historyWithDiff, err := git.CommitHistoryWithDiff(params.FromCommit, params.ToCommit)
	if err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("failed to get commit history with diff: %v", err)
	}

	response, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: params.Model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fmt.Sprintf(Prompt, params.Tags, historyWithDiff)),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
	})

	if err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("failed to AI generate changelog: %v", err)
	}

	if len(response.Choices) == 0 {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("no response from OpenAI")
	}

	resp := response.Choices[0].Message.Content

	var changeLogEntry models.ChangelogEntry
	if err := json.Unmarshal([]byte(resp), &changeLogEntry); err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("Invalid JSON response from OpenAI. Please try again.")
	}

	return GenerateChangelogEntryResponse{
		Entry:        changeLogEntry,
		InputTokens:  int(response.Usage.PromptTokens),
		OutputTokens: int(response.Usage.CompletionTokens),
	}, nil
}

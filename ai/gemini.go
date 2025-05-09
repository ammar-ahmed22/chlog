package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ammar-ahmed22/chlog/git"
	"github.com/ammar-ahmed22/chlog/models"
	"google.golang.org/genai"
)

type GeminiAIClient struct {
	client *genai.Client
}

func NewGeminiAIClient(apiKey string) (*GeminiAIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &GeminiAIClient{client: client}, nil
}

// Compile-time check to ensure GeminiAIClient implements AIClient interface
var _ AIClient = (*GeminiAIClient)(nil)

func (c *GeminiAIClient) GenerateChangelogEntry(params GenerateChangelogEntryParams) (GenerateChangelogEntryResponse, error) {
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"version": {
					Type:        genai.TypeString,
					Description: "The version number of the release. Will be provided in the prompt.",
				},
				"date": {
					Type:        genai.TypeString,
					Description: "The date of the release. Will be provided in the prompt.",
				},
				"changes": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"details": {
								Type:        genai.TypeString,
								Description: "Summarized details of the change. Should only be a single sentence starting with a past-tense verb.",
							},
							"commit_hash": {
								Type:        genai.TypeString,
								Description: "The commit hash associated with this change.",
							},
							"tags": {
								Type: genai.TypeArray,
								Items: &genai.Schema{
									Type:        genai.TypeString,
									Enum:        []string{"added", "changed", "removed", "deprecated", "security", "fixed"},
									Description: "Tags associated with this change",
								},
							},
						},
						Required:         []string{"details", "commit_hash", "tags"},
						PropertyOrdering: []string{"details", "commit_hash", "tags"},
					},
				},
			},
			PropertyOrdering: []string{"version", "date", "changes"},
			Required:         []string{"version", "date", "changes"},
		},
	}

	historyWithDiff, err := git.CommitHistoryWithDiff(params.FromCommit, params.ToCommit)
	if err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("failed to get commit history with diff: %w", err)
	}

	result, err := c.client.Models.GenerateContent(
		context.Background(),
		params.Model,
		genai.Text(fmt.Sprintf(Prompt, params.Version, params.Date, historyWithDiff)),
		config,
	)
	if err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("Failed to AI generate changelog: %v", err)
	}

	resp := result.Text()
	var changelogEntry models.ChangelogEntry
	if err := json.Unmarshal([]byte(resp), &changelogEntry); err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("Invalid JSON response from Gemini. Please try again.\nGenerated response: %s", resp)
	}
	return GenerateChangelogEntryResponse{
		Entry:        changelogEntry,
		InputTokens:  int(result.UsageMetadata.PromptTokenCount),
		OutputTokens: int(result.UsageMetadata.CandidatesTokenCount),
	}, nil
}

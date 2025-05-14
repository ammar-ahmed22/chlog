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
					Description: "The version number of the release. Leave as empty string.",
				},
				"date": {
					Type:        genai.TypeString,
					Description: "The date of the release. Leave as empty string.",
				},
				"from_ref": {
					Type:        genai.TypeString,
					Description: "The starting commit reference for the changelog entry. Leave as empty string.",
				},
				"to_ref": {
					Type:        genai.TypeString,
					Description: "The ending commit reference for the changelog entry. Leave as empty string.",
				},
				"changes": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"title": {
								Type:        genai.TypeString,
								Description: "The title of the change. Should be succinct.",
							},
							"description": {
								Type:        genai.TypeString,
								Description: "End-user friendly description of the change. Should be more verbose.",
							},
							"impact": {
								Type:        genai.TypeString,
								Description: "The impact of the change. Describe what and how the change affects the user or usage of the software.",
							},
							"commits": {
								Type:        genai.TypeArray,
								Description: "List of commit hashes associated with this change. Must have at least one value.",
								Items: &genai.Schema{
									Type: genai.TypeString,
								},
							},
							"tags": {
								Type: genai.TypeArray,
								Items: &genai.Schema{
									Type:        genai.TypeString,
									Description: "Tags associated with this change",
								},
							},
						},
						Required:         []string{"title", "description", "impact", "commits", "tags"},
						PropertyOrdering: []string{"title", "description", "impact", "commits", "tags"},
					},
				},
			},
			PropertyOrdering: []string{"version", "date", "from_ref", "to_ref", "changes"},
			Required:         []string{"changes"},
		},
	}

	historyWithDiff, err := git.CommitHistoryWithDiff(params.FromCommit, params.ToCommit)
	if err != nil {
		return GenerateChangelogEntryResponse{}, fmt.Errorf("failed to get commit history with diff: %w", err)
	}

	result, err := c.client.Models.GenerateContent(
		context.Background(),
		params.Model,
		genai.Text(fmt.Sprintf(Prompt, params.Tags, historyWithDiff)),
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

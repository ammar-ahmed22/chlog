package ai

import (
	"fmt"
	"slices"

	"github.com/ammar-ahmed22/chlog/models"
)

type GenerateChangelogEntryParams struct {
	Version    string
	Date       string
	FromCommit string
	ToCommit   string
	Model      string
	Tags       []string
}

type GenerateChangelogEntryResponse struct {
	Entry        models.ChangelogEntry
	InputTokens  int
	OutputTokens int
}

type AIClient interface {
	GenerateChangelogEntry(params GenerateChangelogEntryParams) (GenerateChangelogEntryResponse, error)
}

var DefaultTags = []string{"feature", "fix", "improvement", "deprecation", "security", "breaking", "documentation"}

var Prompt = `
You are a changelog generation assistant. Based on the provided Git commits and their diffs, generate a structured changelog entry that adheres exactly to the JSON schema.

## Rules:
- Only use the information provided in the commit messages and diffs.
- Each change should include a succint title, a detailed, end-user friendly description, and an impact statement.
- Each change must be tagged appropriately. Valid tags are:
  - %s
- Each change must have at least one tag.
- Each change should include the commit hash or hashes (if multiple) associated with it.
- Each change must be associated with at least one commit.
- You can have multiple changes associated with a single commit, up to your discretion.
- Ordering of changes should be from most recent to oldest (most recent first, oldest last).
- Output must be strictly valid JSON matching the schema. Do not include any explanation or extra text.

## Git Commits:
Each commit is shown below with its hash, message, and code diff separated by "--- COMMIT ---".

%s
	`

func NewAIClient(provider, apiKey string) (AIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for provider: %s", provider)
	}

	switch provider {
	case "openai":
		return NewOpenAIClient(apiKey)
	case "gemini":
		return NewGeminiAIClient(apiKey)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

var ProvidersMap = map[string][]string{
	"openai": {"gpt-4o-mini", "gpt-4.1-mini"},
	"gemini": {"gemini-2.0-flash"},
}

var ProviderEnvVarMap = map[string]string{
	"openai": "OPENAI_API_KEY",
	"gemini": "GEMINI_API_KEY",
}

func SupportedProviders() []string {
	providers := make([]string, 0, len(ProvidersMap))
	for provider := range ProvidersMap {
		providers = append(providers, provider)
	}
	return providers
}

func IsValidProvider(provider string) bool {
	_, ok := ProvidersMap[provider]
	return ok
}

func IsValidModel(provider, model string) bool {
	models, ok := ProvidersMap[provider]
	if !ok {
		return false
	}
	return slices.Contains(models, model)
}

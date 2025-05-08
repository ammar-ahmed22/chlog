package ai

type AIClient interface {
	GenerateChangelog(from, to string) (string, error)
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
	for _, m := range models {
		if m == model {
			return true
		}
	}
	return false
}


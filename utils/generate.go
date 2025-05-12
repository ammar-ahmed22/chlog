package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ammar-ahmed22/chlog/ai"
	"github.com/ammar-ahmed22/chlog/git"
	"github.com/ammar-ahmed22/chlog/models"
	"github.com/spf13/cobra"
)

type GenerateFlags struct {
	From                  string
	To                    string
	Verbose               bool
	Provider              string
	Model                 string
	Date                  string
	APIKey                string
	Pretty                bool
	ExistingChangelog     []models.ChangelogEntry
	ExistingChangelogPath string
}

func ParseGenerateFlags(cmd *cobra.Command) (*GenerateFlags, error) {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}

	err = LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("Error loading config file '%s': %v", configPath, err)
	}

	from, err := cmd.Flags().GetString("from")
	if err != nil {
		return nil, err
	}
	to, err := cmd.Flags().GetString("to")
	if err != nil {
		return nil, err
	}
	err = git.IsValidRef(from)
	if err != nil {
		return nil, fmt.Errorf("Invalid '--from, -f' reference: '%s'. Make sure it's a valid Git commit, tag, or branch (e.g. 'HEAD', 'main', 'v1.0.0', or 'abc1234')", from)
	}

	err = git.IsValidRef(to)
	if err != nil {
		return nil, fmt.Errorf("Invalid '--to, -t' reference: '%s'. Make sure it's a valid Git commit, tag, or branch (e.g. 'HEAD', 'main', 'v1.0.0', or 'abc1234')", to)
	}

	verbose, err := GetConfigFlagBool(cmd, "verbose")
	if err != nil {
		return nil, err
	}

	provider, _, err := GetConfigFlagString(cmd, "provider") 
	if err != nil {
		return nil, err
	}

	ok := ai.IsValidProvider(provider)
	if !ok {
		return nil, fmt.Errorf("Invalid provider '%s'. Supported providers are: %s", provider, ai.SupportedProviders())
	}

	model, _, err := GetConfigFlagString(cmd, "model")
	if err != nil {
		return nil, err
	}

	if model != "" {
		ok = ai.IsValidModel(provider, model)
		if !ok {
			return nil, fmt.Errorf("Invalid model '%s' for provider '%s'. Supported models are: %s", model, provider, ai.ProvidersMap[provider])
		}
	} else {
		model = ai.ProvidersMap[provider][0] // Default to the first model for the provider
	}

	apiKey, _, err := GetConfigFlagString(cmd, "apiKey") 
	if err != nil {
		return nil, err
	}

	if apiKey == "" {
		envVar := ai.ProviderEnvVarMap[provider]
		value, ok := os.LookupEnv(envVar)
		if !ok {
			return nil, fmt.Errorf("API key for provider '%s' is required. Set it using the '--apiKey' flag or the environment variable '%s'", provider, envVar)
		}
		apiKey = value
	}

	date, err := cmd.Flags().GetString("date")
	if err != nil {
		return nil, err
	}

	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("Invalid date format '%s'. Use YYYY-MM-DD format", date)
	}

	pretty, err := GetConfigFlagBool(cmd, "pretty") 
	if err != nil {
		return nil, err
	}

	file, fileFromConfig, err := GetConfigFlagString(cmd, "file")
	if err != nil {
		return nil, err
	}

	var existingChangelog []models.ChangelogEntry
	if file != "" {
		if fileFromConfig {
			// Join the config path with the file path
			configDir := filepath.Dir(configPath)
			joined := filepath.Join(configDir, file)
			absPath, err := filepath.Abs(joined)
			if err != nil {
				return nil, fmt.Errorf("Error getting absolute path for file '%s': %v", file, err)
			}
			file = absPath
		}
		existingChangelog, err = ParseAndValidateChangelogFile(file)
		if err != nil {
			return nil, err
		}
	}

	return &GenerateFlags{
		From:                  from,
		To:                    to,
		Verbose:               verbose,
		Provider:              provider,
		Model:                 model,
		Date:                  date,
		APIKey:                apiKey,
		Pretty:                pretty,
		ExistingChangelog:     existingChangelog,
		ExistingChangelogPath: file,
	}, nil
}


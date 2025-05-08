/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/ammar-ahmed22/chlog/ai"
	"github.com/ammar-ahmed22/chlog/git"
	"github.com/spf13/cobra"
)

type Flags struct {
	From     string
	To       string
	Verbose  bool
	Provider string
	Model    string
	APIKey   string
}

func ParseAndValidateFlags(cmd *cobra.Command) (*Flags, error) {
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

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return nil, err
	}

	provider, err := cmd.Flags().GetString("provider")
	if err != nil {
		return nil, err
	}

	ok := ai.IsValidProvider(provider)
	if !ok {
		return nil, fmt.Errorf("Invalid provider '%s'. Supported providers are: %s", provider, ai.SupportedProviders())
	}

	model, err := cmd.Flags().GetString("model")
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

	apiKey, err := cmd.Flags().GetString("apiKey")
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

	return &Flags{
		From:     from,
		To:       to,
		Verbose:  verbose,
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
	}, nil
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the AI-powered changelog",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := git.IsInstalled()

		flags, err := ParseAndValidateFlags(cmd)
		if err != nil {
			return err
		}

		aiClient, err := ai.NewAIClient(flags.Provider, flags.APIKey)
		if err != nil {
			return err
		}

		logs, err := git.LogRange(flags.From, flags.To)
		if err != nil {
			return fmt.Errorf("Error getting git log: %v", err)
		}

		if flags.Verbose {
			fmt.Println("Generating changelog for commits:")
			for _, log := range logs {
				fmt.Println(log)
			}
		}

		commits, err := git.CommitRange(flags.From, flags.To)
		if err != nil {
			return fmt.Errorf("Error getting commits: %v", err)
		}

		historyWithDiff := ""
		for _, commit := range commits {
			details, err := git.CommitDetails(commit)
			if err != nil {
				return err
			}

			historyWithDiff += fmt.Sprintf("--- COMMIT ---\n%s\n", details)
		}


		changelog, err := aiClient.GenerateChangelog(flags.From, flags.To)
		if err != nil {
			return fmt.Errorf("Error generating changelog: %v", err)
		}

		fmt.Println(changelog)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("from", "f", "HEAD~1", "Starting commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234)")
	generateCmd.Flags().StringP("to", "t", "HEAD", "Ending commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234)")
	generateCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	generateCmd.Flags().StringP("provider", "p", "openai", "LLM provider (see chlog models for available options)")
	generateCmd.Flags().StringP("model", "m", "", "LLM model (see chlog models for available options and defaults)")
	generateCmd.Flags().String("apiKey", "", "API key for the LLM provider (can also be set via environment variable, see chlog models for details)")
}

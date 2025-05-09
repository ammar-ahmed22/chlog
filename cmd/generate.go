package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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
	Date     string
	APIKey   string
	Pretty   bool
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

	date, err := cmd.Flags().GetString("date")
	if err != nil {
		return nil, err
	}

	pretty, err := cmd.Flags().GetBool("pretty")
	if err != nil {
		return nil, err
	}

	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("Invalid date format '%s'. Use YYYY-MM-DD format", date)
	}

	return &Flags{
		From:     from,
		To:       to,
		Verbose:  verbose,
		Provider: provider,
		Model:    model,
		Date:     date,
		APIKey:   apiKey,
		Pretty:   pretty,
	}, nil
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   fmt.Sprintf("generate <VERSION> (default \"%s\")", time.Now().Format("2006-01-02")),
	Short: "Generates the AI-powered changelog entry for the specified version",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := git.IsInstalled()
		if err != nil {
			return err
		}

		var version string
		if len(args) > 0 {
			version = args[0]
		} else {
			version = time.Now().Format("2006-01-02")
		}

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
			fmt.Fprintf(os.Stderr, "Generating changelog entry \"%s\" for commits:\n", version)
			for _, log := range logs {
				fmt.Fprintln(os.Stderr, log)
			}
		}

		if flags.Verbose {
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Starting AI changelog generation")
		}

		response, err := aiClient.GenerateChangelogEntry(ai.GenerateChangelogEntryParams{
			FromCommit: flags.From,
			ToCommit:   flags.To,
			Model:      flags.Model,
			Version:    version,
			Date:       flags.Date,
		})
		if err != nil {
			return fmt.Errorf("Error generating changelog: %v", err)
		}

		if flags.Verbose {
			fmt.Fprintln(os.Stderr, "Completed AI changelog generation")
			fmt.Fprintf(os.Stderr, "tokens used: %d (input: %d, output: %d)\n", response.InputTokens+response.OutputTokens, response.InputTokens, response.OutputTokens)
		}

		if flags.Pretty {
			pretty, err := json.MarshalIndent(response.Entry, "", "  ")
			if err != nil {
				return fmt.Errorf("Error pretty printing JSON: %v", err)
			}
			fmt.Println(string(pretty))
			return nil
		}
		jsonOutput, err := json.Marshal(response.Entry)
		if err != nil {
			return fmt.Errorf("Error generating JSON: %v", err)
		}
		fmt.Println(string(jsonOutput))
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
	generateCmd.Flags().StringP("date", "d", time.Now().Format("2006-01-02"), "Date for the changelog entry in YYYY-MM-DD format")
	generateCmd.Flags().Bool("pretty", false, "Prettified JSON output")
}

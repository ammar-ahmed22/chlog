package cmd

import (
	"fmt"
	"strings"

	"github.com/ammar-ahmed22/chlog/ai"
	"github.com/spf13/cobra"
)

func printModels(provider string, models []string, envVar string) {
	fmt.Printf("Provider: %s", provider)
	if envVar != "" {
		fmt.Printf(" (env var: %s)", envVar)
	}
	fmt.Println()
	for i, model := range models {
		if i == 0 {
			fmt.Printf(" - %s (default)\n", model)
			continue
		}
		fmt.Printf(" - %s\n", model)
	}
	fmt.Println()
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List supported models and providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		provider, err := cmd.Flags().GetString("provider")
		if err != nil {
			return err
		}

		if provider == "all" {
			fmt.Println("Supported providers and models:")
			fmt.Println()
			for p, model := range ai.ProvidersMap {
				envVar := ai.ProviderEnvVarMap[p]
				printModels(p, model, envVar)
			}
			return nil
		}

		if model, ok := ai.ProvidersMap[provider]; ok {
			fmt.Printf("Supported models for provider '%s':\n", provider)
			envVar := ai.ProviderEnvVarMap[provider]
			printModels(provider, model, envVar)
			return nil
		}

		return fmt.Errorf("Invalid provider '%s'. Supported providers are: %s", provider, strings.Join(ai.SupportedProviders(), ", "))
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)

	modelsCmd.Flags().StringP("provider", "p", "all", "The provider to list models for (openai, gemini or all)")
}

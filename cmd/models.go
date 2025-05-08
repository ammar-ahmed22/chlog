/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

// modelsCmd represents the models command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	modelsCmd.Flags().StringP("provider", "p", "all", "The provider to list models for (openai, gemini or all)")
}

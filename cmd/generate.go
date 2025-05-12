package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ammar-ahmed22/chlog/ai"
	"github.com/ammar-ahmed22/chlog/git"
	"github.com/ammar-ahmed22/chlog/models"
	"github.com/ammar-ahmed22/chlog/utils"
	"github.com/spf13/cobra"
)

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

		flags, err := utils.ParseGenerateFlags(cmd)
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
			utils.Eprintf("Generating changelog entry \"%s\" for commits:\n", version)
			for _, log := range logs {
				utils.Eprintln(log)
			}
		}

		if flags.Verbose {
			utils.Eprintln("")
			utils.Eprintf("Starting AI changelog generation (provider: %s, model: %s)\n", flags.Provider, flags.Model)
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
			utils.Eprintln("Completed AI changelog generation")
			utils.Eprintf("tokens used: %d (input: %d, output: %d)\n", response.InputTokens+response.OutputTokens, response.InputTokens, response.OutputTokens)
		}

		jsonOutput, err := json.Marshal(response.Entry)
		if err != nil {
			return fmt.Errorf("Error generating JSON: %v", err)
		}

		if flags.ExistingChangelog != nil {
			if flags.Verbose {
				utils.Eprintf("Writing to existing changelog file '%s'\n", flags.ExistingChangelogPath)
			}

			// NOTE: Adding the new entry to the beginning. This is not good for performance but OK for POC.
			updatedChangelog := append([]models.ChangelogEntry{response.Entry}, flags.ExistingChangelog...)
			err := utils.WriteChangelogFile(flags.ExistingChangelogPath, updatedChangelog)
			if err != nil {
				return fmt.Errorf("Error writing changelog file '%s': %v", flags.ExistingChangelogPath, err)
			}
		}

		if flags.Pretty {
			pretty, err := json.MarshalIndent(response.Entry, "", "  ")
			if err != nil {
				return fmt.Errorf("Error pretty printing JSON: %v", err)
			}
			fmt.Println(string(pretty))
		} else {
			fmt.Println(string(jsonOutput))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("config", "c", "", "Path to config file (optional, will be loaded if present in the current directory)")
	generateCmd.Flags().StringP("from", "f", "HEAD~1", "Starting commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234)")
	generateCmd.Flags().StringP("to", "t", "HEAD", "Ending commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234)")
	generateCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	generateCmd.Flags().StringP("provider", "p", "openai", "LLM provider (see chlog models for available options)")
	generateCmd.Flags().StringP("model", "m", "", "LLM model (see chlog models for available options and defaults)")
	generateCmd.Flags().String("apiKey", "", "API key for the LLM provider (can also be set via environment variable, see chlog models for details)")
	generateCmd.Flags().StringP("date", "d", time.Now().Format("2006-01-02"), "Date for the changelog entry in YYYY-MM-DD format")
	generateCmd.Flags().Bool("pretty", false, "Prettified JSON output")
	generateCmd.Flags().String("file", "", "Path to existing changelog JSON file to update with the new entry (should be an array of changelog entries or empty file)")
}

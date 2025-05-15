package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ammar-ahmed22/chlog/ai"
	"github.com/ammar-ahmed22/chlog/git"
	"github.com/ammar-ahmed22/chlog/models"
	"github.com/ammar-ahmed22/chlog/utils"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
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
			utils.Eprintf("\u2192 Generating changelog entry %s\n", color.CyanString(version))
			utils.Eprintln("\u2192 Using commits:")
			for _, log := range logs {
				parts := strings.SplitN(log, " ", 2)
				hash := parts[0]
				message := parts[1]
				utils.Eprintf(" \u2192 %s %s\n", color.YellowString(hash), message)
			}
		}

		spnr := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		spnr.Writer = os.Stderr
		if flags.Verbose {
			utils.Eprintf("\u2192 Using AI provider: %s\n", color.MagentaString("%s (model: %s)", flags.Provider, flags.Model))
			spnr.Suffix = fmt.Sprintf(" AI Generating changelog entry...")
			spnr.Start()
			defer spnr.Stop()
		}

		response, err := aiClient.GenerateChangelogEntry(ai.GenerateChangelogEntryParams{
			FromCommit: flags.From,
			ToCommit:   flags.To,
			Model:      flags.Model,
			Version:    version,
			Date:       flags.Date,
			Tags:       ai.DefaultTags,
		})
		if err != nil {
			return fmt.Errorf("Error generating changelog: %v", err)
		}

		response.Entry.Version = version
		response.Entry.Date = flags.Date
		response.Entry.FromRef = flags.From
		response.Entry.ToRef = flags.To
		// Add id to each change
		for i, change := range response.Entry.Changes {
			response.Entry.Changes[i].ID = utils.TruncatedKebabCase(change.Title, 40)
		}

		if flags.Verbose {
			spnr.Stop()
			utils.Eprintf("%s AI Generated changelog entry\n", color.GreenString("\u2713"))
			utils.Eprintf("\u2192 Tokens used: %d\n", response.InputTokens+response.OutputTokens)
			utils.Eprintf(" \u2192 Input: %d\n", response.InputTokens)
			utils.Eprintf(" \u2192 Output: %d\n", response.OutputTokens)
		}

		jsonOutput, err := json.Marshal(response.Entry)
		if err != nil {
			return fmt.Errorf("Error generating JSON: %v", err)
		}

		if flags.ExistingChangelog != nil {
			if flags.Verbose {
				utils.Eprintf("\u2192 Writing to changelog file '%s'\n", flags.ExistingChangelogPath)
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

	generateCmd.Flags().StringP("config", "c", "", "Path to config file (optional, chlog.yaml will be loaded if present in the current directory)")
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

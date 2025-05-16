package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/ammar-ahmed22/chlog/ai"
	"github.com/ammar-ahmed22/chlog/models"
	"github.com/ammar-ahmed22/chlog/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ChangelogFile struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Repository  string                  `json:"repository"`
	Entries     []models.ChangelogEntry `json:"entries"`
}

type ChlogConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
	APIKey   string `yaml:"api_key,omitempty"`
	File     string `yaml:"file"`
	Pretty   bool   `yaml:"pretty"`
	Verbose  bool   `yaml:"verbose"`
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a changelog file and generate config interactively",
	Long: `Initialize a new changelog file interactively with metadata fields like title, description, and repository URL.

You can also optionally create a configuration file for the chlog generate command, which includes settings for the generation process, such as the LLM provider, model, API key, and other options to avoid repeating flags during generation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Eprintln("This utility will help you create a changelog file with metadata and configuration.")
		utils.Eprintln("")
		utils.Eprintln("See `chlog init --help` for more information.")
		utils.Eprintln("")
		utils.Eprintln("Press Ctrl+C to cancel at any time.")

		fileName, err := utils.Prompt("changelog file name", "changelog.json")
		if err != nil {
			return err
		}

		var defaultTitle string = "Changelog"
		wd, err := os.Getwd()
		if err == nil {
			dirName := filepath.Base(wd)
			defaultTitle = dirName
		}
		title, err := utils.Prompt("title", defaultTitle)
		if err != nil {
			return err
		}

		description, err := utils.Prompt("description", "")
		if err != nil {
			return err
		}

		repoUrl, err := utils.Prompt("repository", "")
		if err != nil {
			return err
		}

		file := &ChangelogFile{
			Title:       title,
			Description: description,
			Repository:  repoUrl,
			Entries:     []models.ChangelogEntry{},
		}
		fileContent, err := json.MarshalIndent(file, "", "  ")
		if err != nil {
			return fmt.Errorf("Failed to marshal JSON: %v", err)
		}

		err = os.WriteFile(fileName, fileContent, 0644)
		if err != nil {
			return fmt.Errorf("Failed to write file: %v", err)
		}

		utils.Eprintln("")
		utils.Eprintf("\u2192 Writing to '%s'\n", fileName)
		fmt.Println(string(fileContent))

		createConfig, err := utils.Confirm("Would you like to create a chlog configuration file")
		if err != nil {
			return err
		}

		if !createConfig {
			return nil
		}

		configFileName, err := utils.Prompt("chlog config file name", "chlog.yaml")
		if err != nil {
			return err
		}

		provider, err := utils.Select("LLM provider", ai.SupportedProviders())
		if err != nil {
			return err
		}

		model, err := utils.Select("LLM model", ai.SupportedModels(provider))
		if err != nil {
			return err
		}

		apiKey, err := utils.Prompt("API key (can also be set via environment variable)", "")
		if err != nil {
			return err
		}

		enablePrettyPrint, err := utils.Confirm("Enable pretty print")
		if err != nil {
			return err
		}

		enableVerbose, err := utils.Confirm("Enable verbose output")
		if err != nil {
			return err
		}

		config := &ChlogConfig{
			Provider: provider,
			Model:    model,
			APIKey:   apiKey,
			File:     fileName,
			Pretty:   enablePrettyPrint,
			Verbose:  enableVerbose,
		}

		utils.Eprintln("")
		utils.Eprintf("\u2192 Writing to '%s'\n", configFileName)
		configContent, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("Failed to marshal YAML: %v", err)
		}

		err = os.WriteFile(configFileName, configContent, 0644)
		if err != nil {
			return fmt.Errorf("Failed to write config file: %v", err)
		}
		fmt.Println(string(configContent))
		utils.Eprintf("You can now use the `%s` command to generate a changelog entry.\n", color.CyanString("chlog generate"))
		utils.Eprintf("See `%s` for more information.\n", color.CyanString("chlog generate --help"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

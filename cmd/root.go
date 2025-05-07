/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chlog",
	Short: "Generate AI-powered changelogs from your Git history",
	Long: `chlog is a command-line tool that uses AI to generate and update structured changelogs from your Git history.

It can automatically summarize changes based on diffs and commit messages and output structured changelogs in various formats, including Markdown and JSON.
Use it to keep your changelogs clean, consistent, and up-to-date.

Example:
	chlog generate --format json --from HEAD~10 --to HEAD --out changelog.json`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chlog.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}



package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



var rootCmd = &cobra.Command{
	Use:   "chlog",
	Short: "Generate AI-powered changelogs from your Git history",
	Long: `chlog is a command-line tool that uses AI to generate and update structured changelogs from your Git history.

It can automatically summarize changes based on diffs and commit messages and output structured changelogs in JSON formats.
Use it to keep your changelogs clean, consistent, and up-to-date.

Example:
	chlog generate 0.2.0 --from HEAD~10 --to HEAD > changelog.json`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

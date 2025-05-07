/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ammar-ahmed22/chlog/git"
	"github.com/spf13/cobra"
)


type Flags struct {
	From    string
	To      string
	Verbose bool
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

	return &Flags{
		From:    from,
		To:      to,
		Verbose: verbose,
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// generateCmd.Flags().String("fmt", "json", "Output format (json, markdown)")
	generateCmd.Flags().StringP("from", "f", "HEAD~1", "Starting commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234)")
	generateCmd.Flags().StringP("to", "t", "HEAD", "Ending commit reference (e.g. HEAD~3, main, v1.0.0, or abc1234)")
	generateCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}

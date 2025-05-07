/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func isGitInstalled() error {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Git is not installed or not found in PATH: %v", err)
	}
	return nil
}

func isValidGitRef(ref string) error {
	cmd := exec.Command("git", "rev-parse", "--verify", ref)
	return cmd.Run()
}

func getGitLog(from, to string) ([]string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%h %s", from+"..."+to)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error getting git log: %v", err)
	}
	return []string{string(output)}, nil
}

func getCommits(from, to string) ([]string, error) {
	cmd := exec.Command("git", "rev-list", "--reverse", fmt.Sprintf("%s..%s", from, to))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error getting commits: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func getCommitDetails(commit string) (string, error) {
	cmd := exec.Command("git", "show", commit)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error getting commit details: %v", err)
	}
	return string(out), nil
}

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
	err = isValidGitRef(from)
	if err != nil {
		return nil, fmt.Errorf("Invalid '--from, -f' reference: '%s'. Make sure it's a valid Git commit, tag, or branch (e.g. 'HEAD', 'main', 'v1.0.0', or 'abc1234')", from)
	}

	err = isValidGitRef(to)
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
		err := isGitInstalled()
		if err != nil {
			return err
		}

		flags, err := ParseAndValidateFlags(cmd)
		if err != nil {
			return err
		}

		logs, err := getGitLog(flags.From, flags.To)
		if err != nil {
			return fmt.Errorf("Error getting git log: %v", err)
		}

		if flags.Verbose {
			fmt.Println("Generating changelog for commits:")
			for _, log := range logs {
				fmt.Println(log)
			}
		}

		commits, err := getCommits(flags.From, flags.To)
		if err != nil {
			return fmt.Errorf("Error getting commits: %v", err)
		}

		historyWithDiff := ""
		for _, commit := range commits {
			details, err := getCommitDetails(commit)
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

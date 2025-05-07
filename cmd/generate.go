/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os/exec"

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

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the AI-powered changelog",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := isGitInstalled()
		if err != nil {
			return err
		}

		from, err := cmd.Flags().GetString("from")
		if err != nil {
			fmt.Println("Error getting 'from' flag:", err)
			return err
		}
		to, err := cmd.Flags().GetString("to")
		if err != nil {
			fmt.Println("Error getting 'to' flag:", err)
			return err
		}

		err = isValidGitRef(from)
		if err != nil {
			return fmt.Errorf("Invalid '--from, -f' reference: '%s'. Make sure it's a valid Git commit, tag, or branch (e.g. 'HEAD', 'main', 'v1.0.0', or 'abc1234')", from) 
		}

		err = isValidGitRef(to)
		if err != nil {
			return fmt.Errorf("Invalid '--to, -t' reference: '%s'. Make sure it's a valid Git commit, tag, or branch (e.g. 'HEAD', 'main', 'v1.0.0', or 'abc1234')", to) 
		}

		fmt.Println("Generating changelog from", from, "to", to)
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
}

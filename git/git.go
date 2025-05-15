package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func IsInstalled() error {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Git is not installed or not found in PATH: %v", err)
	}
	return nil
}

func IsValidRef(ref string) error {
	cmd := exec.Command("git", "rev-parse", "--verify", ref)
	return cmd.Run()
}

func LogRange(from, to string) ([]string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%h %s", from+"..."+to)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error getting git log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

func CommitRange(from, to string) ([]string, error) {
	cmd := exec.Command("git", "rev-list", "--reverse", fmt.Sprintf("%s..%s", from, to))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error getting commits: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func CommitDetails(commit string) (string, error) {
	cmd := exec.Command("git", "show", commit)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error getting commit details: %v", err)
	}
	return string(out), nil
}

func CommitHistoryWithDiff(from, to string) (string, error) {
	commits, err := CommitRange(from, to)
	if err != nil {
		return "", fmt.Errorf("Error getting commits: %v", err)
	}

	var builder strings.Builder
	for _, commit := range commits {
		details, err := CommitDetails(commit)
		if err != nil {
			return "", err
		}

		builder.WriteString(fmt.Sprintf("--- COMMIT ---\n%s\n", details))
	}

	return builder.String(), nil
}

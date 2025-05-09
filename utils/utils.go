package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ammar-ahmed22/chlog/models"
)

func Eprintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func Eprintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func ParseAndValidateChangelogFile(path string) ([]models.ChangelogEntry, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("File '%s' does not exist", path)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading file '%s': %v", path, err)
	}

	if len(contents) == 0 {
		return []models.ChangelogEntry{}, nil
	}

	var changelog []models.ChangelogEntry
	err = json.Unmarshal(contents, &changelog)
	if err != nil {
		return nil, fmt.Errorf("Changelog file '%s' is not valid JSON. See https://github.com/ammar-ahmed22/chlog#changelog-format for the expected format", path)
	}

	return changelog, nil
}

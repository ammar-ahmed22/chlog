package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ammar-ahmed22/chlog/models"
	"github.com/samber/lo"
)

func Eprintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func Eprintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func TruncatedKebabCase(s string, maxLen int) string {
	kebab := lo.KebabCase(s)
	words := strings.Split(kebab, "-")

	var result strings.Builder
	for i, word := range words {
		if result.Len()+len(word) > maxLen {
			break
		}

		if i > 0 && result.Len() > 0 {
			result.WriteString("-")
		}
		result.WriteString(word)
	}
	return result.String()
}

func ParseAndValidateChangelogFile(path string) ([]models.ChangelogEntry, error) {
	_, err := os.Stat(path)
	if err != nil {
		// File does not exist, create one
		if os.IsNotExist(err) {
			// Create an empty changelog file
			err := os.WriteFile(path, []byte("[]"), 0644)
			if err != nil {
				return nil, fmt.Errorf("Error creating changelog file '%s': %v", path, err)
			}
		} else {
			return nil, fmt.Errorf("Error checking changelog file '%s': %v", path, err)
		}
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

func WriteChangelogFile(path string, changelog []models.ChangelogEntry) error {
	contents, err := json.MarshalIndent(changelog, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshalling changelog to JSON: %v", err)
	}
	err = os.WriteFile(path, contents, 0644)
	if err != nil {
		return fmt.Errorf("Error writing changelog file '%s': %v", path, err)
	}
	return nil
}

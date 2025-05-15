package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ammar-ahmed22/chlog/models"
	"github.com/samber/lo"
	"github.com/tidwall/sjson"
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

func hasEntries(contents []byte) (exists bool, entries []models.ChangelogEntry, err error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(contents, &raw); err != nil {
		return false, nil, err
	}

	rawEntries, exists := raw["entries"]
	if !exists {
		return false, nil, nil
	}

	if err := json.Unmarshal(rawEntries, &entries); err != nil {
		return false, nil, err
	}

	return true, entries, nil
}

func ParseAndValidateChangelogFile(path string) ([]models.ChangelogEntry, bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		// File does not exist, create one
		if os.IsNotExist(err) {
			// Create an empty changelog file
			err := os.WriteFile(path, []byte("[]"), 0644)
			if err != nil {
				return nil, false, fmt.Errorf("Error creating changelog file '%s': %v", path, err)
			}
		} else {
			return nil, false, fmt.Errorf("Error checking changelog file '%s': %v", path, err)
		}
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("Error reading file '%s': %v", path, err)
	}

	if len(contents) == 0 {
		return []models.ChangelogEntry{}, false, nil
	}

	var changelogEntries []models.ChangelogEntry
	err = json.Unmarshal(contents, &changelogEntries)
	if err != nil {
		// Check if the file has JSON with "entries" key
		exists, entries, err := hasEntries(contents)
		if exists && entries != nil && err == nil {
			return entries, true, nil
		}
		return nil, false, fmt.Errorf("Changelog file '%s' is not valid JSON. See https://github.com/ammar-ahmed22/chlog#-json-format for the expected format", path)
	}

	return changelogEntries, false, nil
}

func WriteChangelogFile(path string, entriesKey bool, changelog []models.ChangelogEntry) error {
	newEntries, err := json.MarshalIndent(changelog, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshalling changelog to JSON: %v", err)
	}
	if !entriesKey {
		err = os.WriteFile(path, newEntries, 0644)
		if err != nil {
			return fmt.Errorf("Error writing changelog file '%s': %v", path, err)
		}
		return nil
	}
	// Write to the "entries" key in the JSON file
	fileData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading file '%s': %v", path, err)
	}

	updatedContents, err := sjson.SetBytes(fileData, "entries", json.RawMessage(newEntries))
	if err != nil {
		return fmt.Errorf("Error marshalling changelog to JSON: %v", err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, updatedContents, "", "  "); err != nil {
		return fmt.Errorf("Error formatting JSON: %v", err)
	}

	err = os.WriteFile(path, pretty.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Error writing changelog file '%s': %v", path, err)
	}

	return nil
}

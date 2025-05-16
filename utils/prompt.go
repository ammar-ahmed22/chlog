package utils

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func promptLabelDefault(label, defaultValue string) string {
	dim := color.New(color.Faint)
	return fmt.Sprintf("%s %s", label, dim.Sprintf("\u203A %s", defaultValue))
}

func Prompt(prompt, defaultValue string) (string, error) {
	var label string
	if defaultValue != "" {
		label = promptLabelDefault(prompt, defaultValue)
	} else {
		label = prompt
	}
	promptUI := promptui.Prompt{
		Label: label,
	}

	result, err := promptUI.Run()
	if err != nil {
		return "", err
	}

	if result == "" {
		return defaultValue, nil
	}

	return result, nil
}

func Confirm(prompt string) (bool, error) {
	promptUI := promptui.Prompt{
		Label: prompt,
		IsConfirm: true,
		Default: "y",
	}

	result, err := promptUI.Run()
	if err != nil {
		return false, err 
	}

	if result == "n" || result == "N" {
		return false, nil
	}

	return true, nil
}

func Select(prompt string, items []string) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	promptUI := promptui.Select{
		Label: prompt,
		Items: items,
	}

	_, result, err := promptUI.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

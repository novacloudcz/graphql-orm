package tools

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var promptCache map[string]string

func Confirm(text string) bool {
	return ConfirmWithDefault(text, true)
}

func ConfirmWithDefault(text string, defaultValue bool) bool {
	options := "[y/N]"
	if defaultValue {
		options = "[Y/n]"
	}
	response := strings.ToLower(prompt(text+" "+options, nil, false))
	if response != "" && response != "y" && response != "n" {
		return ConfirmWithDefault(text, defaultValue)
	} else if response == "" {
		return defaultValue
	}
	return response != "n"
}

// Prompt ...
func Prompt(text string) string {
	return prompt(text, nil, false)
}

// PromptSecret ...
func PromptSecret(text string) string {
	return prompt(text, nil, true)
}

// PromptCached ...
func PromptCached(text, key string) string {
	// return PromptSecret(name, configKey, false)
	return prompt(text, &key, false)
}

// PromptSecretCached ...
func PromptSecretCached(text, key string) string {
	return prompt(text, &key, true)
}

// PromptWithChoice ...
func PromptWithChoice(name string, choices []string) (int, error) {
	result := ""
	_choices := []string{}

	for _, choice := range choices {
		if choice == "" {
			choice = "none"
		}
		_choices = append(_choices, choice)
	}

	prompt := &survey.Select{
		Message:  name + ":",
		Options:  _choices,
		PageSize: 20,
	}
	if err := survey.AskOne(prompt, &result, nil); err != nil {
		return -1, err
	}

	for index, choise := range _choices {
		if choise == result {
			return index, nil
		}
	}

	return -1, errors.New("invalid choice selected")
}

// PromptWithMultiChoice ...
func PromptWithMultiChoice(name string, choices []string) ([]int, error) {
	var result []string
	prompt := &survey.MultiSelect{
		Message:  name + ":",
		Options:  choices,
		PageSize: 20,
	}
	if err := survey.AskOne(prompt, &result, nil); err != nil {
		return nil, err
	}

	indexes := []int{}
	for index, choise := range choices {
		for _, answer := range result {
			if choise == answer {
				indexes = append(indexes, index)
			}
		}
	}

	return indexes, nil
}

func prompt(text string, key *string, masked bool) string {
	if key != nil {
		if promptCache == nil {
			promptCache = map[string]string{}
		}

		if val := promptCache[*key]; val != "" {
			return val
		}
	}

	fmt.Print(text + ": ")

	if masked {
		val, err := gopass.GetPasswdMasked()
		if err != nil {
			panic(err)
		}
		return string(val)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.Trim(scanner.Text(), " ")
	}

	return ""
}

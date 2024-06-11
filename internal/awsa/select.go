package awsa

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// SelectSecret prompts the user to select an AWS secret and choose an action
func SelectSecret(cfg aws.Config) (string, string, error) {
	// List secrets using the provided AWS config
	secrets, err := ListSecrets(cfg)
	if err != nil {
		return "", "", err
	}

	// Prompt user to select a secret
	var selectedSecret string
	prompt := &survey.Select{
		Message:  "Choose an AWS secret:",
		Options:  secrets,
		PageSize: 14,
	}
	if err := survey.AskOne(prompt, &selectedSecret); err != nil {
		return "", "", err
	}

	// Prompt user to select an action for the chosen secret
	var action string
	actionPrompt := &survey.Select{
		Message:  "Choose an action:",
		Options:  []string{"Edit latest secret version", "View previous versions"},
		PageSize: 14,
	}
	if err := survey.AskOne(actionPrompt, &action); err != nil {
		return "", "", err
	}

	return selectedSecret, action, nil
}

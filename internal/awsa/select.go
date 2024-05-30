package awsa

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// SelectSecret prompts the user to select an AWS secret and choose an action
func SelectSecret(cfg aws.Config) (string, string, error) {
	secrets, err := ListSecrets(cfg)
	if err != nil {
		return "", "", err
	}

	var selectedSecret string
	prompt := &survey.Select{
		Message: "Choose an AWS secret:",
		Options: secrets,
	}
	err = survey.AskOne(prompt, &selectedSecret, survey.WithPageSize(14))
	if err != nil {
		return "", "", err
	}

	var action string
	actionPrompt := &survey.Select{
		Message: "Choose an action:",
		Options: []string{"Edit latest secret version", "View previous versions"},
	}
	err = survey.AskOne(actionPrompt, &action)
	if err != nil {
		return "", "", err
	}

	return selectedSecret, action, nil
}

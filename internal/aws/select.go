package aws

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// SelectSecret prompts the user to select an AWS secret
func SelectSecret(cfg aws.Config) (string, error) {
	secrets, err := ListSecrets(cfg)
	if err != nil {
		return "", err
	}

	var selectedSecret string
	prompt := &survey.Select{
		Message: "Choose an AWS secret:",
		Options: secrets,
	}
	err = survey.AskOne(prompt, &selectedSecret)
	if err != nil {
		return "", err
	}

	return selectedSecret, nil
}

package aws

import (
	"bufio"
	"context"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"os"
	"path/filepath"
	"strings"
)

// GetProfiles retrieves the list of profiles from the AWS config file
func GetProfiles() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".aws", "config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[profile ") && strings.HasSuffix(line, "]") {
			profile := strings.TrimPrefix(line, "[profile ")
			profile = strings.TrimSuffix(profile, "]")
			profiles = append(profiles, profile)
		}
		// Handle default profile case
		if line == "[default]" {
			profiles = append(profiles, "default")
		}
	}
	return profiles, scanner.Err()
}

// SelectProfile prompts the user to select an AWS profile
func SelectProfile() (string, error) {
	profiles, err := GetProfiles()
	if err != nil {
		return "", err
	}

	var selectedProfile string
	prompt := &survey.Select{
		Message: "Choose an AWS profile:",
		Options: profiles,
	}
	err = survey.AskOne(prompt, &selectedProfile)
	if err != nil {
		return "", err
	}

	return selectedProfile, nil
}

// GetConfigWithProfile loads AWS config with the selected profile
func GetConfigWithProfile(profile string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
}

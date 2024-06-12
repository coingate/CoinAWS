package aws_operations

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// Profile represents the structure of an AWS CLI profile
type Profile struct {
	Name string `json:"Name"`
}

// GetProfiles retrieves the list of profiles using AWS CLI
func GetProfiles() ([]string, error) {
	cmd := exec.Command("aws", "configure", "list-profiles")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	profiles := strings.Split(strings.TrimSpace(out.String()), "\n")
	return profiles, nil
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
	err = survey.AskOne(prompt, &selectedProfile, survey.WithPageSize(14))
	if err != nil {
		return "", err
	}

	return selectedProfile, nil
}

// GetConfigWithProfile loads AWS config with the selected profile and returns both default and shared configs
func GetConfigWithProfile(profile string) (aws.Config, config.SharedConfig, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile))
	if err != nil {
		return aws.Config{}, config.SharedConfig{}, err
	}

	sharedCfg, err := config.LoadSharedConfigProfile(context.TODO(), profile)
	if err != nil {
		return aws.Config{}, config.SharedConfig{}, err
	}

	return cfg, sharedCfg, nil
}

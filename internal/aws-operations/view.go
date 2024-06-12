package aws_operations

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// ViewSecretVersions prompts the user to select a secret version to view in read-only mode
func ViewSecretVersions(cfg aws.Config, secretName, defaultEditor string) error {
	// Create a new Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.DescribeSecretInput{
		SecretId: &secretName,
	}

	// Describe the secret to get version details
	result, err := svc.DescribeSecret(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "AccessDeniedException") {
			return fmt.Errorf("access denied: you do not have permission to describe the secret '%s'", secretName)
		}
		return fmt.Errorf("error describing secret: %w", err)
	}

	versionLabels := make(map[string]string)
	var labels []string
	for versionID, stages := range result.VersionIdsToStages {
		for _, stage := range stages {
			label := fmt.Sprintf("%s (Version ID: %s)", stage, versionID)
			labels = append(labels, label)
			versionLabels[label] = versionID
		}
	}

	// Prompt user to select a version to view
	var selectedLabel string
	prompt := &survey.Select{
		Message:  "Choose a version to view:",
		Options:  labels,
		PageSize: 14,
	}
	if err := survey.AskOne(prompt, &selectedLabel); err != nil {
		return fmt.Errorf("error selecting version: %w", err)
	}

	selectedVersionID := versionLabels[selectedLabel]

	viewInput := &secretsmanager.GetSecretValueInput{
		SecretId:  &secretName,
		VersionId: &selectedVersionID,
	}

	// Get the secret value of the selected version
	viewResult, err := svc.GetSecretValue(context.TODO(), viewInput)
	if err != nil {
		return fmt.Errorf("error getting secret value: %w", err)
	}

	secretValue := *viewResult.SecretString

	// Write the secret value to a temporary file
	tmpfile, err := os.CreateTemp("", "secret-*.txt")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Printf("Warning: failed to remove temporary file: %v", err)
		}
	}()

	if _, err := tmpfile.WriteString(secretValue); err != nil {
		return fmt.Errorf("error writing to temp file: %w", err)
	}
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("error closing temp file: %w", err)
	}

	// Open the secret in the preferred text editor in read-only mode
	editor := defaultEditor
	if editor == "" {
		editor = "nano" // Default to nano if no editor is set
	}

	cmd := exec.Command(editor, "-R", tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running editor: %w", err)
	}

	return nil
}

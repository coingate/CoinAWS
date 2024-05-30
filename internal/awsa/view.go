package awsa

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"io/ioutil"
	"os"
	"os/exec"
)

// ViewSecretVersions prompts the user to select a secret version to view in read-only mode
func ViewSecretVersions(cfg aws.Config, secretName, defaultEditor string) error {
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.DescribeSecretInput{
		SecretId: &secretName,
	}

	result, err := svc.DescribeSecret(context.TODO(), input)
	if err != nil {
		return err
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

	var selectedLabel string
	prompt := &survey.Select{
		Message: "Choose a version to view:",
		Options: labels,
	}
	err = survey.AskOne(prompt, &selectedLabel, survey.WithPageSize(14))
	if err != nil {
		return err
	}

	selectedVersionID := versionLabels[selectedLabel]

	viewInput := &secretsmanager.GetSecretValueInput{
		SecretId:  &secretName,
		VersionId: &selectedVersionID,
	}

	viewResult, err := svc.GetSecretValue(context.TODO(), viewInput)
	if err != nil {
		return err
	}

	secretValue := *viewResult.SecretString

	// Write the secret value to a temporary file
	tmpfile, err := ioutil.TempFile("", "secret-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // Clean up the file afterwards

	if _, err := tmpfile.WriteString(secretValue); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
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
		return err
	}

	return nil
}

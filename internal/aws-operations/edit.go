package aws_operations

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// EditSecret opens the secret in the preferred text editor and saves the changes
func EditSecret(cfg aws.Config, secretName, defaultEditor string) (string, string, error) {
	// Create a new Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	// Retrieve the secret value
	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "AccessDeniedException") {
			return "", "", fmt.Errorf("access denied: you do not have permission to access the secret '%s'", secretName)
		}
		return "", "", fmt.Errorf("error fetching the secret: %w", err)
	}

	originalSecret := *result.SecretString

	// Write the secret value to a temporary file
	tmpfile, err := os.CreateTemp("", "secret-*.txt")
	if err != nil {
		return "", "", err
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Printf("Warning: failed to remove temporary file: %v", err)
		}
	}()

	if _, err := tmpfile.WriteString(originalSecret); err != nil {
		return "", "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", "", err
	}

	// Open the secret in the preferred text editor
	editor := defaultEditor
	if editor == "" {
		editor = "nano" // Default to nano if no editor is set
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	// Read the modified secret back from the file
	editedSecretBytes, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", "", err
	}

	editedSecret := string(editedSecretBytes)

	return originalSecret, editedSecret, nil
}

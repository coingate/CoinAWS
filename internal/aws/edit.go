package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"io/ioutil"
	"os"
	"os/exec"
)

// EditSecret opens the secret in the preferred text editor and saves the changes
func EditSecret(cfg aws.Config, secretName, defaultEditor string) (string, string, error) {
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", "", err
	}

	originalSecret := *result.SecretString

	// Write the secret value to a temporary file
	tmpfile, err := ioutil.TempFile("", "secret-*.txt")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tmpfile.Name()) // Clean up the file afterwards

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
	editedSecretBytes, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return "", "", err
	}

	editedSecret := string(editedSecretBytes)

	return originalSecret, editedSecret, nil
}

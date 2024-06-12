package aws_operations

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// UpdateSecret updates the secret in AWS Secrets Manager
func UpdateSecret(cfg aws.Config, secretName, secretValue string, versionStages []string) error {
	// Create a new Secrets Manager client
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.PutSecretValueInput{
		SecretId:      &secretName,
		SecretString:  &secretValue,
		VersionStages: versionStages,
	}

	// Update the secret value
	_, err := svc.PutSecretValue(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

// GenerateVersionLabel generates a version label with the current timestamp
func GenerateVersionLabel() string {
	return time.Now().Format("20060102T150405")
}

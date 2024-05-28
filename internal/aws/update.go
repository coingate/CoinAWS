package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"time"
)

// UpdateSecret updates the secret in AWS Secrets Manager
func UpdateSecret(cfg aws.Config, secretName, secretValue string, versionStages []string) error {
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.PutSecretValueInput{
		SecretId:      &secretName,
		SecretString:  &secretValue,
		VersionStages: versionStages,
	}

	_, err := svc.PutSecretValue(context.TODO(), input)
	return err
}

// GenerateVersionLabel generates a version label with the current timestamp
func GenerateVersionLabel() string {
	return time.Now().Format("20060102T150405")
}

package aws_operations

import (
	"context"
	"sort"
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

	// Ensure we only keep 15 versions
	err = pruneOldVersions(cfg, secretName)
	if err != nil {
		return err
	}

	return nil
}

// GenerateVersionLabel generates a version label with the current timestamp
func GenerateVersionLabel() string {
	return time.Now().Format("20060102T150405")
}

// pruneOldVersions removes labels from the oldest secret versions, keeping only the latest 15
func pruneOldVersions(cfg aws.Config, secretName string) error {
	svc := secretsmanager.NewFromConfig(cfg)
	maxResults := int32(20)

	// Retrieve all secret versions
	input := &secretsmanager.ListSecretVersionIdsInput{
		SecretId:   &secretName,
		MaxResults: &maxResults,
	}

	result, err := svc.ListSecretVersionIds(context.TODO(), input)
	if err != nil {
		return err
	}

	// Sort versions by creation date
	sort.Slice(result.Versions, func(i, j int) bool {
		return result.Versions[i].CreatedDate.Before(*result.Versions[j].CreatedDate)
	})

	// If there are more than 15 versions, remove labels from the oldest ones
	if len(result.Versions) > 15 {
		for _, version := range result.Versions[:len(result.Versions)-15] {
			for _, versionStage := range version.VersionStages {
				_, err := svc.UpdateSecretVersionStage(context.TODO(), &secretsmanager.UpdateSecretVersionStageInput{
					SecretId:            &secretName,
					VersionStage:        aws.String(versionStage),
					RemoveFromVersionId: version.VersionId,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

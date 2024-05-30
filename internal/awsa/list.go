package awsa

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// ListSecrets lists all AWS Secrets Manager secrets for the given AWS config
func ListSecrets(cfg aws.Config) ([]string, error) {
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.ListSecretsInput{}

	var secrets []string
	paginator := secretsmanager.NewListSecretsPaginator(svc, input)

	for paginator.HasMorePages() {
		result, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, secret := range result.SecretList {
			secrets = append(secrets, *secret.Name)
		}
	}

	return secrets, nil
}

package awsa

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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

// ListEC2Instances lists all EC2 instances for the given AWS config
func ListEC2Instances(cfg aws.Config) ([]types.Instance, error) {
	svc := ec2.NewFromConfig(cfg)
	result, err := svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %w", err)
	}

	var instances []types.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	return instances, nil
}

// getInstanceName extracts the Name tag from an instance
func getInstanceName(instance types.Instance) string {
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return "Unnamed"
}

package awsa

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func CheckAndRefreshToken(profile string) error {
	ctx := context.TODO()

	// Load the Shared AWS Configuration with the specific profile
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create a STS client
	client := sts.NewFromConfig(cfg)

	// Try to get the caller identity to check if the token is valid
	_, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err == nil {
		// Token is valid
		return nil
	}

	// If there is an error, assume the token has expired
	fmt.Println("AWS SSO token has expired, please log in.")

	// Execute the AWS SSO login command
	cmd := exec.Command("aws", "sso", "login", "--profile", profile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute aws sso login: %w", err)
	}

	// Re-load the AWS config after login
	cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return fmt.Errorf("unable to reload SDK config after SSO login: %w", err)
	}

	// Verify the token again
	_, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("unable to verify token after SSO login: %w", err)
	}

	return nil
}

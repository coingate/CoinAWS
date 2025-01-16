package aws_operations

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
)

// ConnectToEC2Instance connects to an EC2 instance using AWS SSM
func ConnectToEC2Instance(cfg config.SharedConfig, instanceID string) error {
	ssmSigChan := make(chan os.Signal, 1)
	signal.Notify(ssmSigChan, syscall.SIGINT)
	defer signal.Stop(ssmSigChan)

	profile := cfg.Profile

	// Form the command to start an SSM session
	command := fmt.Sprintf("aws ssm start-session --target %s --profile %s", instanceID, profile)
	cmd := exec.Command("bash", "-c", command)

	// Capture the command's standard error output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Redirect the command's standard input and output to the current process
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	// Run the command
	if err := cmd.Run(); err != nil {
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "AccessDeniedException") {
			return fmt.Errorf("access denied: you do not have permission to connect to EC2 - '%s'", instanceID)
		}
		return fmt.Errorf("error connecting to instance: %w - %s", err, stderrStr)
	}

	return nil
}

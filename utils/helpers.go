package utils

import (
	"os/exec"
)

// List of known text editors
var knownEditors = []string{
	"vim", "nano", "vi", "emacs", "gedit", "kate", "subl", "code",
}

// CheckAWSCLI checks if the AWS CLI is installed
func CheckAWSCLI() bool {
	_, err := exec.LookPath("aws")
	return err == nil
}

// CheckSsmPlugin checks if the Session Manager plugin is installed
func CheckSsmPlugin() bool {
	_, err := exec.LookPath("session-manager-plugin")
	return err == nil
}

// DetectEditors detects available text editors from the known list
func DetectEditors() []string {
	var availableEditors []string
	for _, editor := range knownEditors {
		if _, err := exec.LookPath(editor); err == nil {
			availableEditors = append(availableEditors, editor)
		}
	}
	return availableEditors
}

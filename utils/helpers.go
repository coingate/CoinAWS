package utils

import (
	"os/exec"
)

var knownEditors = []string{"vim", "nano", "vi", "emacs", "gedit", "kate", "subl", "code"}

// Check if AWS CLI is installed
func CheckAWSCLI() bool {
	_, err := exec.LookPath("aws")
	return err == nil
}

// DetectEditors detects available text editors from the known list
func DetectEditors() []string {
	availableEditors := []string{}
	for _, editor := range knownEditors {
		if _, err := exec.LookPath(editor); err == nil {
			availableEditors = append(availableEditors, editor)
		}
	}
	return availableEditors
}

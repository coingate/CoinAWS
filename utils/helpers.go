package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

// VersionInfo Struct to parse version information from the GitHub API
type VersionInfo struct {
	TagName string `json:"tag_name"`
}

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

// CheckForUpdates compares the current version with the latest GitHub release
func CheckForUpdates(currentVersion string) {
	latestVersion, err := getLatestVersion()
	if err != nil {
		fmt.Println("Could not check for updates, please check your internet connection.")
		return
	}

	if isNewerVersion(currentVersion, latestVersion) {
		fmt.Printf("A newer version (%s) is available. Please upgrade to the latest version.\n", latestVersion)
		fmt.Println("To upgrade, run the following command:")
		fmt.Println("curl -sL https://raw.githubusercontent.com/coingate/CoinAWS/main/install.sh | sudo sh")
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("You are using the latest version.")
	}
}

// Fetch the latest version from GitHub Releases API
func getLatestVersion() (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://api.github.com/repos/coingate/CoinAWS/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get latest version, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var versionInfo VersionInfo
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		return "", err
	}

	return versionInfo.TagName, nil
}

// Compare the current version with the latest version
func isNewerVersion(currentVersion string, latestVersion string) bool {
	return currentVersion < latestVersion
}

package main

import (
	"flag"
	"fmt"
	"os"
	"smeditor/config"
	"smeditor/internal/awsa"
	"smeditor/utils"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	// Define the 'config' flag
	configFlag := flag.Bool("config", false, "Set the default editor")
	flag.Parse()

	// Load the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Handle the 'config' flag
	if *configFlag {
		availableEditors := utils.DetectEditors()
		if len(availableEditors) == 0 {
			fmt.Println("No known editors found on your system.")
			os.Exit(1)
		}

		var selectedEditor string
		prompt := &survey.Select{
			Message: "Choose your default text editor:",
			Options: availableEditors,
		}
		err := survey.AskOne(prompt, &selectedEditor)
		if err != nil {
			fmt.Println("Error selecting editor:", err)
			os.Exit(1)
		}

		cfg.DefaultEditor = selectedEditor
		err = config.SaveConfig(cfg)
		if err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}
		fmt.Printf("Default editor set to '%s'\n", cfg.DefaultEditor)
		os.Exit(0)
	}

	// Check if AWS CLI is installed
	if !utils.CheckAWSCLI() {
		fmt.Println("AWS CLI is not installed. Please install it and try again.")
		os.Exit(1)
	}

	// Select AWS profile
	profile, err := awsa.SelectProfile()
	if err != nil {
		fmt.Println("Error selecting profile:", err)
		os.Exit(1)
	}

	fmt.Println("Selected profile:", profile)

	// Get AWS config with the selected profile
	cfgAWS, err := awsa.GetConfigWithProfile(profile)
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		os.Exit(1)
	}

	for {
		// Select AWS secret and action
		secretName, action, err := awsa.SelectSecret(cfgAWS)
		if err != nil {
			fmt.Println("Error selecting secret:", err)
			os.Exit(1)
		}

		fmt.Println("Selected secret:", secretName)
		fmt.Println("Selected action:", action)

		if action == "Edit latest secret version" {
			// Edit AWS secret
			originalSecret, editedSecret, err := awsa.EditSecret(cfgAWS, secretName, cfg.DefaultEditor)
			if err != nil {
				fmt.Println("Error editing secret:", err)
				os.Exit(1)
			}

			// Compare original and edited secret values
			if originalSecret == editedSecret {
				fmt.Println("No changes made to the secret.")
				continue
			}

			// Prompt for version label
			fmt.Println("Enter version label (leave empty to use timestamp):")
			var versionLabel string
			fmt.Scanln(&versionLabel)

			if versionLabel == "" {
				versionLabel = awsa.GenerateVersionLabel()
			}

			// Update AWS secret with additional label "AWSCURRENT"
			versionStages := []string{versionLabel, "AWSCURRENT"}

			// Update AWS secret
			err = awsa.UpdateSecret(cfgAWS, secretName, editedSecret, versionStages)
			if err != nil {
				fmt.Println("Error updating secret:", err)
				os.Exit(1)
			}

			fmt.Println("Secret updated successfully with version labels:", versionStages)
		} else if action == "View previous versions" {
			// View AWS secret versions
			err := awsa.ViewSecretVersions(cfgAWS, secretName, cfg.DefaultEditor)
			if err != nil {
				fmt.Println("Error viewing secret versions:", err)
				os.Exit(1)
			}
		}
	}
}

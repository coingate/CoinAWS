package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"smeditor/internal/aws-operations"
	"smeditor/utils"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	configFlag := flag.Bool("config", false, "Set the default editor")
	flag.Parse()

	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if *configFlag {
		setDefaultEditor(cfg)
		return
	}

	if !utils.CheckAWSCLI() {
		log.Fatalf("AWS CLI is not installed. Please install it and try again.")
	}

	if !utils.CheckSsmPlugin() {
		log.Fatalf("Session Manager plugin is not installed. Please install it and try again.")
	}

	profile, err := aws_operations.SelectProfile()
	if err != nil {
		log.Fatalf("Error selecting profile: %v", err)
	}

	err = aws_operations.CheckAndRefreshToken(profile)
	if err != nil {
		log.Fatalf("Error checking or refreshing token: %v", err)
	}

	cfgAWS, sharedConfig, err := aws_operations.GetConfigWithProfile(profile)
	if err != nil {
		log.Fatalf("Error loading AWS config: %v", err)
	}

	mainMenu(cfgAWS, &cfg, sharedConfig)
}

func setDefaultEditor(cfg utils.Config) {
	availableEditors := utils.DetectEditors()
	if len(availableEditors) == 0 {
		log.Fatalf("No known editors found on your system.")
	}

	var selectedEditor string
	prompt := &survey.Select{
		Message:  "Choose your default text editor:",
		Options:  availableEditors,
		PageSize: 10,
	}
	err := survey.AskOne(prompt, &selectedEditor)
	if err != nil {
		log.Fatalf("Error selecting editor: %v", err)
	}

	cfg.DefaultEditor = selectedEditor
	err = utils.SaveConfig(cfg)
	if err != nil {
		log.Fatalf("Error saving config: %v", err)
	}
	fmt.Printf("Default editor set to '%s'\n", cfg.DefaultEditor)
}

func mainMenu(cfg aws.Config, appConfig *utils.Config, sharedConfig config.SharedConfig) {
	for {
		var option string
		prompt := &survey.Select{
			Message: "Main Menu",
			Options: []string{"EC2", "SecretsManager", "Exit"},
		}
		err := survey.AskOne(prompt, &option)
		if err != nil {
			log.Fatalf("Error selecting option: %v", err)
		}

		switch option {
		case "EC2":
			handleEC2(cfg, sharedConfig)
		case "SecretsManager":
			handleSecretsManager(cfg, appConfig)
		case "Exit":
			fmt.Println("Exiting.")
			os.Exit(0)
		}
	}
}

func handleEC2(cfg aws.Config, sharedConfig config.SharedConfig) {
	instances, err := aws_operations.ListEC2Instances(cfg)
	if err != nil {
		log.Fatalf("Error listing EC2 instances: %v", err)
	}
	if len(instances) == 0 {
		fmt.Println("No EC2 instances found.")
		return
	}

	var maxNameLen int
	instanceMap := make(map[string]string)
	for _, instance := range instances {
		var name string
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				name = *tag.Value
				break
			}
		}
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
		instanceMap[name] = *instance.InstanceId
	}

	var instanceOptions []string
	for name, instanceID := range instanceMap {
		instanceOption := fmt.Sprintf("%-*s %s", maxNameLen, name, instanceID)
		instanceOptions = append(instanceOptions, instanceOption)
	}

	var selectedInstanceOption string
	prompt := &survey.Select{
		Message:  "Select an EC2 instance to connect to:",
		Options:  instanceOptions,
		PageSize: 14,
	}
	err = survey.AskOne(prompt, &selectedInstanceOption)
	if err != nil {
		log.Fatalf("Error selecting instance: %v", err)
	}

	selectedInstanceID := strings.TrimSpace(selectedInstanceOption[maxNameLen:])
	err = aws_operations.ConnectToEC2Instance(sharedConfig, selectedInstanceID)
	if err != nil {
		log.Println(err)
		fmt.Println("Returning to the main menu...")
		return
	}
}

func handleSecretsManager(cfg aws.Config, appConfig *utils.Config) {
	secretName, action, err := aws_operations.SelectSecret(cfg)
	if err != nil {
		log.Fatalf("Error selecting secret: %v", err)
	}

	fmt.Println("Selected secret:", secretName)
	fmt.Println("Selected action:", action)

	switch action {
	case "Edit latest secret version":
		editSecret(cfg, appConfig, secretName)
	case "View previous versions":
		err := aws_operations.ViewSecretVersions(cfg, secretName, appConfig.DefaultEditor)
		if err != nil {
			log.Println(err)
			fmt.Println("Returning to the main menu...")
			return
		}
	}
}

func editSecret(cfg aws.Config, appConfig *utils.Config, secretName string) {
	originalSecret, editedSecret, err := aws_operations.EditSecret(cfg, secretName, appConfig.DefaultEditor)
	if err != nil {
		log.Println(err)
		fmt.Println("Returning to the main menu...")
		return
	}

	if originalSecret == editedSecret {
		fmt.Println("No changes made to the secret.")
		return
	}

	fmt.Println("Enter version label (leave empty to use timestamp):")
	reader := bufio.NewReader(os.Stdin)
	versionLabel, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading version label: %v", err)
	}

	// Trim the newline character from the input
	versionLabel = strings.TrimSpace(versionLabel)

	if versionLabel == "" {
		versionLabel = aws_operations.GenerateVersionLabel()
	}

	versionStages := []string{versionLabel, "AWSCURRENT"}

	err = aws_operations.UpdateSecret(cfg, secretName, editedSecret, versionStages)
	if err != nil {
		log.Fatalf("Error updating secret: %v", err)
	}

	fmt.Println("Secret updated successfully with version labels:", versionStages)
}

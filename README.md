# CoinAWS

`CoinAWS` is a command-line interface (CLI) tool for managing AWS resources such as EC2, SSM, and Secrets Manager. It
allows users to select AWS profiles, view and edit secrets, manage instances, and handle other AWS operations in an
interactive manner.

## Features

- Interactive selection of AWS profiles
- Viewing and editing AWS Secrets Manager secrets
- Managing EC2 instances
- Managing secret versions
- Customizable default text editor for editing secrets

## Installation

### One-liner Installation

You can install `CoinAWS` directly using a one-liner:

```sh
curl -sL https://raw.githubusercontent.com/coingate/CoinAWS/main/install.sh | sudo sh
```

### Building from Source

To build CoinAWS from source, you need to have Go installed on your system. You can then build and install the tool
using the following steps:
1. Clone the repository:

    ```sh
    git clone https://github.com/coingate/CoinAWS.git
    cd CoinAWS
    ```

2. Build the binary:

    ```sh
    go build -o coinaws
    ```

3. Optionally, move the binary to `/usr/local/bin` to make it accessible from anywhere:

    ```sh
    sudo mv coinaws /usr/local/bin/
    ```

## Usage

### Setting the Default Editor

To set the default text editor for editing secrets, use the `--config` flag:

```sh
coinaws --config
```

You will be prompted to choose your preferred text editor from a list of available editors on your system.

### Starting the Application

To start the coinaws application and manage your AWS resources:

1. #### Run the application:

    ```sh
    ./coinaws
    ```
   If you moved the binary to /usr/local/bin, you can simply run:

    ```sh
    coinaws
    ```

2. #### Select an AWS profile:

   Upon running coinaws, you will be prompted to choose an AWS profile from your AWS configuration.

3. #### Choose an action:

   After selecting a profile, you can choose from various options such as managing EC2 instances, viewing/editing
   secrets, etc.

### Connecting to EC2 Instances

You can manage EC2 instances by selecting the EC2 option from the main menu. This allows you to start SSM session with
selected instance.

### Editing a Secret

When you choose to edit the latest version of a secret, the secret will be opened in your default text editor. After editing and saving the file, you will be prompted to enter a version label (or use a timestamp). The secret will be updated in AWS Secrets Manager with the new value and version label.

### Viewing Older Versions

When you choose to view older versions of a secret, you will be prompted to select a version. The selected version will be opened in read-only mode for viewing.

## Project Structure

- `main.go`: The main entry point for the application. Handles command-line arguments and coordinates the overall flow.
- `install.sh`: A script for installing the `coinaws` binary.
- `internal/aws-operations/`: Contains the core logic for interacting with AWS services and handling user interactions.
  - `profile.go`: Handles AWS profile selection.
  - `view.go`: Handles viewing specific versions of secrets.
  - `edit.go`: Handles editing secrets.
  - `check_token.go`: Handles checking and refreshing AWS SSO tokens.
  - `select.go`: Handles resource selection.
  - `connect.go`: Handles connecting to EC2 instances.
  - `update.go`: Handles updating secrets and generating version labels.
  - `list.go`: Handles listing resources and their details.
- `utils/config.go`: Manages loading and saving configuration settings.
- `utils/helpers.go`: Contains utility functions, such as detecting available text editors and checking for AWS CLI installation.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any suggestions or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
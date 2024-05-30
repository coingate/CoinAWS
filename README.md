# smeditor

`smeditor` is a command-line interface (CLI) tool for managing AWS Secrets Manager secrets. It allows users to select AWS profiles, view and edit secrets, and handle secret versions in an interactive manner.

## Features

- Interactive selection of AWS profiles
- Viewing and editing AWS Secrets Manager secrets
- Managing secret versions
- Customizable default text editor for editing secrets

## Installation

### One-liner Installation

You can install `smeditor` directly using a one-liner:

```sh
curl -sL https://raw.githubusercontent.com/coingate/smeditor/main/install.sh | sudo sh
```

### Building from Source

To build smeditor from source, you need to have Go installed on your system. You can then build and install the tool using the following steps:
1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/smeditor.git
    cd smeditor
    ```

2. Build the binary:

    ```sh
    go build -o smeditor
    ```

3. Optionally, move the binary to `/usr/local/bin` to make it accessible from anywhere:

    ```sh
    sudo mv smeditor /usr/local/bin/
    ```


## Usage

### Setting the Default Editor

To set the default text editor for editing secrets, use the `--config` flag:

```sh
smeditor --config
```

You will be prompted to choose your preferred text editor from a list of available editors on your system.

### Starting the Application

To start the smeditor application and manage your AWS Secrets:

1. #### Run the application:

    ```sh
    ./smeditor
    ```
   If you moved the binary to /usr/local/bin, you can simply run:

    ```sh
    smeditor
    ```

2. #### Select an AWS profile:

    Upon running smeditor, you will be prompted to choose an AWS profile from your AWS configuration.

3. #### Select a secret:

    Next, you will be prompted to choose a secret from AWS Secrets Manager.

4. #### Choose an action:

    After selecting a secret, you can choose to either "Edit latest version" or "View older versions".

### Editing a Secret

When you choose to edit the latest version of a secret, the secret will be opened in your default text editor. After editing and saving the file, you will be prompted to enter a version label (or use a timestamp). The secret will be updated in AWS Secrets Manager with the new value and version label.

### Viewing Older Versions

When you choose to view older versions of a secret, you will be prompted to select a version. The selected version will be opened in read-only mode for viewing.

## Project Structure

- `main.go`: The main entry point for the application. Handles command-line arguments and coordinates the overall flow.
- `install.sh`: A script for installing the `smeditor` binary.
- `internal/awsa/`: Contains the core logic for interacting with AWS Secrets Manager and handling user interactions.
  - `profile.go`: Handles AWS profile selection.
  - `view.go`: Handles viewing specific versions of secrets.
  - `edit.go`: Handles editing secrets.
  - `check_token.go`: Handles checking and refreshing AWS SSO tokens.
  - `select.go`: Handles secret selection.
  - `update.go`: Handles updating secrets and generating version labels.
  - `list.go`: Handles listing secrets and their versions.
- `config/config.go`: Manages loading and saving configuration settings.
- `utils/helpers.go`: Contains utility functions, such as detecting available text editors and checking for AWS CLI installation.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any suggestions or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



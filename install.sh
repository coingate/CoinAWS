#!/bin/sh

set -e

# Determine the OS
OS="$(uname -s)"
case "$OS" in
    Linux*)  OS=linux;;
    Darwin*) OS=macos;;
    *)       echo "Unsupported OS: $OS"; exit 1;;
esac

# Determine the latest version tag from GitHub and find the corresponding asset URL
ASSET_URL=$(curl -s https://api.github.com/repos/coingate/smeditor/releases/latest | grep "browser_download_url.*${OS}.*tar.gz" | cut -d '"' -f 4)

# Check if the ASSET_URL was found
if [ -z "$ASSET_URL" ]; then
    echo "Failed to find the download URL for the latest release."
    exit 1
fi

# Download the appropriate tar.gz file
curl -L -o smeditor.tar.gz "$ASSET_URL"

# Extract the downloaded file
tar -xzf smeditor.tar.gz

# Move the binary to /usr/local/bin
mv smeditor /usr/local/bin/

# Set the correct ownership and permissions
chmod 755 /usr/local/bin/smeditor

# Clean up
rm -rf smeditor.tar.gz smeditor

# Determine the default shell for the current user
DEFAULT_SHELL=$(getent passwd "$USER" | cut -d: -f7)

# Fallback if getent is not available (common on macOS)
if [ -z "$DEFAULT_SHELL" ]; then
    DEFAULT_SHELL=$(dscl . -read /Users/"$USER" UserShell | awk '{print $2}')
fi

# Refresh the appropriate shell configuration file
case "$DEFAULT_SHELL" in
    */bash)
        if [ -f "$HOME/.bashrc" ]; then
            . "$HOME/.bashrc"
        elif [ -f "$HOME/.bash_profile" ]; then
            . "$HOME/.bash_profile"
        fi
        ;;
    */zsh)
        if [ -f "$HOME/.zshrc" ]; then
            . "$HOME/.zshrc"
        fi
        ;;
    */ksh)
        if [ -f "$HOME/.kshrc" ]; then
            . "$HOME/.kshrc"
        fi
        ;;
    */sh)
        if [ -f "$HOME/.profile" ]; then
            . "$HOME/.profile"
        fi
        ;;
    *)
        echo "Unsupported shell: $DEFAULT_SHELL"
        ;;
esac

# Verify installation
if command -v smeditor > /dev/null; then
    echo "smeditor installed successfully!"
else
    echo "Installation failed."
fi

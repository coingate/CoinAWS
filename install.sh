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
ASSET_URL=$(curl -s https://api.github.com/repos/coingate/CoinAWS/releases/latest | grep "browser_download_url.*${OS}.*tar.gz" | cut -d '"' -f 4)

# Check if the ASSET_URL was found
if [ -z "$ASSET_URL" ]; then
    echo "Failed to find the download URL for the latest release."
    exit 1
fi

# Download the appropriate tar.gz file
curl -L -o coinaws.tar.gz "$ASSET_URL"

# Extract the downloaded file
tar -xzf coinaws.tar.gz

# Move the binary to /usr/local/bin
mv coinaws /usr/local/bin/

# Set the correct ownership and permissions
chmod 755 /usr/local/bin/coinaws

# Clean up
rm -rf coinaws.tar.gz coinaws

# Determine the default shell for the original user (not root)
DEFAULT_SHELL=$(getent passwd "$SUDO_USER" | cut -d: -f7)

# Fallback if getent is not available (common on macOS)
if [ -z "$DEFAULT_SHELL" ]; then
    DEFAULT_SHELL=$(dscl . -read /Users/"$SUDO_USER" UserShell | awk '{print $2}')
fi

# Refresh the appropriate shell configuration file for the original user
case "$DEFAULT_SHELL" in
    */bash)
        if [ -f "/home/$SUDO_USER/.bashrc" ]; then
            sudo -u "$SUDO_USER" bash -c "source /home/$SUDO_USER/.bashrc"
        elif [ -f "/home/$SUDO_USER/.bash_profile" ]; then
            sudo -u "$SUDO_USER" bash -c "source /home/$SUDO_USER/.bash_profile"
        fi
        ;;
    */zsh)
        if [ -f "/home/$SUDO_USER/.zshrc" ]; then
            sudo -u "$SUDO_USER" zsh -c "source /home/$SUDO_USER/.zshrc"
        fi
        ;;
    */ksh)
        if [ -f "/home/$SUDO_USER/.kshrc" ]; then
            sudo -u "$SUDO_USER" ksh -c "source /home/$SUDO_USER/.kshrc"
        fi
        ;;
    */sh)
        if [ -f "/home/$SUDO_USER/.profile" ]; then
            sudo -u "$SUDO_USER" sh -c ". /home/$SUDO_USER/.profile"
        fi
        ;;
    *)
        echo "Unsupported shell: $DEFAULT_SHELL"
        ;;
esac

# Verify installation
if command -v coinaws > /dev/null; then
    echo "coinaws installed successfully!"
else
    echo "Installation failed."
fi

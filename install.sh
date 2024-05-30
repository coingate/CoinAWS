#!/bin/sh

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
tar -xzf smeditor.tar.gz -C smeditor --strip-components=1

# Move the binary to /usr/local/bin
sudo mv smeditor/smeditor /usr/local/bin/

# Clean up
rm -rf smeditor.tar.gz smeditor

# Verify installation
if command -v smeditor > /dev/null; then
    echo "smeditor installed successfully!"
else
    echo "Installation failed."
fi

#!/bin/bash

# Function to ensure curl is installed
has_curl() {
    _=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        exit 1
    fi
}

# Function to determine the architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64)        ARCH=amd64;;
        aarch64|arm64) ARCH=arm64;;
        *)         echo "Unsupported architecture"; exit 1;;
    esac
    echo "Detected Architecture: $ARCH"
}

# Function to determine the OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     OS=linux;;
        Darwin*)    OS=darwin;;
        CYGWIN*|MSYS*|MINGW32*|MINGW64*|MINGW*) OS=windows;;
        *)          echo "Unsupported OS"; exit 1;;
    esac
    echo "Detected OS: $OS"
}

# Version of the OpsLevel CLI to install
get_version() {
    VERSION=$(curl -sI https://github.com/OpsLevel/cli/releases/latest | grep -i "location:" | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')
    if [ ! $VERSION ]; then
        echo "Failed while attempting to install OpsLevel's cli. Please manually install:"
        echo ""
        echo "Open your web browser and go to https://github.com/OpsLevel/cli?tab=readme-ov-file#installation for instructions."
        exit 1
    fi
}

# Function to download and install the CLI tool
install_cli() {
    DOWNLOAD_URL="https://github.com/OpsLevel/cli/releases/download/${VERSION}/opslevel-${OS}-${ARCH}.tar.gz"

    # Temporary directory to store the download
    TMP_DIR=$(mktemp -d)

    echo "Downloading OpsLevel CLI from $DOWNLOAD_URL ..."
    curl -L -o "$TMP_DIR/opslevel.tar.gz" "$DOWNLOAD_URL"

    if [ $? -ne 0 ]; then
        echo "Download failed. Please check the version and try again."
        exit 1
    fi

    echo "Extracting the OpsLevel CLI..."
    tar -xzf "$TMP_DIR/opslevel.tar.gz" -C "$TMP_DIR"

    echo "Installing the OpsLevel CLI to /usr/local/bin ..."
    sudo mv "$TMP_DIR/opslevel" /usr/local/bin/

    if [ $? -ne 0 ]; then
        echo "Installation failed."
        exit 1
    fi

    echo "Cleaning up..."
    rm -rf "$TMP_DIR"

    echo "OpsLevel CLI installed successfully!"
}

# Main script execution
has_curl
detect_arch
detect_os
get_version
install_cli

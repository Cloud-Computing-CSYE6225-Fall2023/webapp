#!/bin/bash

# Check root privileges
if [ "$EUID" -ne 0 ]; then
    echo "Please run this script as root or with sudo."
    exit 1
fi

# Update packages
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get clean

# Install Go
echo "Installing the latest version of Go..."
wget https://golang.org/dl/go1.21.1.linux-amd64.tar.gz -O /tmp/go.tar.gz
tar -C /usr/local -xzf /tmp/go.tar.gz
rm /tmp/go.tar.gz

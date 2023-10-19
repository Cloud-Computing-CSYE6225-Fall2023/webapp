#!/bin/bash

# Setup Go Path
echo "Setting up Go environment variables..."
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/home/admin/github.com/shivasaicharanruthala

echo "Check Go version"
go version

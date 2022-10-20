#!/usr/bin/env bash

# NOTE --always falls back to a commit hash when no tags are available 
TAG=$(git describe --abbrev=0 --tags --always | head -n1)
BINARY_NAME="wicket"

set -e
go mod tidy
go build -o $BINARY_NAME 

GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}.linux.amd64
GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}.linux.arm64
GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}.darwin.amd64
GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}.darwin.arm64
GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}.windows.amd64.exe

gh release create -p --generate-notes "${TAG}" ${BINARY_NAME}.*
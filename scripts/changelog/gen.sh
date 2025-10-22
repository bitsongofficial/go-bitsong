#!/bin/bash

# Check if git-cliff is installed
if ! command -v git-cliff &> /dev/null; then
    echo "git-cliff not found. Attempting to install via Cargo..."
    # Check if Cargo is installed
    if ! command -v cargo &> /dev/null; then
        echo "Error: Cargo is not installed. Please install Rust and Cargo first to install git-cliff."
        exit 1
    fi
    # Install git-cliff using Cargo
    cargo install git-cliff
    # Verify installation was successful
    if ! command -v git-cliff &> /dev/null; then
        echo "Error: Failed to install git-cliff. Please check your Cargo setup or install it manually."
        exit 1
    fi
    echo "git-cliff installed successfully."
fi

# Generate the latest changelog and prepend it to CHANGELOG.md
echo "generating changelog.."
git-cliff --config scripts/changelog/cliff.toml -u
echo "changelog generated"

# git cliff --config scripts/changelog/cliff.toml -l --prepend CHANGELOG.md
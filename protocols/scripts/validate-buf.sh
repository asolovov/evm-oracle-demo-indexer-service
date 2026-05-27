#!/bin/bash

set -e

# Script to validate proto files using Buf
# Usage: ./scripts/validate-buf.sh

echo "Validating proto files with Buf..."
echo ""

# Check if buf is installed
if ! command -v buf &> /dev/null; then
    echo "Error: buf is not installed"
    echo "Install it with: make buf-install"
    exit 1
fi

# Run buf lint
echo "Running buf lint..."
if buf lint; then
    echo "✓ Buf lint passed"
else
    echo "✗ Buf lint failed"
    exit 1
fi

echo ""

# Check for breaking changes if we have a git history
if git rev-parse --verify HEAD >/dev/null 2>&1; then
    # Get the latest tag
    LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    
    if [ -n "$LATEST_TAG" ]; then
        echo "Checking for breaking changes against tag: $LATEST_TAG"
        if buf breaking --against ".git#tag=$LATEST_TAG"; then
            echo "✓ No breaking changes detected"
        else
            echo "✗ Breaking changes detected"
            echo ""
            echo "Breaking changes are not allowed without a major version bump."
            echo "If this is intentional, update the version and create a new tag."
            exit 1
        fi
    else
        echo "No tags found, skipping breaking change detection"
    fi
else
    echo "Not a git repository, skipping breaking change detection"
fi

echo ""
echo "✓ All Buf validations passed!"

#!/bin/bash

set -e

# Script to validate proto files using protoc
# Usage: ./scripts/validate-protoc.sh [package-name]

PACKAGE=$1

echo "Validating proto files with protoc..."
echo ""

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Install it from https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Validate specific package or all packages
if [ -n "$PACKAGE" ]; then
    echo "Validating package: $PACKAGE"
    PROTO_FILES=$(find "$PACKAGE" -name "*.proto" 2>/dev/null || true)
    
    if [ -z "$PROTO_FILES" ]; then
        echo "Error: No proto files found in $PACKAGE"
        exit 1
    fi
    
    for proto in $PROTO_FILES; do
        echo "  Validating $proto..."
        protoc --proto_path=. --proto_path="$PACKAGE" --descriptor_set_out=/dev/null "$proto"
    done
    
    echo "✓ Package $PACKAGE validated successfully"
else
    echo "Validating all proto files..."
    FAILED=0
    
    for dir in common user; do
        if [ -d "$dir" ]; then
            echo "  Validating $dir..."
            PROTO_FILES=$(find "$dir" -name "*.proto" 2>/dev/null || true)
            
            for proto in $PROTO_FILES; do
                if ! protoc --proto_path=. --descriptor_set_out=/dev/null "$proto" 2>/dev/null; then
                    echo "  ✗ Failed: $proto"
                    FAILED=1
                else
                    echo "  ✓ Valid: $proto"
                fi
            done
        fi
    done
    
    if [ $FAILED -eq 1 ]; then
        echo ""
        echo "✗ Validation failed for some files"
        exit 1
    fi
    
    echo ""
    echo "✓ All proto files validated successfully"
fi

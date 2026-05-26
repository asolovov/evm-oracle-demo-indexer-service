#!/bin/bash

set -e

# Script to scaffold a new service structure
# Usage: ./scripts/add-service.sh <service-name>

SERVICE_NAME=$1

if [ -z "$SERVICE_NAME" ]; then
    echo "Error: Service name is required"
    echo "Usage: ./scripts/add-service.sh <service-name>"
    echo "Example: ./scripts/add-service.sh product"
    exit 1
fi

# Convert service name to lowercase
SERVICE_NAME=$(echo "$SERVICE_NAME" | tr '[:upper:]' '[:lower:]')

# Check if service already exists
if [ -d "$SERVICE_NAME" ]; then
    echo "Error: Service '$SERVICE_NAME' already exists"
    exit 1
fi

echo "Creating new service: $SERVICE_NAME"

# Create directory structure
mkdir -p "$SERVICE_NAME/v1"

# Create proto file with template
PROTO_FILE="$SERVICE_NAME/v1/$SERVICE_NAME.proto"
PACKAGE_NAME="${SERVICE_NAME}.v1"
GO_PACKAGE="github.com/asolovov/evm-oracle-demo-protocols/${SERVICE_NAME}/v1;${SERVICE_NAME}v1"

# Convert service name to PascalCase for message names
SERVICE_PASCAL=$(echo "$SERVICE_NAME" | sed -r 's/(^|_)([a-z])/\U\2/g')

cat > "$PROTO_FILE" <<EOF
syntax = "proto3";

package $PACKAGE_NAME;

import "google/protobuf/timestamp.proto";

option go_package = "$GO_PACKAGE";

// ${SERVICE_PASCAL}Service provides operations for managing ${SERVICE_NAME}s.
service ${SERVICE_PASCAL}Service {
  // Get${SERVICE_PASCAL} retrieves a ${SERVICE_NAME} by ID.
  rpc Get${SERVICE_PASCAL}(Get${SERVICE_PASCAL}Request) returns (Get${SERVICE_PASCAL}Response);
  
  // Create${SERVICE_PASCAL} creates a new ${SERVICE_NAME}.
  rpc Create${SERVICE_PASCAL}(Create${SERVICE_PASCAL}Request) returns (Create${SERVICE_PASCAL}Response);
}

// ${SERVICE_PASCAL} represents a ${SERVICE_NAME} entity.
message ${SERVICE_PASCAL} {
  // Unique ${SERVICE_NAME} identifier.
  string id = 1;
  
  // ${SERVICE_PASCAL} name.
  string name = 2;
  
  // Timestamp when the ${SERVICE_NAME} was created.
  google.protobuf.Timestamp created_at = 3;
  
  // Timestamp when the ${SERVICE_NAME} was last updated.
  google.protobuf.Timestamp updated_at = 4;
}

// Get${SERVICE_PASCAL}Request is the request message for Get${SERVICE_PASCAL}.
message Get${SERVICE_PASCAL}Request {
  // ID of the ${SERVICE_NAME} to retrieve.
  string id = 1;
}

// Get${SERVICE_PASCAL}Response is the response message for Get${SERVICE_PASCAL}.
message Get${SERVICE_PASCAL}Response {
  // The requested ${SERVICE_NAME}.
  ${SERVICE_PASCAL} ${SERVICE_NAME} = 1;
}

// Create${SERVICE_PASCAL}Request is the request message for Create${SERVICE_PASCAL}.
message Create${SERVICE_PASCAL}Request {
  // ${SERVICE_PASCAL} name (required).
  string name = 1;
}

// Create${SERVICE_PASCAL}Response is the response message for Create${SERVICE_PASCAL}.
message Create${SERVICE_PASCAL}Response {
  // The newly created ${SERVICE_NAME}.
  ${SERVICE_PASCAL} ${SERVICE_NAME} = 1;
}
EOF

# Create README
README_FILE="$SERVICE_NAME/README.md"
cat > "$README_FILE" <<EOF
# ${SERVICE_PASCAL} Service

## Overview

The ${SERVICE_PASCAL} Service provides gRPC APIs for managing ${SERVICE_NAME}s.

## Available Methods

### Get${SERVICE_PASCAL}

Retrieves a ${SERVICE_NAME} by ID.

**Request:**
\`\`\`protobuf
message Get${SERVICE_PASCAL}Request {
  string id = 1;
}
\`\`\`

**Response:**
\`\`\`protobuf
message Get${SERVICE_PASCAL}Response {
  ${SERVICE_PASCAL} ${SERVICE_NAME} = 1;
}
\`\`\`

### Create${SERVICE_PASCAL}

Creates a new ${SERVICE_NAME}.

**Request:**
\`\`\`protobuf
message Create${SERVICE_PASCAL}Request {
  string name = 1;
}
\`\`\`

**Response:**
\`\`\`protobuf
message Create${SERVICE_PASCAL}Response {
  ${SERVICE_PASCAL} ${SERVICE_NAME} = 1;
}
\`\`\`

## Usage Example

\`\`\`go
import ${SERVICE_NAME}v1 "github.com/asolovov/evm-oracle-demo-protocols/${SERVICE_NAME}/v1"

// Create a new ${SERVICE_NAME}
req := &${SERVICE_NAME}v1.Create${SERVICE_PASCAL}Request{
    Name: "Example ${SERVICE_PASCAL}",
}

resp, err := client.Create${SERVICE_PASCAL}(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created ${SERVICE_NAME}: %s\n", resp.${SERVICE_PASCAL}.Id)
\`\`\`

## Integration

To use this service in your microservice:

1. Pull the protocols repository:
   \`\`\`bash
   make proto-update
   \`\`\`

2. Generate Go code:
   \`\`\`bash
   make proto-generate PROTO_PACKAGE=${SERVICE_NAME}
   \`\`\`

3. Import in your code:
   \`\`\`go
   import ${SERVICE_NAME}v1 "your-module/protocols/${SERVICE_NAME}/v1"
   \`\`\`
EOF

echo ""
echo "✓ Service '$SERVICE_NAME' created successfully!"
echo ""
echo "Next steps:"
echo "  1. Review and customize $PROTO_FILE"
echo "  2. Run 'make buf-lint' to validate"
echo "  3. Run 'make buf-generate PACKAGE=$SERVICE_NAME' to generate Go code"
echo "  4. Commit your changes"

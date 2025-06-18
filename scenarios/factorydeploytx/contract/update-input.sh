#!/bin/bash

# Script to update CREATE2Factory.input.json from CREATE2Factory.sol
set -e

# Define file paths
SOL_FILE="CREATE2Factory.sol"
JSON_FILE="CREATE2Factory.input.json"

# Check if files exist
if [[ ! -f "$SOL_FILE" ]]; then
    echo "Error: $SOL_FILE not found!"
    exit 1
fi

if [[ ! -f "$JSON_FILE" ]]; then
    echo "Error: $JSON_FILE not found!"
    exit 1
fi

echo "Updating $JSON_FILE from $SOL_FILE..."

# Read the Solidity file content and escape it for JSON
# This handles newlines, quotes, backslashes, and other special characters
sol_content=$(cat "$SOL_FILE" | jq -Rs .)

# Create a temporary file with the updated JSON
temp_file=$(mktemp)

# Use jq to update the JSON file with the new Solidity content
jq --argjson content "$sol_content" '.sources."CREATE2Factory.sol".content = $content' "$JSON_FILE" > "$temp_file"

# Replace the original file with the updated content
mv "$temp_file" "$JSON_FILE"

echo "Successfully updated $JSON_FILE with content from $SOL_FILE"
echo "Updated contract source code in JSON file." 
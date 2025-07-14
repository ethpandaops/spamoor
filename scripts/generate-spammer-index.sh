#!/bin/bash

# Utility script to generate spammer index
# Usage: ./generate-spammer-index.sh <CONFIG_DIR>

if [ $# -lt 1 ]; then
    echo "Usage: $0 <CONFIG_DIR>"
    echo ""
    echo "Example:"
    echo "  $0 spammer-configs"
    echo ""
    echo "Arguments:"
    echo "  CONFIG_DIR - Directory containing spammer configurations"
    exit 1
fi

CONFIG_DIR="$1"

# Change to script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR/.." || exit 1

# Run the Go script
go run scripts/generate-spammer-index/generate-spammer-index.go "$CONFIG_DIR"